package cmd

import (
	"github.com/codecrafters-io/grep-starter-go/src/file"
	"github.com/codecrafters-io/grep-starter-go/src/logs"
	"github.com/codecrafters-io/grep-starter-go/src/table"
	"github.com/spf13/cobra"
)

const (
	CurrentDir string = "."
)

var (
	searchPath string
	recursive  bool
	depth      int
)

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

		pattern := args[0]

		var files []file.File
		var err error
		if recursive {
			files, err = file.SearchDirRecursively(searchPath, pattern, depth)
		} else {
			files, err = file.SearchFilesInDir(searchPath, pattern)
		}

		if err != nil {
			logs.Fatal(err.Error())
		}

		files = file.SortByDepth(files)
		err = table.PrintTable(files, table.Options{
			Centered: true,
			Border:   true,
		})

		if err != nil {
			logs.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(filesCmd)
	filesCmd.Flags().StringVarP(&searchPath, "path", "p", ".", "The path to search for files")
	filesCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Search recursively")
	filesCmd.Flags().IntVarP(&depth, "depth", "d", 0, "Search recursively up to a certain depth (0 means unlimited)")
}
