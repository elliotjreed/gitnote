package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupGitRepo(t *testing.T) string {
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

func TestIsGitRepo(t *testing.T) {
	tempDir := setupGitRepo(t)
	manager := NewManager(tempDir)
	
	if !manager.IsGitRepo() {
		t.Error("Expected directory to be recognized as git repo")
	}
	
	nonGitDir := t.TempDir()
	manager = NewManager(nonGitDir)
	
	if manager.IsGitRepo() {
		t.Error("Expected directory to not be recognized as git repo")
	}
}

func TestGetStatus(t *testing.T) {
	tempDir := setupGitRepo(t)
	manager := NewManager(tempDir)
	
	status, err := manager.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	
	if len(status) != 0 {
		t.Errorf("Expected empty status for clean repo, got %d items", len(status))
	}
	
	testFile := filepath.Join(tempDir, "test.md")
	os.WriteFile(testFile, []byte("test content"), 0644)
	
	status, err = manager.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed after adding file: %v", err)
	}
	
	if len(status) != 1 {
		t.Fatalf("Expected 1 item in status, got %d", len(status))
	}
	
	if !strings.Contains(status[0], "test.md") {
		t.Errorf("Expected status to contain test.md, got %s", status[0])
	}
}

func TestAddFiles(t *testing.T) {
	tempDir := setupGitRepo(t)
	manager := NewManager(tempDir)
	
	testFile := filepath.Join(tempDir, "test.md")
	os.WriteFile(testFile, []byte("test content"), 0644)
	
	if err := manager.AddFiles([]string{"test.md"}); err != nil {
		t.Fatalf("AddFiles failed: %v", err)
	}
	
	status, err := manager.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	
	if len(status) != 1 {
		t.Fatalf("Expected 1 item in status after add, got %d", len(status))
	}
	
	if !strings.HasPrefix(status[0], "A ") {
		t.Errorf("Expected file to be staged (A), got status: %s", status[0])
	}
}

func TestCommit(t *testing.T) {
	tempDir := setupGitRepo(t)
	manager := NewManager(tempDir)
	
	testFile := filepath.Join(tempDir, "test.md")
	os.WriteFile(testFile, []byte("test content"), 0644)
	
	if err := manager.AddFiles([]string{"test.md"}); err != nil {
		t.Fatalf("AddFiles failed: %v", err)
	}
	
	if err := manager.Commit("Test commit"); err != nil {
		t.Fatalf("Commit failed: %v", err)
	}
	
	status, err := manager.GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}
	
	if len(status) != 0 {
		t.Errorf("Expected clean status after commit, got %d items", len(status))
	}
}

func TestAddFilesEmpty(t *testing.T) {
	tempDir := setupGitRepo(t)
	manager := NewManager(tempDir)
	
	if err := manager.AddFiles([]string{}); err != nil {
		t.Errorf("AddFiles with empty slice should not fail: %v", err)
	}
}

func TestHasMergeConflicts(t *testing.T) {
	tempDir := setupGitRepo(t)
	manager := NewManager(tempDir)
	
	hasConflicts, err := manager.HasMergeConflicts()
	if err != nil {
		t.Fatalf("HasMergeConflicts failed: %v", err)
	}
	
	if hasConflicts {
		t.Error("Expected no merge conflicts in clean repo")
	}
}