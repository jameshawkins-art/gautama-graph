package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGautamaGraphCLI_Help(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "-h")
	out, _ := cmd.CombinedOutput()

	// Exit code from -h flag in Go `flag` package is 2
	if len(out) == 0 {
		t.Errorf("expected help output from CLI, got empty")
	}
}

func TestGautamaGraphCLI_AntigravityInstall_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	cmd := exec.Command("go", "run", ".", "antigravity", "install", "--workspace="+tempDir, "--dry-run")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v, output: %s", err, string(out))
	}

	if !strings.Contains(string(out), "Gautama Graph Antigravity Environment Setup") {
		t.Errorf("expected banner in output, got: %s", string(out))
	}
	if !strings.Contains(string(out), "Dry-run completed") {
		t.Errorf("expected dry-run message, got: %s", string(out))
	}

	// Verify no files on disk
	rulesPath := filepath.Join(tempDir, ".agents", "rules", "graphify.md")
	if _, err := os.Stat(rulesPath); err == nil {
		t.Errorf("expected %s NOT to exist during dry-run", rulesPath)
	}
}

func TestGautamaGraphCLI_AntigravityInstall_Execution(t *testing.T) {
	tempDir := t.TempDir()
	cmd := exec.Command("go", "run", ".", "antigravity", "install", "--workspace="+tempDir, "--project")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("command failed: %v, output: %s", err, string(out))
	}

	if !strings.Contains(string(out), "Files successfully written to workspace") {
		t.Errorf("expected success message, got: %s", string(out))
	}

	// Verify files on disk
	rulesPath := filepath.Join(tempDir, ".agents", "rules", "graphify.md")
	if _, err := os.Stat(rulesPath); err != nil {
		t.Errorf("expected %s to exist on disk", rulesPath)
	}

	scriptPath := filepath.Join(tempDir, "scripts", "graphify_sync.sh")
	if _, err := os.Stat(scriptPath); err != nil {
		t.Errorf("expected %s to exist on disk", scriptPath)
	}
}
