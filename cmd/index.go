package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitnote/internal/index"
)

var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Generate or update the readme.md table of contents",
	Long:  "Scan the directory for markdown notes and create/update readme.md with a table of contents",
	RunE:  runIndex,
}

func runIndex(cmd *cobra.Command, args []string) error {
	generator := index.NewGenerator(".")
	
	upToDate, err := generator.IsReadmeUpToDate()
	if err != nil {
		return fmt.Errorf("failed to check readme status: %w", err)
	}
	
	if upToDate {
		fmt.Println("readme.md is already up to date")
		return nil
	}
	
	if err := generator.GenerateReadme(); err != nil {
		return fmt.Errorf("failed to generate readme: %w", err)
	}
	
	fmt.Println("readme.md has been updated")
	return nil
}