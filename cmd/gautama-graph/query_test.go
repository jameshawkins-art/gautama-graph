package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandleQueryCommand_MissingQuestion(t *testing.T) {
	err := handleQueryCommand([]string{})
	if err == nil {
		t.Fatal("expected error for missing question argument, got nil")
	}
	if !strings.Contains(err.Error(), "question") {
		t.Errorf("expected error mentioning question, got: %v", err)
	}
}

func TestHandlePathCommand_MissingEndpoints(t *testing.T) {
	// Zero args
	err := handlePathCommand([]string{})
	if err == nil {
		t.Fatal("expected error for empty args in path command, got nil")
	}

	// Only 1 arg
	err = handlePathCommand([]string{"nodeA"})
	if err == nil {
		t.Fatal("expected error for single node in path command, got nil")
	}
	if !strings.Contains(err.Error(), "both start and end nodes") {
		t.Errorf("expected error mentioning both nodes, got: %v", err)
	}
}

func TestHandleExplainCommand_MissingConcept(t *testing.T) {
	err := handleExplainCommand([]string{})
	if err == nil {
		t.Fatal("expected error for missing concept in explain command, got nil")
	}
	if !strings.Contains(err.Error(), "concept") {
		t.Errorf("expected error mentioning concept, got: %v", err)
	}
}

func TestHandleQueryCommand_MissingGraphGuidance(t *testing.T) {
	tempDir := t.TempDir()
	err := handleQueryCommand([]string{"--workspace=" + tempDir, "how does this work?"})
	if err == nil {
		t.Fatal("expected error when graphify-out/graph.json is missing, got nil")
	}
	if !strings.Contains(err.Error(), "graphify-out/graph.json") {
		t.Errorf("expected error mentioning graphify-out/graph.json, got: %v", err)
	}
}

func TestHandlePathCommand_MissingGraphGuidance(t *testing.T) {
	tempDir := t.TempDir()
	err := handlePathCommand([]string{"--workspace=" + tempDir, "nodeA", "nodeB"})
	if err == nil {
		t.Fatal("expected error when graphify-out/graph.json is missing, got nil")
	}
	if !strings.Contains(err.Error(), "graphify-out/graph.json") {
		t.Errorf("expected error mentioning graphify-out/graph.json, got: %v", err)
	}
}

func TestHandleExplainCommand_MissingGraphGuidance(t *testing.T) {
	tempDir := t.TempDir()
	err := handleExplainCommand([]string{"--workspace=" + tempDir, "ConceptX"})
	if err == nil {
		t.Fatal("expected error when graphify-out/graph.json is missing, got nil")
	}
	if !strings.Contains(err.Error(), "graphify-out/graph.json") {
		t.Errorf("expected error mentioning graphify-out/graph.json, got: %v", err)
	}
}

func TestGautamaGraphCLI_Query_Subcommand_MissingGraph(t *testing.T) {
	tempDir := t.TempDir()
	cmd := exec.Command("go", "run", ".", "query", "--workspace="+tempDir, "How does it work?")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected command to exit with non-zero code when graph.json is missing, got success. Output: %s", string(out))
	}
	outStr := string(out)
	if !strings.Contains(outStr, "gautama-graph") && !strings.Contains(outStr, "make graphify-update") {
		t.Errorf("expected output to contain guidance to run gautama-graph or make graphify-update, got: %s", outStr)
	}
}

func TestGautamaGraphCLI_Path_Subcommand_MissingEndpoints(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "path")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected command to fail when missing endpoints, got success. Output: %s", string(out))
	}
	if !strings.Contains(string(out), "start and end nodes") {
		t.Errorf("expected error mentioning start and end nodes, got: %s", string(out))
	}
}

func TestGautamaGraphCLI_Explain_Subcommand_MissingConcept(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "explain")
	out, err := cmd.CombinedOutput()
	if err == nil {
		t.Fatalf("expected command to fail when missing concept, got success. Output: %s", string(out))
	}
	if !strings.Contains(string(out), "concept") {
		t.Errorf("expected error mentioning concept, got: %s", string(out))
	}
}

func TestGautamaGraphCLI_Query_NominalWithRealGraph(t *testing.T) {
	// Create temporary workspace with graphify-out/graph.json
	tempDir := t.TempDir()
	graphDir := filepath.Join(tempDir, "graphify-out")
	if err := os.MkdirAll(graphDir, 0755); err != nil {
		t.Fatalf("failed to create graphify-out dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(graphDir, "graph.json"), []byte(`{"nodes":[], "edges":[]}`), 0644); err != nil {
		t.Fatalf("failed to write dummy graph.json: %v", err)
	}

	cmd := exec.Command("go", "run", ".", "query", "--workspace="+tempDir, "Test Question")
	out, _ := cmd.CombinedOutput()
	// Even if dummy graph has no nodes or binary runs, ensure command executed without panic
	outStr := string(out)
	if strings.Contains(outStr, "panic:") {
		t.Fatalf("unexpected panic during query CLI execution: %s", outStr)
	}
}

func TestStringSliceFlag(t *testing.T) {
	var f stringSliceFlag
	if err := f.Set("tag1"); err != nil {
		t.Fatalf("unexpected error setting flag: %v", err)
	}
	if err := f.Set("tag2"); err != nil {
		t.Fatalf("unexpected error setting flag: %v", err)
	}
	if f.String() != "tag1,tag2" {
		t.Errorf("expected 'tag1,tag2', got %q", f.String())
	}
}

func TestHandleQueryCommand_InProcess_WithRealGraph(t *testing.T) {
	tempDir := t.TempDir()
	graphDir := filepath.Join(tempDir, "graphify-out")
	if err := os.MkdirAll(graphDir, 0755); err != nil {
		t.Fatalf("failed to create graphify-out dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(graphDir, "graph.json"), []byte(`{"nodes":[], "edges":[]}`), 0644); err != nil {
		t.Fatalf("failed to write dummy graph.json: %v", err)
	}

	// In-process query
	_ = handleQueryCommand([]string{"--workspace=" + tempDir, "--dfs", "--budget=1000", "--context=ast", "Test Question"})

	// In-process path
	_ = handlePathCommand([]string{"--workspace=" + tempDir, "nodeA", "nodeB"})

	// In-process explain
	_ = handleExplainCommand([]string{"--workspace=" + tempDir, "ConceptX"})
}

