package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/codecrafters-io/grep-starter-go/src/matcher"
)

type File struct {
	Name    string
	Size    string
	ModTime string
	Perms   string
	Path    string
}

// fileFromInfo creates a File struct from fs.FileInfo
func fileFromInfo(basePath, currentFilePath string, info fs.FileInfo) File {
	relPath, _ := filepath.Rel(basePath, currentFilePath)

	return File{
		Name:    info.Name(),
		Size:    formatSize(info.Size()),
		ModTime: info.ModTime().Format("Jan 02 15:04"),
		Perms:   info.Mode().String(),
		Path:    relPath,
	}
}

// SearchFilesInDir searches for files matching the pattern in the given directory path.
// It returns a slice of matching File structs and any error encountered.
// This function does not search subdirectories.
func SearchFilesInDir(searchPath string, pattern string) ([]File, error) {
	files, err := os.ReadDir(searchPath)
	absPath, _ := filepath.Abs(searchPath)

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

		foundFiles = append(foundFiles, fileFromInfo(absPath, filepath.Join(absPath, file.Name()), info))
	}

	return foundFiles, nil
}

// SearchDirRecursively searches for files matching the pattern in the given directory path
// and all its subdirectories recursively. It returns a slice of matching File structs
// and any error encountered.
func SearchDirRecursively(searchPath, pattern string, maxDepth int) ([]File, error) {
	var foundFiles []File
	basePath, err := filepath.Abs(searchPath)
	if err != nil {
		return nil, fmt.Errorf("error getting absolute path: %v", err)
	}

	err = filepath.WalkDir(searchPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		fullPath := filepath.Join(basePath, path)
		if !isWithinDepth(basePath, fullPath, maxDepth) {
			return nil
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

		foundFiles = append(foundFiles, fileFromInfo(basePath, fullPath, info))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	return foundFiles, nil
}

// formatSize converts a size in bytes to a human-readable string with appropriate units (B, KB, MB, etc)
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

// isWithinDepth checks if the current path is within the maximum depth from the base path.
// If maxDepth is 0, it returns true, meaning no depth limit is enforced.
func isWithinDepth(basePath, currentPath string, maxDepth int) bool {
	if maxDepth == 0 {
		return true
	}

	baseDepth, err := calcDepth(basePath)
	if err != nil {
		return false
	}

	currentDepth, err := calcDepth(currentPath)
	if err != nil {
		return false
	}

	return currentDepth-baseDepth <= maxDepth
}

// calcDepth calculates the depth of a path by counting the number of directories in the path
// starting from the root directory.
func calcDepth(path string) (int, error) {
	cleaned, err := filepath.Abs(path)
	if err != nil {
		return 0, fmt.Errorf("error getting absolute path: %v", err)
	}
	return len(strings.Split(cleaned, string(os.PathSeparator))) - 1, nil
}

// SortByDepth sorts the files by their depth in the directory tree.
// If two files have the same depth, they are sorted by their path.
func SortByDepth(files []File) []File {
	sort.Slice(files, func(i, j int) bool {
		di, _ := calcDepth(files[i].Path)
		dj, _ := calcDepth(files[j].Path)

		if di == dj {
			return files[i].Path < files[j].Path
		}

		return di < dj
	})
	return files
}
