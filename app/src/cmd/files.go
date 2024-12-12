package cmd

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/grep-starter-go/src/file"
	"github.com/codecrafters-io/grep-starter-go/src/logs"
	"github.com/codecrafters-io/grep-starter-go/src/table"
	"github.com/spf13/cobra"
)

const (
	CurrentDir string = "."
)

var (
	searchPath     string
	recursive      bool
	depth          int
	invert         bool
	caseSensitive  bool
	hidden         bool
	maxSize        int64
	minSize        int64
	modifiedAfter  string
	modifiedBefore string
)

func parseTime(timeStr string) (time.Time, error) {
	t, err := time.Parse(time.DateOnly, timeStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %v", err)
	}
	return t, nil
}

var filesCmd = &cobra.Command{
	Use:   "ls",
	Short: "Search for files in a directory",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logs.Fatal("No pattern provided to search for do something like 'gep ls -p . -r \"*.go\"'")
		}

		if searchPath == "" {
			searchPath = CurrentDir
		}

		ma, err := parseTime(modifiedAfter)
		if err != nil {
			logs.Fatal(err.Error())
		}

		mb, err := parseTime(modifiedBefore)
		if err != nil {
			logs.Fatal(err.Error())
		}

		pattern := args[0]
		options := file.SearchOptions{
			Recursive: recursive,
			Invert:    invert,
			MaxDepth:  depth,
			FileFilter: file.SearchWithFileProperty{
				CaseSensitive:  caseSensitive,
				Hidden:         hidden,
				MaxSize:        maxSize,
				MinSize:        minSize,
				ModifiedAfter:  ma,
				ModifiedBefore: mb,
			},
		}

		files, err := file.SearchWithPattern(searchPath, pattern, options)

		if err != nil {
			logs.Fatal(err.Error())
		}

		files = file.SortByDepth(files)
		if err := table.PrintTable(files, table.Options{
			Centered: true,
			Border:   true,
		}); err != nil {
			logs.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.Flags().StringVarP(&searchPath, "path", "p", ".", "The path to search for files")
	filesCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Search recursively")
	filesCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Search recursively up to a certain depth (0 means unlimited)")
	filesCmd.Flags().BoolVarP(&invert, "invert", "i", false, "Invert the search so it matches files that don't match the pattern")
	filesCmd.Flags().BoolVarP(&caseSensitive, "case-sensitive", "c", false, "Case sensitive search (default is case insensitive)")
	filesCmd.Flags().BoolVarP(&hidden, "hidden", "h", false, "Include hidden files in the search")
	filesCmd.Flags().Int64VarP(&maxSize, "max-size", "s", 0, "Maximum file size to search for")
	filesCmd.Flags().Int64VarP(&minSize, "min-size", "m", 0, "Minimum file size to search for")
	filesCmd.Flags().StringVarP(&modifiedAfter, "modified-after", "a", "", "Search for files modified after a certain date")
	filesCmd.Flags().StringVarP(&modifiedBefore, "modified-before", "b", "", "Search for files modified before a certain date")
}
