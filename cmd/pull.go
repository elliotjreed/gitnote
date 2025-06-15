package cmd

import (
	"fmt"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"gitnote/internal/git"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from remote repository",
	Long:  "Pull changes from remote repository with merge conflict handling",
	RunE:  runPull,
}

func runPull(cmd *cobra.Command, args []string) error {
	gitManager := git.NewManager(".")
	
	if !gitManager.IsGitRepo() {
		return fmt.Errorf("current directory is not a git repository")
	}
	
	output, err := gitManager.Pull()
	if err != nil {
		hasConflicts, conflictErr := gitManager.HasMergeConflicts()
		if conflictErr != nil {
			return fmt.Errorf("failed to check merge conflicts: %w", conflictErr)
		}
		
		if hasConflicts {
			fmt.Printf("Pull failed with merge conflicts:\n%s\n", output)
			return handleMergeConflicts(gitManager)
		}
		
		return fmt.Errorf("pull failed: %w\nOutput: %s", err, output)
	}
	
	fmt.Printf("Pull completed successfully:\n%s", output)
	
	if strings.Contains(output, "files changed") {
		status, err := gitManager.GetStatus()
		if err == nil && len(status) == 0 {
			fmt.Println("\nRepository is up to date.")
		}
	}
	
	return nil
}

func handleMergeConflicts(gitManager *git.Manager) error {
	fmt.Println("\nMerge conflicts detected. Choose an option:")
	
	prompt := promptui.Select{
		Label: "Select action",
		Items: []string{
			"Manually fix conflicts (you will need to resolve conflicts and commit)",
			"Roll back (discard changes and reset to HEAD)",
		},
	}
	
	_, result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("failed to get user choice: %w", err)
	}
	
	switch result {
	case "Roll back (discard changes and reset to HEAD)":
		if err := gitManager.Reset(); err != nil {
			return fmt.Errorf("failed to reset repository: %w", err)
		}
		fmt.Println("Repository has been reset to previous state")
		return nil
	default:
		fmt.Println("Please resolve the merge conflicts manually and then run 'git commit' to complete the merge")
		return nil
	}
}