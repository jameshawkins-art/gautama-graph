package auditor

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCalculateCanonicalRelPath(t *testing.T) {
	tests := []struct {
		name       string
		sourceFile string
		targetFile string
		expected   string
	}{
		{
			name:       "sibling in same directory",
			sourceFile: "docs/specs/001-reqs.md",
			targetFile: "docs/specs/001-arch.md",
			expected:   "./001-arch.md",
		},
		{
			name:       "child to sibling directory",
			sourceFile: "docs/specs/001-reqs.md",
			targetFile: "docs/roadmap/001-roadmap.md",
			expected:   "../roadmap/001-roadmap.md",
		},
		{
			name:       "nested child to root",
			sourceFile: "docs/specs/sub/test.md",
			targetFile: "README.md",
			expected:   "../../../README.md",
		},
		{
			name:       "root to nested child",
			sourceFile: "README.md",
			targetFile: "docs/specs/001.md",
			expected:   "./docs/specs/001.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := CalculateCanonicalRelPath(tt.sourceFile, tt.targetFile)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGenerateHeadingSlug_GFM(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "1. Executive Summary & Strategic Objective",
			expected: "1-executive-summary-strategic-objective",
		},
		{
			input:    "Streaming AST IPC Bridge (V1.3.0)",
			expected: "streaming-ast-ipc-bridge-v130",
		},
		{
			input:    "Heading: With Colons, Brackets [Test], and Symbols!",
			expected: "heading-with-colons-brackets-test-and-symbols",
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			res := GenerateHeadingSlug(tt.input)
			if res != tt.expected {
				t.Errorf("expected slug %q, got %q", tt.expected, res)
			}
		})
	}
}

func TestTarjanSCCDetector_Cycles(t *testing.T) {
	detector := NewTarjanSCCDetector()

	// 1. Graph with 3-node cycle: A -> B -> C -> A
	graph := &DocGraph{
		Nodes: map[string]DocNode{
			"A.md": {ID: "A.md", FilePath: "A.md"},
			"B.md": {ID: "B.md", FilePath: "B.md"},
			"C.md": {ID: "C.md", FilePath: "C.md"},
			"D.md": {ID: "D.md", FilePath: "D.md"},
		},
		Edges: []DocEdge{
			{SourceID: "A.md", TargetID: "B.md"},
			{SourceID: "B.md", TargetID: "C.md"},
			{SourceID: "C.md", TargetID: "A.md"},
			{SourceID: "C.md", TargetID: "D.md"}, // D is a leaf
		},
	}

	report := detector.FindCycles(graph)
	if report.TotalCycles != 1 {
		t.Fatalf("expected 1 cycle, got %d (%v)", report.TotalCycles, report.Cycles)
	}

	cycle := report.Cycles[0]
	if cycle.Length != 3 {
		t.Errorf("expected cycle length 3, got %d", cycle.Length)
	}

	// 2. Empty graph
	emptyReport := detector.FindCycles(&DocGraph{})
	if emptyReport.TotalCycles != 0 {
		t.Errorf("expected 0 cycles on empty graph, got %d", emptyReport.TotalCycles)
	}
}

func TestDocRemediator_PlanAndApplyRemediation(t *testing.T) {
	tmpDir := t.TempDir()

	// Create structure:
	// README.md -> links to docs/INDEX.md (broken stepping)
	// docs/index.md (target file)
	// docs/specs/spec1.md -> links to ../../README.md (broken stepping: should be ../../README.md)
	// docs/roadmap/road1.md -> links to ../specs/spec1.md

	_ = os.MkdirAll(filepath.Join(tmpDir, "docs", "specs"), 0755)
	_ = os.MkdirAll(filepath.Join(tmpDir, "docs", "roadmap"), 0755)

	readmeContent := `# Root README
[Docs Index](./docs/INDEX.md)
`
	indexContent := `# Docs Index
# Overview
[Spec 1](./specs/spec1.md)
`
	specContent := `# Spec 1
## Section Details
[Roadmap 1](../../docs/roadmap/road1.md)
`
	roadContent := `# Roadmap 1
[Spec 1](../specs/spec1.md#section-details)
`

	_ = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(readmeContent), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "docs", "index.md"), []byte(indexContent), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "docs", "specs", "spec1.md"), []byte(specContent), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "docs", "roadmap", "road1.md"), []byte(roadContent), 0644)

	remediator := NewDefaultDocRemediatorService()
	ctx := context.Background()

	// 1. Heading anchor index test
	anchors, errAnchor := remediator.IndexHeadingAnchors(ctx, tmpDir)
	if errAnchor != nil {
		t.Fatalf("failed indexing heading anchors: %v", errAnchor)
	}
	if len(anchors) < 3 {
		t.Errorf("expected at least 3 anchor tables, got %d", len(anchors))
	}

	// 2. Cycle detection test
	cycleReport, errCycle := remediator.DetectCycles(ctx, tmpDir)
	if errCycle != nil {
		t.Fatalf("failed detecting cycles: %v", errCycle)
	}
	if cycleReport == nil {
		t.Fatalf("expected non-nil cycle report")
	}

	// 3. Plan remediation test
	plan, errPlan := remediator.PlanRemediation(ctx, tmpDir, true)
	if errPlan != nil {
		t.Fatalf("failed planning remediation: %v", errPlan)
	}

	if plan.TotalActions == 0 {
		t.Fatalf("expected remediation actions for broken links, got 0")
	}

	// 4. Apply remediation test
	if errApply := remediator.ApplyRemediation(ctx, plan); errApply != nil {
		t.Fatalf("failed applying remediation: %v", errApply)
	}

	// Verify README was updated
	updatedReadme, _ := os.ReadFile(filepath.Join(tmpDir, "README.md"))
	if !strings.Contains(string(updatedReadme), "./docs/index.md") {
		t.Errorf("expected README.md link to be corrected to ./docs/index.md, got:\n%s", string(updatedReadme))
	}
}

func TestDocRemediator_EdgeCases(t *testing.T) {
	remediator := NewDefaultDocRemediatorService()

	// 1. Nil / Empty plan
	if err := remediator.ApplyRemediation(context.Background(), nil); err != nil {
		t.Errorf("expected nil error for nil plan, got %v", err)
	}

	if err := remediator.ApplyRemediation(context.Background(), &DocRemediationPlan{}); err != nil {
		t.Errorf("expected nil error for empty plan, got %v", err)
	}

	// 2. Context Cancellation in PlanRemediation
	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	_, errCanceled := remediator.PlanRemediation(canceledCtx, ".", false)
	if errCanceled == nil {
		t.Errorf("expected error for canceled context, got nil")
	}

	// 3. Self-Loop Cycle
	graph := &DocGraph{
		Nodes: map[string]DocNode{
			"loop.md": {ID: "loop.md", FilePath: "loop.md"},
		},
		Edges: []DocEdge{
			{SourceID: "loop.md", TargetID: "loop.md"},
		},
	}
	detector := NewTarjanSCCDetector()
	cycleReport := detector.FindCycles(graph)
	if cycleReport.TotalCycles != 1 {
		t.Errorf("expected 1 self-loop cycle, got %d", cycleReport.TotalCycles)
	}
}

func TestGraphifyDocAuditCLI_Flags(t *testing.T) {
	tmpDir := t.TempDir()

	readmeContent := `# Root README
[Docs Index](./docs/INDEX.md)
`
	_ = os.MkdirAll(filepath.Join(tmpDir, "docs"), 0755)
	_ = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte(readmeContent), 0644)
	_ = os.WriteFile(filepath.Join(tmpDir, "docs", "index.md"), []byte("# Index\n"), 0644)

	remediator := NewDefaultDocRemediatorService()
	ctx := context.Background()

	// Dry run plan
	plan, err := remediator.PlanRemediation(ctx, tmpDir, true)
	if err != nil || plan.TotalActions == 0 {
		t.Fatalf("expected dry-run plan actions, got: %v, actions: %d", err, plan.TotalActions)
	}

	// Real fix
	if err := remediator.ApplyRemediation(ctx, plan); err != nil {
		t.Fatalf("failed applying remediation: %v", err)
	}

	// Cycles and anchors
	cycles, errCycles := remediator.DetectCycles(ctx, tmpDir)
	if errCycles != nil || cycles == nil {
		t.Fatalf("detect cycles failed: %v", errCycles)
	}

	anchors, errAnchors := remediator.IndexHeadingAnchors(ctx, tmpDir)
	if errAnchors != nil || len(anchors) == 0 {
		t.Fatalf("index heading anchors failed: %v", errAnchors)
	}
}
