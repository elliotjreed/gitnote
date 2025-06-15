package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"gitnote/internal/git"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit new or updated files to git",
	Long:  "Add and commit any new or updated files with a descriptive commit message",
	RunE:  runCommit,
}

func runCommit(cmd *cobra.Command, args []string) error {
	gitManager := git.NewManager(".")
	
	if !gitManager.IsGitRepo() {
		return fmt.Errorf("current directory is not a git repository")
	}
	
	status, err := gitManager.GetStatus()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}
	
	if len(status) == 0 {
		fmt.Println("No changes to commit")
		return nil
	}
	
	var newFiles []string
	var modifiedFiles []string
	var filesToAdd []string
	
	for _, line := range status {
		if len(line) < 3 {
			continue
		}
		
		statusCode := line[:2]
		filename := strings.TrimSpace(line[3:])
		
		switch statusCode {
		case "??":
			newFiles = append(newFiles, filename)
			filesToAdd = append(filesToAdd, filename)
		case " M", "MM":
			modifiedFiles = append(modifiedFiles, filename)
			filesToAdd = append(filesToAdd, filename)
		case "M ", "A ":
		default:
			filesToAdd = append(filesToAdd, filename)
		}
	}
	
	if len(filesToAdd) > 0 {
		if err := gitManager.AddFiles(filesToAdd); err != nil {
			return fmt.Errorf("failed to add files: %w", err)
		}
	}
	
	commitMessage := buildCommitMessage(newFiles, modifiedFiles)
	
	if err := gitManager.Commit(commitMessage); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	
	fmt.Printf("Committed changes with message: %s\n", commitMessage)
	return nil
}

func buildCommitMessage(newFiles, modifiedFiles []string) string {
	var messageParts []string
	
	if len(newFiles) > 0 {
		if len(newFiles) == 1 {
			messageParts = append(messageParts, fmt.Sprintf("Add %s", newFiles[0]))
		} else {
			messageParts = append(messageParts, fmt.Sprintf("Add %d new files", len(newFiles)))
		}
	}
	
	if len(modifiedFiles) > 0 {
		if len(modifiedFiles) == 1 {
			messageParts = append(messageParts, fmt.Sprintf("Update %s", modifiedFiles[0]))
		} else {
			messageParts = append(messageParts, fmt.Sprintf("Update %d files", len(modifiedFiles)))
		}
	}
	
	if len(messageParts) == 0 {
		return "Update files"
	}
	
	return strings.Join(messageParts, " and ")
}