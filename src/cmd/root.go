package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

var rootCmd = cobra.Command{
	Use:   "gep",
	Short: "A simple grep implementation in Go",
	Long:  "A simple grep implementation in Go that supports basic patterns",
}

func StartCommand() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing  command: %v", err)
	}
}
