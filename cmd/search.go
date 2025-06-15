package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"gitnote/internal/note"
)

var (
	searchFull bool
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search for notes by title or content",
	Long:  "Search for notes by title, or optionally by content with --full flag",
	Args:  cobra.ExactArgs(1),
	RunE:  runSearch,
}

func init() {
	searchCmd.Flags().BoolVar(&searchFull, "full", false, "Search in file content as well as titles")
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]
	noteManager := note.NewManager(".")
	
	results, err := noteManager.SearchNotes(query, searchFull)
	if err != nil {
		return fmt.Errorf("failed to search notes: %w", err)
	}
	
	if len(results) == 0 {
		fmt.Println("No notes found matching the query")
		return nil
	}
	
	for _, note := range results {
		fmt.Println(note.Path)
	}
	
	return nil
}