package fw

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewFileWatcher tests the NewFileWatcher function.
func TestNewFileWatcher(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filewatcher_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fw, err := NewFileWatcher(tempDir)
	require.NoError(t, err)
	defer fw.Close()

	assert.NotNil(t, fw.watcher)
	assert.Equal(t, tempDir, fw.rootPath)
	assert.NotNil(t, fw.eventCh)
	assert.Equal(t, DefaultEventBufferSize, cap(fw.eventCh))
}

// TestFileCreation tests the file creation event.
func TestFileCreation(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filewatcher_test_*")
	fmt.Println("tempDir", tempDir)
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fw, err := NewFileWatcher(tempDir)
	require.NoError(t, err)
	defer fw.Close()

	eventCount := 0
	done := make(chan bool)

	go func() {
		// Wait for up to 2 events over 1 second to see if we get create+modify
		timeout := time.After(1 * time.Second)
		for {
			select {
			case event := <-fw.Events():
				eventCount++
				assert.Equal(t, filepath.Join(tempDir, "test.txt"), event.Path)
				assert.False(t, event.IsDir)
				if eventCount == 1 {
					assert.Equal(t, FileModified, event.Type)
				}
			case <-timeout:
				done <- true
				return
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	<-done
	assert.Equal(t, 1, eventCount, "Should receive exactly one event (create) when writing new file")
}

// TestFileDeletion tests the file deletion event.
func TestFileDeletion(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filewatcher_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	fw, err := NewFileWatcher(tempDir)
	require.NoError(t, err)
	defer fw.Close()

	done := make(chan bool)

	go func() {
		event := <-fw.Events()
		assert.Equal(t, FileDeleted, event.Type)
		assert.Equal(t, testFile, event.Path)
		done <- true
	}()

	time.Sleep(100 * time.Millisecond)
	err = os.Remove(testFile)
	require.NoError(t, err)

	select {
	case <-done:
		// Test passed
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for file deletion event")
	}
}

// TestIgnoredPaths tests the ignored paths feature.
func TestIgnoredPaths(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filewatcher_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fw, err := NewFileWatcher(tempDir)
	require.NoError(t, err)
	defer fw.Close()

	ignoredDir := filepath.Join(tempDir, ".git")
	fw.ignorePaths[ignoredDir] = true

	err = os.Mkdir(ignoredDir, 0755)
	require.NoError(t, err)

	timeout := make(chan bool)
	go func() {
		time.Sleep(500 * time.Millisecond)
		timeout <- true
	}()

	ignoredFile := filepath.Join(ignoredDir, "test.txt")
	err = os.WriteFile(ignoredFile, []byte("test content"), 0644)
	require.NoError(t, err)

	select {
	case event := <-fw.Events():
		t.Fatalf("Received unexpected event: %+v", event)
	case <-timeout:
		// Test passed - no events received
	}
}

// TestEventBatching tests that rapid successive events are properly batched.
// It verifies that multiple quick modifications to the same file result in a single event.
func TestEventBatching(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "filewatcher_test_*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	fw, err := NewFileWatcher(tempDir)
	require.NoError(t, err)
	defer fw.Close()

	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("initial content"), 0644)
	require.NoError(t, err)

	eventCount := make(chan int)

	go func() {
		count := 0
		timeout := time.After(DefaultFlushInterval * 2)
		for {
			select {
			case <-fw.Events():
				count++
			case <-timeout:
				eventCount <- count
				return
			}
		}
	}()

	for i := 0; i < 5; i++ {
		err = os.WriteFile(testFile, []byte("modified content"), 0644)
		require.NoError(t, err)
		time.Sleep(10 * time.Millisecond)
	}

	// Verify we receive only one batched event
	count := <-eventCount
	assert.Equal(t, 1, count, "Expected one batched event, got %d events", count)
}
