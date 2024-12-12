package fileutils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/grep-starter-go/src/logs"
)

// CloseFile closes a file and logs any errors that occur.
func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		logs.Error("error closing file: %v", err)
	}
}

// ReadFileContent reads the content of a file and returns it as a byte slice.
// It returns an error if the file cannot be read.
func ReadFileContent(path string) ([]byte, error) {
	if !CheckFileExists(path) {
		return nil, fmt.Errorf("file does not exist: %s", path)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}
	return content, nil
}

// CheckFileExists checks if a file exists at the given path.
// It returns true if the file exists, false otherwise.
func CheckFileExists(path string) bool {
	_, err := os.Stat(filepath.Clean(path))
	return err == nil
}

// GetFileExtension returns the file extension (including the dot).
// Returns empty string if there is no extension.
func GetFileExtension(path string) string {
	return filepath.Ext(path)
}

// GetFileSize returns the size of a file in bytes.
// It returns an error if the file cannot be read.
func GetFileSize(path string) (int64, error) {
	info, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return 0, fmt.Errorf("error getting file size: %v", err)
	}
	return info.Size(), nil
}
