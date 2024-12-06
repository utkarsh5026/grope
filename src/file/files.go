package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/codecrafters-io/grep-starter-go/src/matcher"
)

type File struct {
	Name    string
	Size    string
	ModTime string
	Perms   string
}

// fileFromInfo creates a File struct from fs.FileInfo
func fileFromInfo(info fs.FileInfo) File {
	return File{
		Name:    info.Name(),
		Size:    formatSize(info.Size()),
		ModTime: info.ModTime().Format("Jan 02 15:04"),
		Perms:   info.Mode().String(),
	}
}

// SearchFilesInDir searches for files matching the pattern in the given directory path.
// It returns a slice of matching File structs and any error encountered.
// This function does not search subdirectories.
func SearchFilesInDir(path string, pattern string) ([]File, error) {
	files, err := os.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	var foundFiles []File
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if !matcher.Match([]byte(file.Name()), pattern) {
			continue
		}

		info, err := file.Info()
		if err != nil {
			return nil, fmt.Errorf("error getting file info: %v", err)
		}

		foundFiles = append(foundFiles, fileFromInfo(info))
	}

	return foundFiles, nil
}

// SearchDirRecursively searches for files matching the pattern in the given directory path
// and all its subdirectories recursively. It returns a slice of matching File structs
// and any error encountered.
func SearchDirRecursively(path string, pattern string) ([]File, error) {
	var foundFiles []File
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !matcher.Match([]byte(d.Name()), pattern) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		foundFiles = append(foundFiles, fileFromInfo(info))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	return foundFiles, nil
}

// formatSize converts a size in bytes to a human readable string with appropriate units (B, KB, MB, etc)
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
