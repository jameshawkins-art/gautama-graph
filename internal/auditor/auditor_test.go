package auditor_test

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"gautama-graph/internal/auditor"
)

func TestASTInferredRelationshipAuditor(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Create a sample Go test source file with explicit selector call (t.Run) and helper call (at)
	sampleGoSource := `package sample_test

import "testing"

func at() {}

func TestSample(t *testing.T) {
	at()
	t.Run("subtest", func(t *testing.T) {
		// nested
	})
}
`
	sampleFilePath := filepath.Join(tempDir, "sample_test.go")
	if err := os.WriteFile(sampleFilePath, []byte(sampleGoSource), 0644); err != nil {
		t.Fatalf("failed to write sample Go test file: %v", err)
	}

	// 2. Instantiate Auditor components
	parser := auditor.NewDefaultASTParser(tempDir)
	evaluator := auditor.NewDefaultSelectorEvaluator(50)
	store := auditor.NewJSONGraphStore()

	cfg := auditor.Config{
		MaxASTDepth:       50,
		MinConfidence:     0.7,
		AuditorTimeout:    5 * time.Second,
		WorkspaceRootPath: tempDir,
	}

	engine := auditor.NewEngine(parser, evaluator, store, cfg)

	// 3. Define candidate inferred edges
	candidates := []auditor.CandidateEdge{
		{
			ID:              "sample_test.go:TestSample->testing.T:Run",
			SourceFile:      sampleFilePath,
			SourceSymbol:    "t",
			TargetSymbol:    "Run",
			InitialRelation: "CALLS",
		},
		{
			ID:              "sample_test.go:TestSample->sample_test:at",
			SourceFile:      sampleFilePath,
			SourceSymbol:    "",
			TargetSymbol:    "at",
			InitialRelation: "CALLS",
		},
		{
			ID:              "sample_test.go:TestSample->phantom:d",
			SourceFile:      sampleFilePath,
			SourceSymbol:    "phantom",
			TargetSymbol:    "d",
			InitialRelation: "INFERRED",
		},
	}

	// 4. Audit Candidates
	ctx := context.Background()
	audited, err := engine.AuditCandidates(ctx, candidates)
	if err != nil {
		t.Fatalf("AuditCandidates failed: %v", err)
	}

	if len(audited) != 3 {
		t.Fatalf("expected 3 audited edges, got %d", len(audited))
	}

	// Assert edge 0: t.Run -> VERIFIED / EXTRACTED_AST (confidence 1.0)
	if audited[0].ProvenanceStatus != "EXTRACTED_AST" || audited[0].Confidence != 1.0 {
		t.Errorf("edge 0 expected EXTRACTED_AST with confidence 1.0, got %s (%.1f)", audited[0].ProvenanceStatus, audited[0].Confidence)
	}

	// Assert edge 1: at() helper -> VERIFIED / EXTRACTED_AST (confidence 1.0)
	if audited[1].ProvenanceStatus != "EXTRACTED_AST" || audited[1].Confidence != 1.0 {
		t.Errorf("edge 1 expected EXTRACTED_AST with confidence 1.0, got %s (%.1f)", audited[1].ProvenanceStatus, audited[1].Confidence)
	}

	// Assert edge 2: phantom.d -> PRUNED_PHANTOM (confidence 0.0)
	if audited[2].ProvenanceStatus != "PRUNED_PHANTOM" || audited[2].Confidence != 0.0 {
		t.Errorf("edge 2 expected PRUNED_PHANTOM with confidence 0.0, got %s (%.1f)", audited[2].ProvenanceStatus, audited[2].Confidence)
	}

	// 5. Test JSON GraphStore update
	targetGraphJSON := filepath.Join(tempDir, "graph.json")
	if err := store.SaveAuditedEdges(ctx, targetGraphJSON, audited); err != nil {
		t.Fatalf("SaveAuditedEdges failed: %v", err)
	}

	if _, err := os.Stat(targetGraphJSON); os.IsNotExist(err) {
		t.Fatalf("expected graph.json artifact to exist at %s", targetGraphJSON)
	}
}

func TestASTAuditorPathBoundarySafety(t *testing.T) {
	tempDir := t.TempDir()
	parser := auditor.NewDefaultASTParser(tempDir)
	ctx := context.Background()

	// Attempt to parse path outside workspace root
	_, _, err := parser.ParseFile(ctx, "/etc/passwd")
	if err == nil {
		t.Errorf("expected security error when parsing path outside workspace root, got nil")
	}
}

func TestAuditGraphFile_FullPipeline(t *testing.T) {
	tempDir := t.TempDir()

	// 1. Create a sample Go source file
	sampleGoSource := `package mathutil

func HelperAdd(a, b int) int {
	return a + b
}

func CalculateTotal(items []int) int {
	sum := 0
	for _, it := range items {
		sum = HelperAdd(sum, it)
	}
	return sum
}
`
	sampleGoPath := filepath.Join(tempDir, "mathutil.go")
	if err := os.WriteFile(sampleGoPath, []byte(sampleGoSource), 0644); err != nil {
		t.Fatalf("failed to write mathutil.go: %v", err)
	}

	// 2. Create mock graph.json
	graph := auditor.GraphData{
		Nodes: []map[string]interface{}{
			{
				"id":          "node_calculate_total",
				"label":       "CalculateTotal",
				"source_file": "mathutil.go",
				"file_type":   "code",
			},
			{
				"id":          "node_helper_add",
				"label":       "HelperAdd",
				"source_file": "mathutil.go",
				"file_type":   "code",
			},
			{
				"id":          "node_phantom_target",
				"label":       "NonExistentFunc",
				"source_file": "mathutil.go",
				"file_type":   "code",
			},
		},
		Links: []map[string]interface{}{
			{
				"source":           "node_calculate_total",
				"target":           "node_helper_add",
				"relation":         "calls",
				"confidence":       "INFERRED",
				"confidence_score": 0.5,
				"source_file":      "mathutil.go",
			},
			{
				"source":           "node_calculate_total",
				"target":           "node_phantom_target",
				"relation":         "inferred",
				"confidence":       "INFERRED",
				"confidence_score": 0.5,
				"source_file":      "mathutil.go",
			},
		},
	}

	graphBytes, err := json.MarshalIndent(graph, "", "  ")
	if err != nil {
		t.Fatalf("failed to marshal mock graph: %v", err)
	}

	graphPath := filepath.Join(tempDir, "graph.json")
	if err := os.WriteFile(graphPath, graphBytes, 0644); err != nil {
		t.Fatalf("failed to write graph.json: %v", err)
	}

	// 3. Instantiate DefaultEngine
	cfg := auditor.Config{
		WorkspaceRootPath: tempDir,
		AuditorTimeout:    10 * time.Second,
		MinConfidence:     0.8,
	}

	engine := auditor.NewDefaultEngine(cfg)
	report, err := engine.AuditGraphFile(context.Background(), graphPath, true)
	if err != nil {
		t.Fatalf("AuditGraphFile failed: %v", err)
	}

	if report.TotalEdges != 2 {
		t.Errorf("expected total edges 2, got %d", report.TotalEdges)
	}
	if report.VerifiedASTCount != 1 {
		t.Errorf("expected verified AST count 1, got %d", report.VerifiedASTCount)
	}
	if report.PrunedPhantomCount != 1 {
		t.Errorf("expected pruned phantom count 1, got %d", report.PrunedPhantomCount)
	}

	// 4. Verify persisted graph.json only contains verified link (phantom pruned)
	updatedData, err := os.ReadFile(graphPath)
	if err != nil {
		t.Fatalf("failed to read updated graph.json: %v", err)
	}

	var updatedGraph auditor.GraphData
	if err := json.Unmarshal(updatedData, &updatedGraph); err != nil {
		t.Fatalf("failed to unmarshal updated graph.json: %v", err)
	}

	if len(updatedGraph.Links) != 1 {
		t.Errorf("expected 1 remaining link after phantom pruning, got %d", len(updatedGraph.Links))
	}
	if updatedGraph.Links[0]["target"] != "node_helper_add" {
		t.Errorf("expected remaining link to target node_helper_add, got %v", updatedGraph.Links[0]["target"])
	}
}

func TestGraphifyASTAuditCLI(t *testing.T) {
	tempDir := t.TempDir()
	graphDir := filepath.Join(tempDir, "graphify-out")
	if err := os.MkdirAll(graphDir, 0755); err != nil {
		t.Fatalf("failed to create graphify-out: %v", err)
	}
	mockGraph := auditor.GraphData{}
	b, err := json.Marshal(mockGraph)
	if err != nil {
		t.Fatalf("failed to marshal mock graph: %v", err)
	}
	if err := os.WriteFile(filepath.Join(graphDir, "graph.json"), b, 0644); err != nil {
		t.Fatalf("failed to write graph.json: %v", err)
	}

	rootDir, err := filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatalf("failed to locate workspace root directory: %v", err)
	}

	cmd := exec.Command("go", "run", "cmd/graphify-ast-audit/main.go", "--workspace", tempDir)
	cmd.Dir = rootDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected cmd/graphify-ast-audit to succeed, got error: %v, output: %s", err, string(output))
	}
}
