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
	"log/slog"
)

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
	// created. Hidden subdirectories are skipped and symbolic links are not followed.
	AddDirPath(dirPath string) error
	// AddFilePath watches the single file, which does not need to exist yet.
	// Adding it again re-establishes the watch if its directory was recreated.
	AddFilePath(filePath string) error
	// SetExcludeDirPaths replaces the set of directories within which changes
	// are not reported. Files added with AddFilePath are reported regardless.
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
