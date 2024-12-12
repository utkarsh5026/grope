package fw

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/codecrafters-io/grep-starter-go/src/logs"
	"github.com/fsnotify/fsnotify"
)

const (
	DefaultEventBufferSize = 100
	DefaultFlushInterval   = 100 * time.Millisecond
)

type EventType int

const (
	FileCreated EventType = iota
	FileModified
	FileDeleted
)

func (e EventType) String() string {
	return []string{"File Created", "FileModified", "FileDeleted"}[e]
}

type FileEvent struct {
	Type       EventType
	Path       string
	ModTime    time.Time
	IsDir      bool
	ChangeType string
}

type FileWatcher struct {
	watcher     *fsnotify.Watcher // The underlying fsnotify watcher.
	rootPath    string            // The root path that the watcher is watching.
	eventCh     chan FileEvent    // The channel to send file events to.
	ignorePaths map[string]bool   // A map of paths to ignore.
}

// NewFileWatcher creates a new FileWatcher instance.
// It initializes the fsnotify watcher, sets up the event channel, and starts the watching process.
// The function returns an error if the watcher creation fails.
func NewFileWatcher(rootPath string) (*FileWatcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	fw := &FileWatcher{
		watcher:     watcher,
		rootPath:    rootPath,
		eventCh:     make(chan FileEvent, DefaultEventBufferSize),
		ignorePaths: make(map[string]bool),
	}

	if err := fw.watchRecursively(rootPath); err != nil {
		watcher.Close()
		return nil, fmt.Errorf("failed to watch recursively: %w", err)
	}

	go fw.start()

	return fw, nil
}

// watchRecursively watches the given root path recursively and adds all directories to the watcher.
// It also adds all files to the watcher if they are not ignored.
func (fw *FileWatcher) watchRecursively(rootPath string) error {
	return filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && !fw.ignorePaths[path] {
			fw.watcher.Add(path)
		}

		return nil
	})
}

// shouldIgnore returns true if the given path should be ignored.
// eg. .gitignore, .DS_Store, etc.
func (fw *FileWatcher) shouldIgnore(path string) bool {
	for ignore := range fw.ignorePaths {
		if strings.HasPrefix(path, ignore) {
			return true
		}
	}
	return false
}

// start begins the file watching process. It runs in a separate goroutine and handles
// incoming file system events from the watcher. Events are batched together over a short
// interval (DefaultFlushInterval) to prevent event flooding.
//
// The function implements debouncing by maintaining a map of pending events and using
// a timer to flush them periodically. This helps coalesce rapid sequences of events
// for the same file into a single event.
//
// Events are processed as follows:
// - Create events generate FileCreated events
// - Write events generate FileModified events
// - Remove events generate FileDeleted events
//
// Any errors from the underlying watcher are logged but do not stop the watching process.
// The function will return if either the Events or Errors channel is closed.
func (fw *FileWatcher) start() {
	var timer *time.Timer

	pendingEvents := make(map[string]FileEvent)

	for {
		select {
		case event, ok := <-fw.watcher.Events:

			if !ok {
				return
			}

			if fw.shouldIgnore(event.Name) {
				continue
			}

			fe := FileEvent{
				Path:    event.Name,
				ModTime: time.Now(),
			}

			switch {
			case event.Op&fsnotify.Create == fsnotify.Create:
				fe.Type = FileCreated
			case event.Op&fsnotify.Write == fsnotify.Write:
				fe.Type = FileModified
			case event.Op&fsnotify.Remove == fsnotify.Remove:
				fe.Type = FileDeleted
			}

			pendingEvents[event.Name] = fe

			if timer != nil {
				timer.Stop()
			}

			timer = time.AfterFunc(DefaultFlushInterval, func() {
				fw.flushEvents(pendingEvents)
				pendingEvents = make(map[string]FileEvent)
			})

		case err, ok := <-fw.watcher.Errors:
			if !ok {
				return
			}
			logs.Error("file watcher error: %v", err)
		}

	}
}

// flushEvents flushes the pending events to the event channel.
func (fw *FileWatcher) flushEvents(events map[string]FileEvent) {
	for _, event := range events {
		event.IsDir = isDir(event.Path)
		fw.eventCh <- event
	}
}

// isDir returns true if the given path is a directory.
func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// Events returns a channel that emits FileEvents.
func (fw *FileWatcher) Events() <-chan FileEvent {
	return fw.eventCh
}

// Close closes the file watcher and releases all resources.
func (fw *FileWatcher) Close() error {
	return fw.watcher.Close()
}

// StartWatching starts watching the given root path and calls the given function for each file event.
// It blocks indefinitely while watching for file events.
func StartWatching(rootPath string, onEvent func(FileEvent) error) error {
	fw, err := NewFileWatcher(rootPath)
	if err != nil {
		return err
	}
	defer fw.Close()

	// Block indefinitely while processing events
	for event := range fw.Events() {
		if err := onEvent(event); err != nil {
			return err
		}
	}
	return nil
}
