package index

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	generator := NewGenerator("")
	if generator.workingDir != "." {
		t.Errorf("Expected working dir to be '.', got %s", generator.workingDir)
	}
	
	generator = NewGenerator("/test")
	if generator.workingDir != "/test" {
		t.Errorf("Expected working dir to be '/test', got %s", generator.workingDir)
	}
}

func TestGenerateReadme(t *testing.T) {
	tempDir := t.TempDir()
	generator := NewGenerator(tempDir)
	
	os.Mkdir(filepath.Join(tempDir, "work"), 0755)
	os.Mkdir(filepath.Join(tempDir, "work", "meetings"), 0755)
	
	note1Path := filepath.Join(tempDir, "2025-01-01 root note.md")
	note2Path := filepath.Join(tempDir, "work", "2025-01-02 work note.md")
	note3Path := filepath.Join(tempDir, "work", "meetings", "2025-01-03 standup.md")
	
	os.WriteFile(note1Path, []byte("# root note"), 0644)
	os.WriteFile(note2Path, []byte("# work note"), 0644)
	os.WriteFile(note3Path, []byte("# standup"), 0644)
	
	if err := generator.GenerateReadme(); err != nil {
		t.Fatalf("GenerateReadme failed: %v", err)
	}
	
	readmePath := filepath.Join(tempDir, "readme.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read generated readme: %v", err)
	}
	
	readmeContent := string(content)
	
	if !strings.Contains(readmeContent, "# Notes Index") {
		t.Error("Expected readme to contain main heading")
	}
	
	if !strings.Contains(readmeContent, "[root note](/2025-01-01 root note.md)") {
		t.Error("Expected readme to contain root note link")
	}
	
	if !strings.Contains(readmeContent, "## work") {
		t.Error("Expected readme to contain work category heading")
	}
	
	if !strings.Contains(readmeContent, "### meetings") {
		t.Error("Expected readme to contain meetings subcategory heading")
	}
	
	if !strings.Contains(readmeContent, "[work note](/work/2025-01-02 work note.md)") {
		t.Error("Expected readme to contain work note link")
	}
	
	if !strings.Contains(readmeContent, "[standup](/work/meetings/2025-01-03 standup.md)") {
		t.Error("Expected readme to contain standup note link")
	}
}

func testIsReadmeUpToDate(t *testing.T) {
	tempDir := t.TempDir()
	generator := NewGenerator(tempDir)
	
	upToDate, err := generator.IsReadmeUpToDate()
	if err != nil {
		t.Fatalf("IsReadmeUpToDate failed: %v", err)
	}
	
	if upToDate {
		t.Error("Expected readme to not be up to date when it doesn't exist")
	}
	
	note1Path := filepath.Join(tempDir, "2025-01-01 test note.md")
	os.WriteFile(note1Path, []byte("# test note"), 0644)
	
	if err := generator.GenerateReadme(); err != nil {
		t.Fatalf("GenerateReadme failed: %v", err)
	}
	
	upToDate, err = generator.IsReadmeUpToDate()
	if err != nil {
		t.Fatalf("IsReadmeUpToDate failed: %v", err)
	}
	
	if !upToDate {
		t.Error("Expected readme to be up to date when content matches")
	}
	
	note2Path := filepath.Join(tempDir, "2025-01-02 another note.md")
	os.WriteFile(note2Path, []byte("# another note"), 0644)
	
	upToDate, err = generator.IsReadmeUpToDate()
	if err != nil {
		t.Fatalf("IsReadmeUpToDate failed: %v", err)
	}
	
	if upToDate {
		t.Error("Expected readme to not be up to date when notes were added")
	}
}

func TestBuildTableOfContents(t *testing.T) {
	content := "# Notes Index\n\n[root note](/2025-01-01 root note.md)\n\n## work\n\n[work note](/work/2025-01-02 work note.md)\n\n### meetings\n\n[meeting note](/work/meetings/2025-01-03 meeting.md)\n\n"
	
	if !strings.Contains(content, "# Notes Index") {
		t.Error("Expected content to contain main heading")
	}
	
	if !strings.Contains(content, "[root note]") {
		t.Error("Expected content to contain root note")
	}
	
	if !strings.Contains(content, "## work") {
		t.Error("Expected content to contain work heading")
	}
}