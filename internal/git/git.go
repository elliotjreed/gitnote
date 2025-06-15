package git

import (
	"fmt"
	"os/exec"
	"strings"
)

type Manager struct {
	workingDir string
}

func NewManager(workingDir string) *Manager {
	if workingDir == "" {
		workingDir = "."
	}
	return &Manager{workingDir: workingDir}
}

func (g *Manager) IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Dir = g.workingDir
	return cmd.Run() == nil
}

func (g *Manager) GetStatus() ([]string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.workingDir
	
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get git status: %w", err)
	}
	
	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return []string{}, nil
	}
	
	return lines, nil
}

func (g *Manager) AddFiles(files []string) error {
	if len(files) == 0 {
		return nil
	}
	
	args := append([]string{"add"}, files...)
	cmd := exec.Command("git", args...)
	cmd.Dir = g.workingDir
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}
	
	return nil
}

func (g *Manager) Commit(message string) error {
	cmd := exec.Command("git", "commit", "-m", message)
	cmd.Dir = g.workingDir
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}
	
	return nil
}

func (g *Manager) Pull() (string, error) {
	cmd := exec.Command("git", "pull")
	cmd.Dir = g.workingDir
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("pull failed: %w", err)
	}
	
	return string(output), nil
}

func (g *Manager) HasMergeConflicts() (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = g.workingDir
	
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check merge conflicts: %w", err)
	}
	
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "UU ") || strings.HasPrefix(line, "AA ") || 
		   strings.HasPrefix(line, "DD ") || strings.HasPrefix(line, "AU ") ||
		   strings.HasPrefix(line, "UA ") || strings.HasPrefix(line, "DU ") ||
		   strings.HasPrefix(line, "UD ") {
			return true, nil
		}
	}
	
	return false, nil
}

func (g *Manager) Reset() error {
	cmd := exec.Command("git", "reset", "--hard", "HEAD")
	cmd.Dir = g.workingDir
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to reset: %w", err)
	}
	
	return nil
}