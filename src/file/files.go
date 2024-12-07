package file

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/codecrafters-io/grep-starter-go/src/matcher"
)

type File struct {
	Name    string
	Size    string
	ModTime string
	Perms   string
	Path    string
	AbsPath string
}

type SearchOptions struct {
	Recursive  bool
	Invert     bool
	MaxDepth   int
	FileFilter SearchWithFileProperty
}

type SearchWithFileProperty struct {
	CaseSensitive  bool
	Hidden         bool
	MaxSize        int64
	MinSize        int64
	ModifiedAfter  time.Time
	ModifiedBefore time.Time
}

// FromInfo creates a File struct from fs.FileInfo
func FromInfo(basePath, currentFilePath string, info fs.FileInfo) File {
	relPath, _ := filepath.Rel(basePath, currentFilePath)

	return File{
		Name:    info.Name(),
		Size:    formatSize(info.Size()),
		ModTime: info.ModTime().Format("Jan 02 15:04"),
		Perms:   info.Mode().String(),
		Path:    relPath,
	}
}

// SearchWithPattern searches for files matching the pattern in the given directory path
// and returns a slice of matching File structs.
func SearchWithPattern(searchPath, pattern string, options SearchOptions) ([]File, error) {
	if !options.Recursive {
		options.MaxDepth = 1
	}
	return searchFiles(searchPath, pattern, options)
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

// searchFiles searches for files matching the pattern in the given directory path
// and returns a slice of matching File structs.
//
// Parameters:
//   - searchPath: The directory path to search in
//   - pattern: The pattern to match files against
//   - filter: A function that determines whether a file matches the pattern
//   - options: The search options
//
// Returns:
//   - []File: A slice of matching File structs
//   - error: An error if something goes wrong
func searchFiles(searchPath, pattern string, options SearchOptions) ([]File, error) {
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
		if !isWithinDepth(basePath, fullPath, options.MaxDepth) {
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if filterFile(d, pattern, options.Invert, options.FileFilter) {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		foundFiles = append(foundFiles, FromInfo(basePath, fullPath, info))
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking directory: %v", err)
	}

	return foundFiles, nil
}

// filterFile filters a file based on the given options and returns true if the file matches the pattern.
// If invert is true, the function returns true if the file does not match the pattern.
func filterFile(file os.DirEntry, pattern string, invert bool, options SearchWithFileProperty) bool {
	info, err := file.Info()
	if err != nil {
		return false
	}

	if options.MaxSize > 0 && info.Size() > options.MaxSize {
		return false
	}

	if options.MinSize > 0 && info.Size() < options.MinSize {
		return false
	}

	modTime := info.ModTime()
	if !options.ModifiedBefore.IsZero() && modTime.After(options.ModifiedBefore) {
		return false
	}

	if !options.ModifiedAfter.IsZero() && modTime.Before(options.ModifiedAfter) {
		return false
	}

	fileName := file.Name()
	if !options.Hidden && strings.HasPrefix(fileName, ".") {
		return false
	}

	if !options.CaseSensitive {
		pattern = strings.ToLower(pattern)
		fileName = strings.ToLower(fileName)
	}

	match := matcher.Match([]byte(fileName), pattern)
	return invert != match
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
