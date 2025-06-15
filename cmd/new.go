package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"

	"gitnote/internal/note"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new note",
	Long:  "Create a new note with category selection and automatic date formatting",
	RunE:  runNew,
}

func runNew(cmd *cobra.Command, args []string) error {
	noteManager := note.NewManager(".")

	categoryPath, err := selectCategory(noteManager)
	if err != nil {
		return fmt.Errorf("failed to select category: %w", err)
	}

	title, err := promptForTitle()
	if err != nil {
		return fmt.Errorf("failed to get note title: %w", err)
	}

	notePath, err := noteManager.CreateNote(categoryPath, title)
	if err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Printf("Created note: %s\n", notePath)
	return nil
}

func selectCategory(noteManager *note.Manager) (string, error) {
	categories, err := noteManager.GetCategories()
	if err != nil {
		return "", err
	}

	var currentPath string

	for {
		options := make([]string, 0, len(categories)+1)
		options = append(options, "Create new category")
		options = append(options, categories...)

		if currentPath != "" {
			options = append(options, "Use current category")
		}

		prompt := promptui.Select{
			Label: "Select category",
			Items: options,
		}

		_, result, err := prompt.Run()
		if err != nil {
			return "", err
		}

		if result == "Create new category" {
			newCategory, err := promptForNewCategory()
			if err != nil {
				return "", err
			}

			categoryPath := filepath.Join(currentPath, newCategory)
			if err := noteManager.CreateCategory(categoryPath); err != nil {
				return "", fmt.Errorf("failed to create category: %w", err)
			}

			return categoryPath, nil
		}

		if result == "Use current category" {
			return currentPath, nil
		}

		newPath := filepath.Join(currentPath, result)

		subcategories, err := noteManager.GetSubcategories(newPath)
		if err != nil {
			return "", err
		}

		if len(subcategories) == 0 {
			return newPath, nil
		}

		currentPath = newPath
		categories = subcategories
	}
}

func promptForNewCategory() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter new category name",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("category name cannot be empty")
			}
			if strings.ContainsAny(input, "/\\:*?\"<>|") {
				return fmt.Errorf("category name contains invalid characters")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func promptForTitle() (string, error) {
	prompt := promptui.Prompt{
		Label: "Enter note title",
		Validate: func(input string) error {
			if strings.TrimSpace(input) == "" {
				return fmt.Errorf("note title cannot be empty")
			}
			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}
