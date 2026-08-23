package auditor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestDocGraphAuditor_ParseWorkspaceDocs(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test markdown files
	// docA.md links to docB.md
	docAPath := filepath.Join(tmpDir, "docA.md")
	docBPath := filepath.Join(tmpDir, "docB.md")
	docOrphanPath := filepath.Join(tmpDir, "docOrphan.md")

	errA := os.WriteFile(docAPath, []byte("# Doc A\nSee [Doc B](./docB.md) for details."), 0644)
	errB := os.WriteFile(docBPath, []byte("# Doc B\nBack to [Doc A](./docA.md)."), 0644)
	errOrphan := os.WriteFile(docOrphanPath, []byte("# Orphan Doc\nNo links here."), 0644)

	if errA != nil || errB != nil || errOrphan != nil {
		t.Fatalf("Failed creating temp test markdown files: %v %v %v", errA, errB, errOrphan)
	}

	parser := NewDefaultDocGraphParser(tmpDir)
	nodes, broken, err := parser.ParseWorkspaceDocs(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("ParseWorkspaceDocs returned unexpected error: %v", err)
	}

	if len(broken) > 0 {
		t.Errorf("Expected 0 broken links, got %d: %v", len(broken), broken)
	}

	nodeMap := make(map[string]DocNodeResult)
	for _, n := range nodes {
		nodeMap[n.FilePath] = n
	}

	if len(nodeMap) != 3 {
		t.Fatalf("Expected 3 nodes, got %d", len(nodeMap))
	}

	// docA should have in_degree 1, out_degree 1, not orphan
	docA := nodeMap["docA.md"]
	if docA.InDegree != 1 || docA.OutDegree != 1 || docA.IsOrphan {
		t.Errorf("docA.md status mismatch: %+v", docA)
	}

	// docB should have in_degree 1, out_degree 1, not orphan
	docB := nodeMap["docB.md"]
	if docB.InDegree != 1 || docB.OutDegree != 1 || docB.IsOrphan {
		t.Errorf("docB.md status mismatch: %+v", docB)
	}

	// docOrphan should have in_degree 0, out_degree 0, is_orphan true
	orphan := nodeMap["docOrphan.md"]
	if orphan.InDegree != 0 || orphan.OutDegree != 0 || !orphan.IsOrphan {
		t.Errorf("docOrphan.md status mismatch: %+v", orphan)
	}
}

func TestDocGraphAuditor_BrokenLinkDetection(t *testing.T) {
	tmpDir := t.TempDir()

	docPath := filepath.Join(tmpDir, "doc.md")
	err := os.WriteFile(docPath, []byte("# Doc\nLink to [NonExistent](./missing.md)"), 0644)
	if err != nil {
		t.Fatalf("Failed writing test doc: %v", err)
	}

	parser := NewDefaultDocGraphParser(tmpDir)
	_, broken, err := parser.ParseWorkspaceDocs(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("ParseWorkspaceDocs returned unexpected error: %v", err)
	}

	if len(broken) != 1 {
		t.Fatalf("Expected 1 broken link, got %d", len(broken))
	}

	if broken[0].SourceFile != "doc.md" || broken[0].LinkTarget != "./missing.md" {
		t.Errorf("Broken link details mismatch: %+v", broken[0])
	}
}

func TestDocGraphAuditor_SecurityPathTraversal(t *testing.T) {
	tmpDir := t.TempDir()

	docPath := filepath.Join(tmpDir, "doc.md")
	// Malicious relative link attempting to escape root
	err := os.WriteFile(docPath, []byte("# Malicious Doc\nLink to [Escape](../../../../etc/passwd)"), 0644)
	if err != nil {
		t.Fatalf("Failed writing test doc: %v", err)
	}

	parser := NewDefaultDocGraphParser(tmpDir)
	_, broken, err := parser.ParseWorkspaceDocs(context.Background(), tmpDir)
	if err != nil {
		t.Fatalf("ParseWorkspaceDocs returned unexpected error: %v", err)
	}

	if len(broken) != 1 {
		t.Fatalf("Expected 1 broken security link, got %d", len(broken))
	}

	if broken[0].ErrorReason == "" {
		t.Errorf("Expected path traversal error reason, got empty string")
	}
}

func TestDocGraphAuditor_AuditDocGraph(t *testing.T) {
	tmpDir := t.TempDir()

	docPath := filepath.Join(tmpDir, "README.md")
	err := os.WriteFile(docPath, []byte("# README\nNo links."), 0644)
	if err != nil {
		t.Fatalf("Failed writing README: %v", err)
	}

	auditor := NewDocGraphAuditor(tmpDir)
	report, err := auditor.AuditDocGraph(context.Background())
	if err != nil {
		t.Fatalf("AuditDocGraph returned unexpected error: %v", err)
	}

	if report.TotalDocNodes != 1 || report.OrphanCount != 1 {
		t.Errorf("Audit report mismatch: Total=%d Orphan=%d", report.TotalDocNodes, report.OrphanCount)
	}

	outputJSON := filepath.Join(tmpDir, "graphify-out", "doc_graph_audit.json")
	if _, statErr := os.Stat(outputJSON); os.IsNotExist(statErr) {
		t.Errorf("Expected output JSON %s to exist on disk", outputJSON)
	}
}
