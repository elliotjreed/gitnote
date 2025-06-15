package index

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gitnote/internal/note"
)

type Generator struct {
	workingDir string
	noteManager *note.Manager
}

func NewGenerator(workingDir string) *Generator {
	if workingDir == "" {
		workingDir = "."
	}
	return &Generator{
		workingDir: workingDir,
		noteManager: note.NewManager(workingDir),
	}
}

func (g *Generator) GenerateReadme() error {
	notes, err := g.noteManager.FindNotes()
	if err != nil {
		return fmt.Errorf("failed to find notes: %w", err)
	}
	
	content := g.buildTableOfContents(notes)
	
	readmePath := filepath.Join(g.workingDir, "readme.md")
	
	if err := os.WriteFile(readmePath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write readme.md: %w", err)
	}
	
	return nil
}

func (g *Generator) buildTableOfContents(notes []note.Note) string {
	var content strings.Builder
	content.WriteString("# Notes Index\n\n")
	
	notesByCategory := make(map[string][]note.Note)
	
	for _, note := range notes {
		category := note.Category
		if category == "" {
			category = "root"
		}
		notesByCategory[category] = append(notesByCategory[category], note)
	}
	
	var categories []string
	for category := range notesByCategory {
		categories = append(categories, category)
	}
	sort.Strings(categories)
	
	for _, category := range categories {
		categoryNotes := notesByCategory[category]
		
		if category == "root" {
			for _, note := range categoryNotes {
				link := fmt.Sprintf("[%s](/%s)", note.Title, note.Path)
				content.WriteString(fmt.Sprintf("%s\n", link))
			}
			if len(categoryNotes) > 0 {
				content.WriteString("\n")
			}
		} else {
			g.buildCategorySection(&content, category, categoryNotes)
		}
	}
	
	return content.String()
}

func (g *Generator) buildCategorySection(content *strings.Builder, category string, notes []note.Note) {
	parts := strings.Split(category, string(filepath.Separator))
	
	subcategoryNotes := make(map[string][]note.Note)
	var directNotes []note.Note
	
	for _, note := range notes {
		noteParts := strings.Split(note.Category, string(filepath.Separator))
		if len(noteParts) == len(parts) {
			directNotes = append(directNotes, note)
		} else if len(noteParts) > len(parts) {
			subcat := strings.Join(noteParts[:len(parts)+1], string(filepath.Separator))
			subcategoryNotes[subcat] = append(subcategoryNotes[subcat], note)
		}
	}
	
	for i, part := range parts {
		level := strings.Repeat("#", i+2)
		content.WriteString(fmt.Sprintf("%s %s\n\n", level, part))
	}
	
	for _, note := range directNotes {
		link := fmt.Sprintf("[%s](/%s)", note.Title, note.Path)
		content.WriteString(fmt.Sprintf("%s\n", link))
	}
	
	if len(directNotes) > 0 {
		content.WriteString("\n")
	}
	
	var subcats []string
	for subcat := range subcategoryNotes {
		subcats = append(subcats, subcat)
	}
	sort.Strings(subcats)
	
	for _, subcat := range subcats {
		g.buildCategorySection(content, subcat, subcategoryNotes[subcat])
	}
}

func (g *Generator) IsReadmeUpToDate() (bool, error) {
	readmePath := filepath.Join(g.workingDir, "readme.md")
	
	if _, err := os.Stat(readmePath); os.IsNotExist(err) {
		return false, nil
	}
	
	existingContent, err := os.ReadFile(readmePath)
	if err != nil {
		return false, fmt.Errorf("failed to read existing readme: %w", err)
	}
	
	notes, err := g.noteManager.FindNotes()
	if err != nil {
		return false, fmt.Errorf("failed to find notes: %w", err)
	}
	
	expectedContent := g.buildTableOfContents(notes)
	
	return string(existingContent) == expectedContent, nil
}