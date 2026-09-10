package runner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type queryMockBinaryManager struct {
	binaryPath string
	version    string
	err        error
	calledWith RunnerConfig
}

func (m *queryMockBinaryManager) EnsureBinary(ctx context.Context, cfg RunnerConfig) (string, string, error) {
	m.calledWith = cfg
	if m.err != nil {
		return "", "", m.err
	}
	return m.binaryPath, m.version, nil
}

type queryMockSubprocessRunner struct {
	stdout     []byte
	stderr     []byte
	err        error
	calledBin  string
	calledRoot string
	calledArgs []string
}

func (m *queryMockSubprocessRunner) ExecuteCommand(ctx context.Context, binaryPath, workspaceRoot string, args ...string) ([]byte, []byte, error) {
	m.calledBin = binaryPath
	m.calledRoot = workspaceRoot
	m.calledArgs = args
	if m.err != nil {
		return m.stdout, m.stderr, m.err
	}
	return m.stdout, m.stderr, nil
}

func setupTestWorkspaceWithGraph(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	graphDir := filepath.Join(dir, "graphify-out")
	if err := os.MkdirAll(graphDir, 0755); err != nil {
		t.Fatalf("failed to create graphify-out dir: %v", err)
	}
	graphFile := filepath.Join(graphDir, "graph.json")
	if err := os.WriteFile(graphFile, []byte(`{"nodes":[], "edges":[]}`), 0644); err != nil {
		t.Fatalf("failed to write graph.json: %v", err)
	}
	return dir
}

func TestDefaultQueryService_Query_BFS_Nominal(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{stdout: []byte("Traversal: BFS depth=2 | Found 3 nodes")}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	res, err := svc.Query(context.Background(), "How does runner work?", QueryOptions{
		WorkspaceRoot: ws,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil {
		t.Fatal("expected non-nil QueryResult")
	}
	if !strings.Contains(res.Output, "Found 3 nodes") {
		t.Errorf("expected output to contain 'Found 3 nodes', got %q", res.Output)
	}
	if res.BinaryVersion != "v1.5.0" {
		t.Errorf("expected version v1.5.0, got %q", res.BinaryVersion)
	}

	// Verify argv
	expectedArgs := []string{"query", "How does runner work?"}
	if len(mockRunner.calledArgs) != len(expectedArgs) {
		t.Fatalf("expected %d args, got %d: %v", len(expectedArgs), len(mockRunner.calledArgs), mockRunner.calledArgs)
	}
	for i, arg := range expectedArgs {
		if mockRunner.calledArgs[i] != arg {
			t.Errorf("arg[%d]: expected %q, got %q", i, arg, mockRunner.calledArgs[i])
		}
	}
}

func TestDefaultQueryService_Query_DFS_Budget_And_Context(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{stdout: []byte("DFS traversal output")}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	res, err := svc.Query(context.Background(), "Trace call chain", QueryOptions{
		WorkspaceRoot: ws,
		DFS:           true,
		BudgetTokens:  4000,
		Context:       []string{"ast", "doc"},
		JSONOutput:    true,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil || res.Output != "DFS traversal output" {
		t.Errorf("unexpected result: %+v", res)
	}

	// Verify flags in args
	argsStr := strings.Join(mockRunner.calledArgs, " ")
	if !strings.Contains(argsStr, "--dfs") {
		t.Errorf("expected --dfs flag in args: %v", mockRunner.calledArgs)
	}
	if !strings.Contains(argsStr, "--budget 4000") {
		t.Errorf("expected --budget 4000 in args: %v", mockRunner.calledArgs)
	}
	if !strings.Contains(argsStr, "--context ast") || !strings.Contains(argsStr, "--context doc") {
		t.Errorf("expected context tags in args: %v", mockRunner.calledArgs)
	}
	if !strings.Contains(argsStr, "--json") {
		t.Errorf("expected --json flag in args: %v", mockRunner.calledArgs)
	}
}

func TestDefaultQueryService_Path_Nominal(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{stdout: []byte("Path: A -> B")}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	res, err := svc.Path(context.Background(), "nodeA", "nodeB", PathOptions{
		WorkspaceRoot: ws,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil || res.Output != "Path: A -> B" {
		t.Errorf("unexpected result: %+v", res)
	}

	expectedArgs := []string{"path", "nodeA", "nodeB"}
	if strings.Join(mockRunner.calledArgs, " ") != strings.Join(expectedArgs, " ") {
		t.Errorf("expected args %v, got %v", expectedArgs, mockRunner.calledArgs)
	}
}

func TestDefaultQueryService_Explain_Nominal(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{stdout: []byte("Explanation of concept")}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	res, err := svc.Explain(context.Background(), "ConceptX", ExplainOptions{
		WorkspaceRoot: ws,
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res == nil || res.Output != "Explanation of concept" {
		t.Errorf("unexpected result: %+v", res)
	}

	expectedArgs := []string{"explain", "ConceptX"}
	if strings.Join(mockRunner.calledArgs, " ") != strings.Join(expectedArgs, " ") {
		t.Errorf("expected args %v, got %v", expectedArgs, mockRunner.calledArgs)
	}
}

func TestDefaultQueryService_MissingGraphGuidance(t *testing.T) {
	emptyDir := t.TempDir() // no graphify-out/graph.json
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	_, err := svc.Query(context.Background(), "test question", QueryOptions{
		WorkspaceRoot: emptyDir,
	})
	if err == nil {
		t.Fatal("expected error when graph.json is missing, got nil")
	}

	if !errors.Is(err, ErrGraphNotFound) {
		t.Errorf("expected ErrGraphNotFound, got %v", err)
	}

	if !strings.Contains(err.Error(), "run 'gautama-graph' or 'make graphify-update' first") {
		t.Errorf("expected actionable guidance in error message, got %q", err.Error())
	}
}

func TestDefaultQueryService_PathBoundaryRejection(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	_, err := svc.Query(context.Background(), "test question", QueryOptions{
		WorkspaceRoot: ws,
		GraphPath:     filepath.Clean("/etc/passwd"),
	})
	if err == nil {
		t.Fatal("expected error for path traversal escape, got nil")
	}
	if !errors.Is(err, ErrPathOutOfBounds) {
		t.Errorf("expected ErrPathOutOfBounds, got %v", err)
	}
}

func TestDefaultQueryService_EmptyInputs(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	svc := NewDefaultQueryService(nil, nil)

	if _, err := svc.Query(context.Background(), "", QueryOptions{WorkspaceRoot: ws}); err == nil {
		t.Error("expected error for empty question")
	}
	if _, err := svc.Path(context.Background(), "", "B", PathOptions{WorkspaceRoot: ws}); err == nil {
		t.Error("expected error for empty nodeA")
	}
	if _, err := svc.Path(context.Background(), "A", "", PathOptions{WorkspaceRoot: ws}); err == nil {
		t.Error("expected error for empty nodeB")
	}
	if _, err := svc.Explain(context.Background(), "", ExplainOptions{WorkspaceRoot: ws}); err == nil {
		t.Error("expected error for empty concept")
	}
}

func TestDefaultQueryService_SubprocessError(t *testing.T) {
	ws := setupTestWorkspaceWithGraph(t)
	mockMgr := &queryMockBinaryManager{binaryPath: "/bin/graphify", version: "v1.5.0"}
	mockRunner := &queryMockSubprocessRunner{err: errors.New("command failed: exit status 1"), stderr: []byte("node not found")}

	svc := NewDefaultQueryService(mockMgr, mockRunner)
	_, err := svc.Query(context.Background(), "Where is Foo?", QueryOptions{
		WorkspaceRoot: ws,
	})
	if err == nil {
		t.Fatal("expected error on subprocess failure, got nil")
	}
	if !strings.Contains(err.Error(), "command failed") {
		t.Errorf("expected command failure error, got: %v", err)
	}
}
