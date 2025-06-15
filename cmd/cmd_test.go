package cmd

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

func setupTestRepo(t *testing.T) string {
	tempDir := t.TempDir()
	
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set git email: %v", err)
	}
	
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to set git name: %v", err)
	}
	
	return tempDir
}

func TestIndexCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	os.Mkdir("work", 0755)
	
	note1Path := "2025-01-01 root note.md"
	note2Path := "work/2025-01-02 work note.md"
	
	os.WriteFile(note1Path, []byte("# root note"), 0644)
	os.WriteFile(note2Path, []byte("# work note"), 0644)
	
	if err := runIndex(nil, []string{}); err != nil {
		t.Fatalf("runIndex failed: %v", err)
	}
	
	readmeContent, err := os.ReadFile("readme.md")
	if err != nil {
		t.Fatalf("Failed to read generated readme: %v", err)
	}
	
	content := string(readmeContent)
	if !strings.Contains(content, "# Notes Index") {
		t.Error("Expected readme to contain main heading")
	}
	
	if !strings.Contains(content, "[root note](/2025-01-01 root note.md)") {
		t.Error("Expected readme to contain root note link")
	}
	
	if !strings.Contains(content, "## work") {
		t.Error("Expected readme to contain work category heading")
	}
}

func TestSearchCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	os.Mkdir("work", 0755)
	
	note1Path := "2025-01-01 meeting notes.md"
	note2Path := "work/2025-01-02 project update.md"
	note3Path := "2025-01-03 random thoughts.md"
	
	os.WriteFile(note1Path, []byte("# meeting notes\nDiscussed project timeline"), 0644)
	os.WriteFile(note2Path, []byte("# project update\nMeeting went well"), 0644)
	os.WriteFile(note3Path, []byte("# random thoughts\nJust some ideas"), 0644)
	
	if err := runSearch(nil, []string{"meeting"}); err != nil {
		t.Fatalf("runSearch failed: %v", err)
	}
	
	searchFull = true
	if err := runSearch(nil, []string{"project"}); err != nil {
		t.Fatalf("runSearch with full flag failed: %v", err)
	}
	searchFull = false
}

func TestCommitCommand(t *testing.T) {
	tempDir := setupTestRepo(t)
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	
	testFile := "test.md"
	os.WriteFile(testFile, []byte("# test note"), 0644)
	
	if err := runCommit(nil, []string{}); err != nil {
		t.Fatalf("runCommit failed: %v", err)
	}
	
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Failed to get git status: %v", err)
	}
	
	if len(strings.TrimSpace(string(output))) != 0 {
		t.Error("Expected clean git status after commit")
	}
}

func TestBuildCommitMessage(t *testing.T) {
	tests := []struct {
		name         string
		newFiles     []string
		modifiedFiles []string
		expected     string
	}{
		{
			name:     "single new file",
			newFiles: []string{"test.md"},
			expected: "Add test.md",
		},
		{
			name:     "multiple new files",
			newFiles: []string{"test1.md", "test2.md"},
			expected: "Add 2 new files",
		},
		{
			name:          "single modified file",
			modifiedFiles: []string{"existing.md"},
			expected:      "Update existing.md",
		},
		{
			name:          "multiple modified files",
			modifiedFiles: []string{"file1.md", "file2.md"},
			expected:      "Update 2 files",
		},
		{
			name:          "new and modified files",
			newFiles:      []string{"new.md"},
			modifiedFiles: []string{"existing.md"},
			expected:      "Add new.md and Update existing.md",
		},
		{
			name:     "no files",
			expected: "Update files",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildCommitMessage(tt.newFiles, tt.modifiedFiles)
			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}