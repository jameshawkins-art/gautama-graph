package main

import (
	"os/exec"
	"testing"
)

func TestGautamaGraphCLI_Help(t *testing.T) {
	cmd := exec.Command("go", "run", "main.go", "-h")
	out, _ := cmd.CombinedOutput()

	// Exit code from -h flag in Go `flag` package is 2
	if len(out) == 0 {
		t.Errorf("expected help output from CLI, got empty")
	}
}
