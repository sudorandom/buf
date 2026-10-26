// Copyright 2020-2026 Buf Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package filewatch

import (
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"buf.build/go/standard/xlog/xslog"
	"github.com/bufbuild/buf/private/pkg/normalpath"
	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	logger          *slog.Logger
	fsnotifyWatcher *fsnotify.Watcher

	eventChan chan Event
	errChan   chan error
	// closeChan is closed by Close so that run stops sending on eventChan and errChan.
	closeChan chan struct{}
	closeOnce sync.Once
	closeErr  error
	waitGroup sync.WaitGroup

	lock sync.RWMutex
	// rootDirPaths are the directories added with AddDirPath.
	rootDirPaths map[string]struct{}
	// recursiveDirPaths are the directories whose direct children are reported,
	// and whose new subdirectories are watched automatically.
	//
	// This includes rootDirPaths and their subdirectories that are not excluded.
	recursiveDirPaths map[string]struct{}
	// pendingDirPaths are the non-existent directories (rootDirPaths or their
	// missing ancestors) whose parent directories are watched for their creation.
	pendingDirPaths map[string]struct{}
	// filePaths are the individual files that are reported.
	filePaths map[string]struct{}
	// excludeDirPaths are the directories given to SetExcludeDirPaths, in
	// normalpath form for use with normalpath.EqualsOrContainsPath.
	excludeDirPaths []string
	// activeExcludeDirPaths are the excludeDirPaths that apply, which are those
	// that do not contain a root directory.
	activeExcludeDirPaths []string
}

func newWatcher(logger *slog.Logger) (*watcher, error) {
	fsnotifyWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("create file watcher: %w", err)
	}
	watcher := &watcher{
		logger:            logger,
		fsnotifyWatcher:   fsnotifyWatcher,
		eventChan:         make(chan Event, 64),
		errChan:           make(chan error),
		closeChan:         make(chan struct{}),
		rootDirPaths:      make(map[string]struct{}),
		recursiveDirPaths: make(map[string]struct{}),
		pendingDirPaths:   make(map[string]struct{}),
		filePaths:         make(map[string]struct{}),
	}
	watcher.waitGroup.Go(watcher.run)
	return watcher, nil
}

func (w *watcher) AddDirPath(dirPath string) error {
	absDirPath, err := filepath.Abs(dirPath)
	if err != nil {
		return err
	}
	// WalkDir does not follow a root directory that is a symbolic link.
	fileInfo, err := os.Lstat(absDirPath)
	if err == nil && fileInfo.Mode()&fs.ModeSymlink != 0 {
		absDirPath, err = filepath.EvalSymlinks(absDirPath)
		if err != nil {
			return err
		}
	}
	w.lock.Lock()
	if _, ok := w.rootDirPaths[absDirPath]; ok {
		w.lock.Unlock()
		return nil
	}
	w.rootDirPaths[absDirPath] = struct{}{}
	// An excluded directory that contains the new root no longer applies.
	excludesChanged := w.updateActiveExcludeDirPathsLocked()
	w.lock.Unlock()
	if excludesChanged {
		err = w.watchRootDirPaths()
	} else {
		err = w.watchRootDirPath(absDirPath, false)
	}
	if err != nil {
		// The root may be added again to retry.
		w.lock.Lock()
		delete(w.rootDirPaths, absDirPath)
		w.lock.Unlock()
		return err
	}
	return nil
}

func (w *watcher) AddFilePath(filePath string) error {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return err
	}
	w.lock.Lock()
	defer w.lock.Unlock()
	// The directory is always added, as its watch is lost when it is removed.
	if err := w.fsnotifyWatcher.Add(filepath.Dir(absFilePath)); err != nil {
		return fmt.Errorf("watch %s: %w", absFilePath, err)
	}
	w.filePaths[absFilePath] = struct{}{}
	return nil
}

func (w *watcher) SetExcludeDirPaths(dirPaths ...string) error {
	normalizedExcludeDirPaths := make([]string, 0, len(dirPaths))
	for _, dirPath := range dirPaths {
		absDirPath, err := filepath.Abs(dirPath)
		if err != nil {
			return err
		}
		normalizedExcludeDirPaths = append(normalizedExcludeDirPaths, normalpath.Normalize(absDirPath))
	}
	slices.Sort(normalizedExcludeDirPaths)
	w.lock.Lock()
	w.excludeDirPaths = slices.Compact(normalizedExcludeDirPaths)
	changed := w.updateActiveExcludeDirPathsLocked()
	w.lock.Unlock()
	if !changed {
		return nil
	}
	return w.watchRootDirPaths()
}

// updateActiveExcludeDirPathsLocked recomputes the excluded directories that
// apply and stops watching the directories that are now excluded.
//
// It returns whether the active excluded directories changed.
func (w *watcher) updateActiveExcludeDirPathsLocked() bool {
	var activeExcludeDirPaths []string
	for _, excludeDirPath := range w.excludeDirPaths {
		if !w.containsRootDirPathLocked(excludeDirPath) {
			activeExcludeDirPaths = append(activeExcludeDirPaths, excludeDirPath)
		}
	}
	if slices.Equal(w.activeExcludeDirPaths, activeExcludeDirPaths) {
		return false
	}
	w.activeExcludeDirPaths = activeExcludeDirPaths
	for dirPath := range w.recursiveDirPaths {
		if !w.isExcludedPathLocked(dirPath) {
			continue
		}
		delete(w.recursiveDirPaths, dirPath)
		w.removeUnneededWatchLocked(dirPath)
	}
	return true
}

func (w *watcher) containsRootDirPathLocked(normalizedDirPath string) bool {
	for rootDirPath := range w.rootDirPaths {
		if normalpath.EqualsOrContainsPath(normalizedDirPath, normalpath.Normalize(rootDirPath), normalpath.Absolute) {
			return true
		}
	}
	return false
}

// watchRootDirPaths watches all root directories again, as after the excluded
// directories change. Directories that are already watched are skipped.
func (w *watcher) watchRootDirPaths() error {
	w.lock.RLock()
	rootDirPaths := slices.Sorted(maps.Keys(w.rootDirPaths))
	w.lock.RUnlock()
	for _, rootDirPath := range rootDirPaths {
		if err := w.watchRootDirPath(rootDirPath, false); err != nil {
			return err
		}
	}
	return nil
}

// removeUnneededWatchLocked stops watching dirPath in fsnotifyWatcher if it is
// no longer needed for recursive watching, a watched file, or a pending directory.
func (w *watcher) removeUnneededWatchLocked(dirPath string) {
	if _, ok := w.recursiveDirPaths[dirPath]; ok {
		return
	}
	for filePath := range w.filePaths {
		if filepath.Dir(filePath) == dirPath {
			return
		}
	}
	for pendingDirPath := range w.pendingDirPaths {
		if filepath.Dir(pendingDirPath) == dirPath {
			return
		}
	}
	if err := w.fsnotifyWatcher.Remove(dirPath); err != nil && !errors.Is(err, fsnotify.ErrNonExistentWatch) {
		w.logger.Debug("cannot unwatch directory", slog.String("path", dirPath), xslog.ErrorAttr(err))
	}
}

func (w *watcher) Events() <-chan Event {
	return w.eventChan
}

func (w *watcher) Errors() <-chan error {
	return w.errChan
}

func (w *watcher) Close() error {
	w.closeOnce.Do(func() {
		close(w.closeChan)
		w.closeErr = w.fsnotifyWatcher.Close()
		w.waitGroup.Wait()
	})
	return w.closeErr
}

// run forwards events from the underlying fsnotify.Watcher until it is closed.
func (w *watcher) run() {
	defer close(w.eventChan)
	defer close(w.errChan)
	fsnotifyEventChan := w.fsnotifyWatcher.Events
	fsnotifyErrChan := w.fsnotifyWatcher.Errors
	for fsnotifyEventChan != nil || fsnotifyErrChan != nil {
		select {
		case fsnotifyEvent, ok := <-fsnotifyEventChan:
			if !ok {
				fsnotifyEventChan = nil
				continue
			}
			w.handleFsnotifyEvent(fsnotifyEvent)
		case err, ok := <-fsnotifyErrChan:
			if !ok {
				fsnotifyErrChan = nil
				continue
			}
			w.sendError(err)
		case <-w.closeChan:
			return
		}
	}
}

func (w *watcher) handleFsnotifyEvent(fsnotifyEvent fsnotify.Event) {
	if fsnotifyEvent.Op == fsnotify.Chmod {
		return
	}
	path := filepath.Clean(fsnotifyEvent.Name)
	if !w.isReportedPath(path) {
		return
	}
	switch {
	case fsnotifyEvent.Has(fsnotify.Create):
		w.handleCreate(path)
	case fsnotifyEvent.Has(fsnotify.Remove), fsnotifyEvent.Has(fsnotify.Rename):
		isRemovedDir, hadPendingDir := w.removeDirPathRecursive(path)
		if isRemovedDir {
			w.sendEvent(Event{Path: path, IsRemovedDir: true})
			if w.isRecursiveDirPath(filepath.Dir(path)) {
				if fileInfo, err := os.Stat(path); err == nil && fileInfo.IsDir() {
					w.handleCreate(path)
				}
			}
		} else if !hadPendingDir {
			w.sendEvent(Event{Path: path})
		}
		if isRemovedDir || hadPendingDir {
			w.watchMissingRootDirPaths()
		}
	default:
		if !w.isPendingParentDirPath(path) {
			w.sendEvent(Event{Path: path})
		}
	}
}

// handleCreate watches a directory created within a recursively watched directory
// or at a pending root directory path, and reports the files already within it.
func (w *watcher) handleCreate(path string) {
	if w.takePendingDirPath(path) {
		w.watchMissingRootDirPaths()
		if !w.isFilePath(path) {
			return
		}
	}
	fileInfo, err := os.Stat(path)
	if err != nil || !fileInfo.IsDir() {
		// A file, or a path that was removed again already.
		w.sendEvent(Event{Path: path})
		return
	}
	if !w.isRecursiveDirPath(filepath.Dir(path)) || strings.HasPrefix(filepath.Base(path), ".") {
		return
	}
	// The directory may be removed again before it is walked.
	if err := w.addDirPathRecursive(path); err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			w.sendError(err)
		}
		return
	}
	if err := w.sendEventsForFilesInDirPath(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
		w.sendError(err)
	}
}

func (w *watcher) sendEvent(event Event) {
	select {
	case w.eventChan <- event:
	case <-w.closeChan:
	}
}

func (w *watcher) sendError(err error) {
	select {
	case w.errChan <- err:
	case <-w.closeChan:
	}
}

func (w *watcher) isReportedPath(path string) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if _, ok := w.filePaths[path]; ok {
		return true
	}
	if _, ok := w.pendingDirPaths[path]; ok {
		return true
	}
	for pendingDirPath := range w.pendingDirPaths {
		if filepath.Dir(pendingDirPath) == path {
			return true
		}
	}
	if w.isExcludedPathLocked(path) {
		return false
	}
	if _, ok := w.recursiveDirPaths[filepath.Dir(path)]; ok {
		return true
	}
	// A recursively watched directory itself may be removed or renamed.
	_, ok := w.recursiveDirPaths[path]
	return ok
}

func (w *watcher) isFilePath(path string) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	_, ok := w.filePaths[path]
	return ok
}

func (w *watcher) isRecursiveDirPath(dirPath string) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	_, ok := w.recursiveDirPaths[dirPath]
	return ok
}

func (w *watcher) isPendingParentDirPath(dirPath string) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	for pendingDirPath := range w.pendingDirPaths {
		if filepath.Dir(pendingDirPath) == dirPath {
			return true
		}
	}
	return false
}

// watchRootDirPath watches absDirPath and its subdirectories if it exists, or
// watches its nearest existing ancestor for its creation if it does not.
//
// If sendEvents is true and absDirPath exists, an Event is sent for each file
// already within it.
func (w *watcher) watchRootDirPath(absDirPath string, sendEvents bool) error {
	for {
		err := w.addDirPathRecursive(absDirPath)
		if err == nil {
			if sendEvents {
				if err := w.sendEventsForFilesInDirPath(absDirPath); err != nil && !errors.Is(err, fs.ErrNotExist) {
					return err
				}
			}
			return nil
		}
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		watchPath := absDirPath
		for {
			parentDirPath := filepath.Dir(watchPath)
			if parentDirPath == watchPath {
				return err
			}
			addErr := w.addPendingDirPath(watchPath)
			if addErr == nil {
				break
			}
			if !errors.Is(addErr, fs.ErrNotExist) {
				return addErr
			}
			watchPath = parentDirPath
		}
		// If watchPath was created before its parent was watched, try again.
		if _, statErr := os.Stat(watchPath); statErr != nil {
			if errors.Is(statErr, fs.ErrNotExist) {
				return nil
			}
			return statErr
		}
		w.takePendingDirPath(watchPath)
	}
}

func (w *watcher) watchMissingRootDirPaths() {
	for _, rootDirPath := range w.missingRootDirPaths() {
		if err := w.watchRootDirPath(rootDirPath, true); err != nil {
			w.sendError(err)
		}
	}
}

func (w *watcher) missingRootDirPaths() []string {
	w.lock.RLock()
	defer w.lock.RUnlock()
	var dirPaths []string
	for rootDirPath := range w.rootDirPaths {
		if _, ok := w.recursiveDirPaths[rootDirPath]; !ok && !w.isExcludedPathLocked(rootDirPath) {
			dirPaths = append(dirPaths, rootDirPath)
		}
	}
	slices.Sort(dirPaths)
	return dirPaths
}

func (w *watcher) addPendingDirPath(dirPath string) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if err := w.fsnotifyWatcher.Add(filepath.Dir(dirPath)); err != nil {
		return fmt.Errorf("watch %s: %w", dirPath, err)
	}
	w.pendingDirPaths[dirPath] = struct{}{}
	return nil
}

func (w *watcher) takePendingDirPath(dirPath string) bool {
	w.lock.Lock()
	defer w.lock.Unlock()
	if _, ok := w.pendingDirPaths[dirPath]; !ok {
		return false
	}
	delete(w.pendingDirPaths, dirPath)
	w.removeUnneededWatchLocked(filepath.Dir(dirPath))
	return true
}

// addDirPathRecursive watches absDirPath and all of its subdirectories.
// The lock is held per directory so that events keep being handled during the walk.
func (w *watcher) addDirPathRecursive(absDirPath string) error {
	return filepath.WalkDir(
		absDirPath,
		func(path string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				if path == absDirPath {
					return fmt.Errorf("watch %s: %w", absDirPath, err)
				}
				if !errors.Is(err, fs.ErrNotExist) {
					w.logger.Warn(
						"skipping unreadable directory",
						slog.String("path", path),
						xslog.ErrorAttr(err),
					)
				}
				return filepath.SkipDir
			}
			if !dirEntry.IsDir() {
				return nil
			}
			if path != absDirPath && strings.HasPrefix(dirEntry.Name(), ".") {
				return filepath.SkipDir
			}
			if err := w.addDirPath(path); err != nil {
				if path == absDirPath || errors.Is(err, filepath.SkipDir) {
					return err
				}
				if !errors.Is(err, fs.ErrNotExist) {
					w.logger.Warn(
						"skipping unreadable directory",
						slog.String("path", path),
						xslog.ErrorAttr(err),
					)
				}
				return filepath.SkipDir
			}
			return nil
		},
	)
}

// addDirPath watches the single directory absDirPath.
//
// It returns filepath.SkipDir if the directory is excluded.
func (w *watcher) addDirPath(absDirPath string) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	if w.isExcludedPathLocked(absDirPath) {
		return filepath.SkipDir
	}
	if _, ok := w.recursiveDirPaths[absDirPath]; ok {
		return nil
	}
	if err := w.fsnotifyWatcher.Add(absDirPath); err != nil {
		return fmt.Errorf("watch %s: %w", absDirPath, err)
	}
	w.recursiveDirPaths[absDirPath] = struct{}{}
	return nil
}

// removeDirPathRecursive stops watching absDirPath and all of its subdirectories.
//
// It returns whether absDirPath was a recursively watched directory and whether
// any pending directory watch under absDirPath was removed.
func (w *watcher) removeDirPathRecursive(absDirPath string) (bool, bool) {
	normalizedDirPath := normalpath.Normalize(absDirPath)
	w.lock.Lock()
	defer w.lock.Unlock()
	_, isRemovedDir := w.recursiveDirPaths[absDirPath]
	if isRemovedDir {
		for dirPath := range w.recursiveDirPaths {
			if !normalpath.EqualsOrContainsPath(normalizedDirPath, normalpath.Normalize(dirPath), normalpath.Absolute) {
				continue
			}
			delete(w.recursiveDirPaths, dirPath)
			w.removeUnneededWatchLocked(dirPath)
		}
	}
	hadPendingDir := false
	for pendingDirPath := range w.pendingDirPaths {
		if !normalpath.EqualsOrContainsPath(normalizedDirPath, normalpath.Normalize(filepath.Dir(pendingDirPath)), normalpath.Absolute) {
			continue
		}
		delete(w.pendingDirPaths, pendingDirPath)
		w.removeUnneededWatchLocked(filepath.Dir(pendingDirPath))
		hadPendingDir = true
	}
	return isRemovedDir, hadPendingDir
}

// sendEventsForFilesInDirPath sends an Event for every file within absDirPath
// and its watched subdirectories.
func (w *watcher) sendEventsForFilesInDirPath(absDirPath string) error {
	return filepath.WalkDir(
		absDirPath,
		func(path string, dirEntry fs.DirEntry, err error) error {
			if err != nil {
				if path != absDirPath && errors.Is(err, fs.ErrNotExist) {
					return nil
				}
				return err
			}
			if dirEntry.IsDir() {
				if !w.isRecursiveDirPath(path) {
					return filepath.SkipDir
				}
				return nil
			}
			w.sendEvent(Event{Path: path})
			return nil
		},
	)
}

func (w *watcher) isExcludedPathLocked(path string) bool {
	if len(w.activeExcludeDirPaths) == 0 {
		return false
	}
	normalizedPath := normalpath.Normalize(path)
	for _, normalizedExcludeDirPath := range w.activeExcludeDirPaths {
		if normalpath.EqualsOrContainsPath(normalizedExcludeDirPath, normalizedPath, normalpath.Absolute) {
			return true
		}
	}
	return false
}
