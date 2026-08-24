package auditor

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func createSyntheticWorkspace(t *testing.T) string {
	tmpDir := t.TempDir()

	// Package A (pkgA) with interface definition and constant/var
	pkgADir := filepath.Join(tmpDir, "pkgA")
	_ = os.MkdirAll(pkgADir, 0755)
	pkgACode := `package pkgA

const DefaultTimeout = 30
var ActiveCount = 1

type Service interface {
	DoWork() error
	GetStatus() string
}

func HelperA() string {
	return "helperA"
}
`
	_ = os.WriteFile(filepath.Join(pkgADir, "service.go"), []byte(pkgACode), 0644)

	// Package B (pkgB) with concrete struct implementing Service and calling HelperA
	pkgBDir := filepath.Join(tmpDir, "pkgB")
	_ = os.MkdirAll(pkgBDir, 0755)
	pkgBCode := `package pkgB

import (
	"pkgA"
	aliasA "pkgA"
)

type ConcreteService struct {
	ID string
}

func (s *ConcreteService) DoWork() error {
	_ = pkgA.HelperA()
	_ = aliasA.HelperA()
	return nil
}

func (s *ConcreteService) GetStatus() string {
	return "active"
}

func NewConcreteService() *ConcreteService {
	return &ConcreteService{ID: "test"}
}
`
	_ = os.WriteFile(filepath.Join(pkgBDir, "impl.go"), []byte(pkgBCode), 0644)

	// Package C (pkgC) with dot import
	pkgCDir := filepath.Join(tmpDir, "pkgC")
	_ = os.MkdirAll(pkgCDir, 0755)
	pkgCCode := `package pkgC

import (
	. "pkgA"
)

func CallWithDot() string {
	return HelperA()
}
`
	_ = os.WriteFile(filepath.Join(pkgCDir, "dot.go"), []byte(pkgCCode), 0644)

	// Package Incomplete (pkgIncomplete) with missing method
	pkgIncDir := filepath.Join(tmpDir, "pkgIncomplete")
	_ = os.MkdirAll(pkgIncDir, 0755)
	pkgIncCode := `package pkgIncomplete

type IncompleteService struct{}

func (s *IncompleteService) DoWork() error {
	return nil
}
`
	_ = os.WriteFile(filepath.Join(pkgIncDir, "inc.go"), []byte(pkgIncCode), 0644)

	return tmpDir
}

func TestPackageSymbolIndexer_IndexWorkspace(t *testing.T) {
	wsDir := createSyntheticWorkspace(t)

	indexer := NewDefaultPackageSymbolIndexer()
	tables, err := indexer.IndexWorkspace(context.Background(), wsDir)
	if err != nil {
		t.Fatalf("unexpected error indexing workspace: %v", err)
	}

	if len(tables) < 4 {
		t.Errorf("expected at least 4 packages, got %d", len(tables))
	}

	// Verify pkgA symbols
	tableA, foundA := indexer.GetPackageTable("pkgA")
	if !foundA || tableA == nil {
		t.Fatalf("expected pkgA table to be found")
	}

	if _, exists := tableA.Symbols["Service"]; !exists {
		t.Errorf("expected Service interface in pkgA symbols")
	}

	if _, exists := tableA.Symbols["HelperA"]; !exists {
		t.Errorf("expected HelperA func in pkgA symbols")
	}

	if _, exists := tableA.Symbols["DefaultTimeout"]; !exists {
		t.Errorf("expected DefaultTimeout const in pkgA symbols")
	}

	// Verify pkgB symbols and method sets
	tableB, foundB := indexer.GetPackageTable("pkgB")
	if !foundB || tableB == nil {
		t.Fatalf("expected pkgB table to be found")
	}

	if _, exists := tableB.Symbols["ConcreteService"]; !exists {
		t.Errorf("expected ConcreteService struct in pkgB symbols")
	}

	methods := tableB.MethodSets["ConcreteService"]
	if len(methods) != 2 {
		t.Errorf("expected 2 methods on ConcreteService, got %d (%v)", len(methods), methods)
	}

	// Verify GetAllPackages
	all := indexer.GetAllPackages()
	if len(all) == 0 {
		t.Errorf("expected non-empty GetAllPackages")
	}
}

func TestCrossPackageEvaluator_EvaluateCrossPackageCall(t *testing.T) {
	wsDir := createSyntheticWorkspace(t)

	indexer := NewDefaultPackageSymbolIndexer()
	_, _ = indexer.IndexWorkspace(context.Background(), wsDir)

	evaluator := NewDefaultCrossPackageEvaluator(indexer)
	srcFile := filepath.Join(wsDir, "pkgB", "impl.go")

	// Test cross package call: ConcreteService.DoWork -> pkgA.HelperA
	matched, prov, pattern, err := evaluator.EvaluateCrossPackageCall(context.Background(), srcFile, "DoWork", "pkgA", "HelperA")
	if err != nil {
		t.Fatalf("unexpected error in cross package call eval: %v", err)
	}

	if !matched {
		t.Errorf("expected cross package call to match")
	}

	if prov != ProvenanceResolvedCrossPackage {
		t.Errorf("expected provenance %s, got %s", ProvenanceResolvedCrossPackage, prov)
	}

	if pattern == "" {
		t.Errorf("expected non-empty AST pattern")
	}

	// Test dot import call in pkgC
	dotFile := filepath.Join(wsDir, "pkgC", "dot.go")
	dotMatched, dotProv, _, dotErr := evaluator.EvaluateCrossPackageCall(context.Background(), dotFile, "CallWithDot", "pkgA", "HelperA")
	if dotErr != nil || !dotMatched {
		t.Errorf("expected dot import call to match, got %v, err: %v", dotMatched, dotErr)
	}
	if dotProv != ProvenanceResolvedCrossPackage {
		t.Errorf("expected provenance %s, got %s", ProvenanceResolvedCrossPackage, dotProv)
	}

	// Test missing target symbol
	matchedMissing, _, _, _ := evaluator.EvaluateCrossPackageCall(context.Background(), srcFile, "DoWork", "pkgA", "NonExistentFunc")
	if matchedMissing {
		t.Errorf("expected false for non-existent target symbol")
	}

	// Test missing file error
	_, _, _, errMissingFile := evaluator.EvaluateCrossPackageCall(context.Background(), filepath.Join(wsDir, "missing.go"), "DoWork", "pkgA", "HelperA")
	if errMissingFile == nil {
		t.Errorf("expected error for missing file, got nil")
	}
}

func TestInterfaceResolver_CheckImplementation_And_Find(t *testing.T) {
	wsDir := createSyntheticWorkspace(t)

	indexer := NewDefaultPackageSymbolIndexer()
	_, _ = indexer.IndexWorkspace(context.Background(), wsDir)

	resolver := NewDefaultInterfaceResolver(indexer)

	// 1. Success case: ConcreteService in pkgB satisfies Service in pkgA
	isImpl, binding, err := resolver.CheckImplementation(context.Background(), "pkgB", "ConcreteService", "pkgA", "Service")
	if err != nil {
		t.Fatalf("unexpected error checking interface impl: %v", err)
	}

	if !isImpl {
		t.Errorf("expected ConcreteService to implement Service")
	}

	if binding == nil || len(binding.MatchedMethods) != 2 {
		t.Errorf("expected binding with 2 matched methods, got %v", binding)
	}

	// 2. Incomplete case: IncompleteService in pkgIncomplete lacks GetStatus
	isImplInc, _, errInc := resolver.CheckImplementation(context.Background(), "pkgIncomplete", "IncompleteService", "pkgA", "Service")
	if errInc != nil {
		t.Fatalf("unexpected error on incomplete service check: %v", errInc)
	}
	if isImplInc {
		t.Errorf("expected IncompleteService to NOT implement Service")
	}

	// 3. FindImplementations test
	foundImpls, findErr := resolver.FindImplementations(context.Background(), "pkgA", "Service")
	if findErr != nil {
		t.Fatalf("unexpected error finding implementations: %v", findErr)
	}
	if len(foundImpls) == 0 {
		t.Errorf("expected at least 1 implementation found, got %d", len(foundImpls))
	}

	// 4. Error cases in CheckImplementation
	_, _, errMissingIface := resolver.CheckImplementation(context.Background(), "pkgB", "ConcreteService", "pkgMissing", "Service")
	if errMissingIface == nil {
		t.Errorf("expected error for missing interface package, got nil")
	}

	_, _, errMissingConcrete := resolver.CheckImplementation(context.Background(), "pkgMissing", "MissingService", "pkgA", "Service")
	if errMissingConcrete == nil {
		t.Errorf("expected error for missing concrete package, got nil")
	}

	_, _, errNonIface := resolver.CheckImplementation(context.Background(), "pkgB", "ConcreteService", "pkgA", "HelperA")
	if errNonIface == nil {
		t.Errorf("expected error when symbol is not an interface, got nil")
	}

	_, _, errNonStruct := resolver.CheckImplementation(context.Background(), "pkgA", "HelperA", "pkgA", "Service")
	if errNonStruct == nil {
		t.Errorf("expected error when concrete symbol is not a struct, got nil")
	}
}

func TestEngine_MultiPackage_AuditGraphFile(t *testing.T) {
	wsDir := createSyntheticWorkspace(t)
	graphPath := filepath.Join(wsDir, "graph.json")

	graphJSON := `{
  "nodes": [
    {"id": "pkgB.DoWork", "label": "DoWork", "source_file": "pkgB/impl.go"},
    {"id": "pkgA.HelperA", "label": "HelperA", "source_file": "pkgA/service.go"},
    {"id": "pkgB.ConcreteService", "label": "ConcreteService", "source_file": "pkgB/impl.go"},
    {"id": "pkgA.Service", "label": "Service", "source_file": "pkgA/service.go"},
    {"id": "pkgB.Unknown", "label": "Unknown", "source_file": "pkgB/impl.go"}
  ],
  "links": [
    {"source": "pkgB.DoWork", "target": "pkgA.HelperA", "relation": "calls"},
    {"source": "pkgB.ConcreteService", "target": "pkgA.Service", "relation": "implements"},
    {"source": "pkgB.DoWork", "target": "pkgB.Unknown", "relation": "calls"}
  ]
}`
	_ = os.WriteFile(graphPath, []byte(graphJSON), 0644)

	cfg := Config{
		WorkspaceRootPath: wsDir,
		AuditorTimeout:    10 * time.Second,
	}
	engine := NewDefaultEngine(cfg)

	report, err := engine.AuditGraphFile(context.Background(), graphPath, false)
	if err != nil {
		t.Fatalf("AuditGraphFile failed: %v", err)
	}

	if report.TotalEdges != 3 {
		t.Errorf("expected 3 total edges, got %d", report.TotalEdges)
	}

	if report.VerifiedASTCount != 2 {
		t.Errorf("expected 2 verified AST edges (1 cross-call + 1 iface-impl), got %d", report.VerifiedASTCount)
	}

	if report.PrunedPhantomCount != 1 {
		t.Errorf("expected 1 pruned phantom edge, got %d", report.PrunedPhantomCount)
	}
}

func TestEngine_MultiPackage_AuditCandidates(t *testing.T) {
	wsDir := createSyntheticWorkspace(t)

	cfg := Config{
		WorkspaceRootPath: wsDir,
		AuditorTimeout:    10 * time.Second,
	}
	engine := NewDefaultEngine(cfg)

	// Setters test
	engine.SetPythonBridge(NewDefaultPythonASTBridge(wsDir))
	engine.SetIndexer(NewDefaultPackageSymbolIndexer())
	engine.SetCrossEvaluator(NewDefaultCrossPackageEvaluator(engine.indexer))
	engine.SetInterfaceResolver(NewDefaultInterfaceResolver(engine.indexer))

	// Pre-index workspace
	_, _ = engine.indexer.IndexWorkspace(context.Background(), wsDir)

	candidates := []CandidateEdge{
		{
			ID:           "c1",
			SourceFile:   filepath.Join("pkgB", "impl.go"),
			SourceSymbol: "DoWork",
			TargetSymbol: "HelperA",
		},
		{
			ID:           "c2",
			SourceFile:   filepath.Join("pkgB", "impl.go"),
			SourceSymbol: "ConcreteService",
			TargetSymbol: "Service",
		},
		{
			ID:           "c3",
			SourceFile:   filepath.Join("pkgB", "impl.go"),
			SourceSymbol: "DoWork",
			TargetSymbol: "NonExistentMethod",
		},
	}

	results, err := engine.AuditCandidates(context.Background(), candidates)
	if err != nil {
		t.Fatalf("audit candidates failed: %v", err)
	}

	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// Verify c1 -> RESOLVED_CROSS_PACKAGE_CALL
	if results[0].ProvenanceStatus != string(ProvenanceResolvedCrossPackage) {
		t.Errorf("expected c1 provenance %s, got %s", ProvenanceResolvedCrossPackage, results[0].ProvenanceStatus)
	}

	// Verify c2 -> RESOLVED_INTERFACE_IMPL
	if results[1].ProvenanceStatus != string(ProvenanceResolvedInterfaceImpl) {
		t.Errorf("expected c2 provenance %s, got %s", ProvenanceResolvedInterfaceImpl, results[1].ProvenanceStatus)
	}

	// Verify c3 -> PRUNED_PHANTOM
	if results[2].ProvenanceStatus != string(ProvenancePrunedPhantom) {
		t.Errorf("expected c3 provenance %s, got %s", ProvenancePrunedPhantom, results[2].ProvenanceStatus)
	}

	// Test AuditAndSave
	savedResults, errSave := engine.AuditAndSave(context.Background(), filepath.Join(wsDir, "out_test.json"), candidates)
	if errSave != nil || len(savedResults) != 3 {
		t.Errorf("expected AuditAndSave to succeed, got %v, err: %v", savedResults, errSave)
	}
}
