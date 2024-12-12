package cmd

import (
	"github.com/codecrafters-io/grep-starter-go/src/logs"
	"github.com/spf13/cobra"
)

var rootCmd = cobra.Command{
	Use:   "gep",
	Short: "A simple grep implementation in Go",
	Long:  "A simple grep implementation in Go that supports basic patterns",
}

func StartCommand() {
	if err := rootCmd.Execute(); err != nil {
		logs.Fatal(err.Error())
	}
}
