package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitnote",
	Short: "A CLI tool for managing markdown notes in a Git repository",
	Long: `GitNote is a command-line tool that helps you organise and manage
markdown notes in a Git repository with automatic naming conventions
and directory-based organisation.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(indexCmd)
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(pullCmd)
}
