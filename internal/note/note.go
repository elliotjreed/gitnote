package note

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Note struct {
	Title    string
	Path     string
	Category string
	Date     time.Time
}

type Manager struct {
	workingDir string
}

func NewManager(workingDir string) *Manager {
	if workingDir == "" {
		workingDir = "."
	}
	return &Manager{workingDir: workingDir}
}

func (m *Manager) GetCategories() ([]string, error) {
	var categories []string
	
	entries, err := os.ReadDir(m.workingDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read working directory: %w", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			categories = append(categories, entry.Name())
		}
	}
	
	sort.Strings(categories)
	return categories, nil
}

func (m *Manager) GetSubcategories(category string) ([]string, error) {
	var subcategories []string
	categoryPath := filepath.Join(m.workingDir, category)
	
	entries, err := os.ReadDir(categoryPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read category directory: %w", err)
	}
	
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			subcategories = append(subcategories, entry.Name())
		}
	}
	
	sort.Strings(subcategories)
	return subcategories, nil
}

func (m *Manager) CreateCategory(path string) error {
	fullPath := filepath.Join(m.workingDir, path)
	return os.MkdirAll(fullPath, 0755)
}

func (m *Manager) CreateNote(categoryPath, title string) (string, error) {
	now := time.Now()
	filename := fmt.Sprintf("%s %s.md", now.Format("2006-01-02"), title)
	
	fullPath := filepath.Join(m.workingDir, categoryPath, filename)
	
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create directory: %w", err)
	}
	
	content := fmt.Sprintf("# %s\n", title)
	
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to create note file: %w", err)
	}
	
	relativePath, err := filepath.Rel(m.workingDir, fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %w", err)
	}
	
	return relativePath, nil
}

func (m *Manager) FindNotes() ([]Note, error) {
	var notes []Note
	
	err := filepath.Walk(m.workingDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".md") {
			relativePath, err := filepath.Rel(m.workingDir, path)
			if err != nil {
				return err
			}
			
			if strings.HasPrefix(relativePath, ".") {
				return nil
			}
			
			note := Note{
				Path: relativePath,
				Date: info.ModTime(),
			}
			
			dirPath := filepath.Dir(relativePath)
			if dirPath != "." {
				note.Category = dirPath
			}
			
			filename := filepath.Base(relativePath)
			filename = strings.TrimSuffix(filename, ".md")
			
			if len(filename) > 10 && filename[4] == '-' && filename[7] == '-' {
				note.Title = strings.TrimSpace(filename[11:])
			} else {
				note.Title = filename
			}
			
			notes = append(notes, note)
		}
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to walk directory: %w", err)
	}
	
	sort.Slice(notes, func(i, j int) bool {
		return notes[i].Path < notes[j].Path
	})
	
	return notes, nil
}

func (m *Manager) SearchNotes(query string, searchContent bool) ([]Note, error) {
	allNotes, err := m.FindNotes()
	if err != nil {
		return nil, err
	}
	
	var matchingNotes []Note
	queryLower := strings.ToLower(query)
	
	for _, note := range allNotes {
		titleMatch := strings.Contains(strings.ToLower(note.Title), queryLower)
		
		if titleMatch {
			matchingNotes = append(matchingNotes, note)
			continue
		}
		
		if searchContent {
			fullPath := filepath.Join(m.workingDir, note.Path)
			content, err := os.ReadFile(fullPath)
			if err != nil {
				continue
			}
			
			if strings.Contains(strings.ToLower(string(content)), queryLower) {
				matchingNotes = append(matchingNotes, note)
			}
		}
	}
	
	return matchingNotes, nil
}