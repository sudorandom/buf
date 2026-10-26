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

// Package filewatch provides notifications of changes to files on the local file system.
package filewatch

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"slices"
	"time"

	"buf.build/go/standard/xlog/xslog"
	"github.com/fsnotify/fsnotify"
)

// ErrEventOverflow is sent on Errors when changes were dropped because they
// were not read quickly enough, so callers should assume any watched file changed.
var ErrEventOverflow = fsnotify.ErrEventOverflow

// Event is a change to a watched file, or to a file or directory within a watched directory.
type Event struct {
	// Path is the absolute path of the file or directory that changed.
	Path string
	// IsRemovedDir is true if Path was a watched directory that was removed or renamed.
	IsRemovedDir bool
}

// Watcher watches files and directories for changes.
//
// Metadata-only changes are not reported. A change may be reported more than
// once, or be missed right after a file is created, so callers should coalesce
// Events and read the files afterwards.
type Watcher interface {
	// AddDirPath watches the directory and its subdirectories, including those
	// created later. A directory that does not exist yet is watched once it is
	// created. Hidden subdirectories are skipped and symbolic links within it
	// are not followed. If the directory itself is a symbolic link, it is
	// resolved and changes are reported with the resolved path.
	AddDirPath(dirPath string) error
	// AddFilePath watches the single file, which does not need to exist yet.
	// Adding it again re-establishes the watch if its directory was recreated.
	AddFilePath(filePath string) error
	// SetExcludeDirPaths replaces the set of directories within which changes
	// are not reported. A directory that contains a watched directory is not
	// excluded. Files added with AddFilePath are reported regardless.
	SetExcludeDirPaths(dirPaths ...string) error
	// Events returns the channel on which Events are sent. It is closed by Close.
	Events() <-chan Event
	// Errors returns the channel on which errors encountered while watching are sent.
	// It is closed by Close.
	Errors() <-chan error
	// Close stops watching and releases all resources.
	Close() error
}

// NewWatcher returns a new Watcher, which must be closed when no longer needed.
func NewWatcher(logger *slog.Logger) (Watcher, error) {
	return newWatcher(logger)
}

// Run calls onChanges with the paths that changed, until ctx is done or the
// Watcher is closed.
//
// Events for which isRelevant returns false are ignored. onChanges is called
// with the sorted paths of the relevant Events once no relevant Event has been
// received for debounceDuration. The paths are empty if changes were dropped,
// in which case any watched file may have changed. Errors are logged.
func Run(
	ctx context.Context,
	logger *slog.Logger,
	watcher Watcher,
	debounceDuration time.Duration,
	isRelevant func(Event) bool,
	onChanges func(context.Context, []string),
) {
	timer := time.NewTimer(debounceDuration)
	timer.Stop()
	changedPaths := make(map[string]struct{})
	for {
		select {
		case <-ctx.Done():
			return
		case event, ok := <-watcher.Events():
			if !ok {
				return
			}
			if !isRelevant(event) {
				continue
			}
			changedPaths[event.Path] = struct{}{}
			timer.Reset(debounceDuration)
		case err, ok := <-watcher.Errors():
			if !ok {
				return
			}
			logger.WarnContext(ctx, "file watcher error", xslog.ErrorAttr(err))
			if errors.Is(err, ErrEventOverflow) {
				timer.Reset(debounceDuration)
			}
		case <-timer.C:
			onChanges(ctx, slices.Sorted(maps.Keys(changedPaths)))
			clear(changedPaths)
		}
	}
}
