package note

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	manager := NewManager("")
	if manager.workingDir != "." {
		t.Errorf("Expected working dir to be '.', got %s", manager.workingDir)
	}
	
	manager = NewManager("/test")
	if manager.workingDir != "/test" {
		t.Errorf("Expected working dir to be '/test', got %s", manager.workingDir)
	}
}

func TestGetCategories(t *testing.T) {
	tempDir := t.TempDir()
	
	os.Mkdir(filepath.Join(tempDir, "work"), 0755)
	os.Mkdir(filepath.Join(tempDir, "personal"), 0755)
	os.Mkdir(filepath.Join(tempDir, ".hidden"), 0755)
	
	os.WriteFile(filepath.Join(tempDir, "file.md"), []byte("test"), 0644)
	
	manager := NewManager(tempDir)
	categories, err := manager.GetCategories()
	if err != nil {
		t.Fatalf("GetCategories failed: %v", err)
	}
	
	expected := []string{"personal", "work"}
	if len(categories) != len(expected) {
		t.Fatalf("Expected %d categories, got %d", len(expected), len(categories))
	}
	
	for i, cat := range expected {
		if categories[i] != cat {
			t.Errorf("Expected category %s, got %s", cat, categories[i])
		}
	}
}

func TestGetSubcategories(t *testing.T) {
	tempDir := t.TempDir()
	
	workDir := filepath.Join(tempDir, "work")
	os.Mkdir(workDir, 0755)
	os.Mkdir(filepath.Join(workDir, "meetings"), 0755)
	os.Mkdir(filepath.Join(workDir, "projects"), 0755)
	
	manager := NewManager(tempDir)
	subcategories, err := manager.GetSubcategories("work")
	if err != nil {
		t.Fatalf("GetSubcategories failed: %v", err)
	}
	
	expected := []string{"meetings", "projects"}
	if len(subcategories) != len(expected) {
		t.Fatalf("Expected %d subcategories, got %d", len(expected), len(subcategories))
	}
	
	for i, subcat := range expected {
		if subcategories[i] != subcat {
			t.Errorf("Expected subcategory %s, got %s", subcat, subcategories[i])
		}
	}
}

func TestCreateCategory(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)
	
	err := manager.CreateCategory("work/meetings")
	if err != nil {
		t.Fatalf("CreateCategory failed: %v", err)
	}
	
	categoryPath := filepath.Join(tempDir, "work", "meetings")
	if _, err := os.Stat(categoryPath); os.IsNotExist(err) {
		t.Errorf("Category directory was not created: %s", categoryPath)
	}
}

func TestCreateNote(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)
	
	notePath, err := manager.CreateNote("work", "test note")
	if err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}
	
	expectedPattern := "work" + string(filepath.Separator) + time.Now().Format("2006-01-02") + " test note.md"
	if notePath != expectedPattern {
		t.Errorf("Expected note path to match pattern %s, got %s", expectedPattern, notePath)
	}
	
	content, err := os.ReadFile(filepath.Join(tempDir, notePath))
	if err != nil {
		t.Fatalf("Failed to read created note: %v", err)
	}
	
	expectedContent := "# test note\n"
	if string(content) != expectedContent {
		t.Errorf("Expected content %q, got %q", expectedContent, string(content))
	}
}

func TestFindNotes(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)
	
	os.Mkdir(filepath.Join(tempDir, "work"), 0755)
	
	note1Path := filepath.Join(tempDir, "2025-01-01 root note.md")
	note2Path := filepath.Join(tempDir, "work", "2025-01-02 work note.md")
	
	os.WriteFile(note1Path, []byte("# root note"), 0644)
	os.WriteFile(note2Path, []byte("# work note"), 0644)
	
	notes, err := manager.FindNotes()
	if err != nil {
		t.Fatalf("FindNotes failed: %v", err)
	}
	
	if len(notes) != 2 {
		t.Fatalf("Expected 2 notes, got %d", len(notes))
	}
	
	note1 := notes[0]
	if note1.Title != "root note" {
		t.Errorf("Expected title 'root note', got %s", note1.Title)
	}
	if note1.Category != "" {
		t.Errorf("Expected empty category, got %s", note1.Category)
	}
	
	note2 := notes[1]
	if note2.Title != "work note" {
		t.Errorf("Expected title 'work note', got %s", note2.Title)
	}
	if note2.Category != "work" {
		t.Errorf("Expected category 'work', got %s", note2.Category)
	}
}

func TestSearchNotes(t *testing.T) {
	tempDir := t.TempDir()
	manager := NewManager(tempDir)
	
	os.Mkdir(filepath.Join(tempDir, "work"), 0755)
	
	note1Path := filepath.Join(tempDir, "2025-01-01 meeting notes.md")
	note2Path := filepath.Join(tempDir, "work", "2025-01-02 project update.md")
	note3Path := filepath.Join(tempDir, "2025-01-03 random thoughts.md")
	
	os.WriteFile(note1Path, []byte("# meeting notes\nDiscussed project timeline"), 0644)
	os.WriteFile(note2Path, []byte("# project update\nMeeting went well"), 0644)
	os.WriteFile(note3Path, []byte("# random thoughts\nJust some ideas"), 0644)
	
	results, err := manager.SearchNotes("meeting", false)
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}
	
	if len(results) != 1 {
		t.Fatalf("Expected 1 result for title search, got %d", len(results))
	}
	if results[0].Title != "meeting notes" {
		t.Errorf("Expected 'meeting notes', got %s", results[0].Title)
	}
	
	results, err = manager.SearchNotes("meeting", true)
	if err != nil {
		t.Fatalf("SearchNotes with content failed: %v", err)
	}
	
	if len(results) != 2 {
		t.Fatalf("Expected 2 results for content search, got %d", len(results))
	}
	
	results, err = manager.SearchNotes("nonexistent", false)
	if err != nil {
		t.Fatalf("SearchNotes failed: %v", err)
	}
	
	if len(results) != 0 {
		t.Fatalf("Expected 0 results for nonexistent search, got %d", len(results))
	}
}