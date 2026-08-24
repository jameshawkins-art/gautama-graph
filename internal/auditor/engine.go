package auditor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Engine orchestrates AST parsing, selector evaluation, and graph edge metadata updates.
type Engine struct {
	parser         ASTParser
	evaluator      SelectorEvaluator
	indexer        PackageSymbolIndexer
	crossEvaluator CrossPackageEvaluator
	ifaceResolver  InterfaceResolver
	pyBridge       PythonASTBridge
	store          GraphStore
	cfg            Config
}

// NewEngine initializes a new Engine instance.
func NewEngine(p ASTParser, e SelectorEvaluator, s GraphStore, cfg Config) *Engine {
	if cfg.AuditorTimeout <= 0 {
		cfg.AuditorTimeout = 60 * time.Second
	}
	if cfg.MaxASTDepth <= 0 {
		cfg.MaxASTDepth = 50
	}
	idx := NewDefaultPackageSymbolIndexer()
	return &Engine{
		parser:         p,
		evaluator:      e,
		indexer:        idx,
		crossEvaluator: NewDefaultCrossPackageEvaluator(idx),
		ifaceResolver:  NewDefaultInterfaceResolver(idx),
		pyBridge:       NewDefaultPythonASTBridge(cfg.WorkspaceRootPath),
		store:          s,
		cfg:            cfg,
	}
}

// NewDefaultEngine constructs an Engine instance with default standard adapters.
func NewDefaultEngine(cfg Config) *Engine {
	if cfg.AuditorTimeout <= 0 {
		cfg.AuditorTimeout = 60 * time.Second
	}
	if cfg.MaxASTDepth <= 0 {
		cfg.MaxASTDepth = 50
	}
	if cfg.MinConfidence <= 0 {
		cfg.MinConfidence = 0.8
	}
	idx := NewDefaultPackageSymbolIndexer()
	return &Engine{
		parser:         NewDefaultASTParser(cfg.WorkspaceRootPath),
		evaluator:      NewDefaultSelectorEvaluator(cfg.MaxASTDepth),
		indexer:        idx,
		crossEvaluator: NewDefaultCrossPackageEvaluator(idx),
		ifaceResolver:  NewDefaultInterfaceResolver(idx),
		pyBridge:       NewDefaultPythonASTBridge(cfg.WorkspaceRootPath),
		store:          NewJSONGraphStore(),
		cfg:            cfg,
	}
}

// SetPythonBridge allows overriding the PythonASTBridge adapter (e.g. for testing).
func (e *Engine) SetPythonBridge(b PythonASTBridge) {
	e.pyBridge = b
}

// SetIndexer allows overriding the PackageSymbolIndexer adapter.
func (e *Engine) SetIndexer(idx PackageSymbolIndexer) {
	e.indexer = idx
	e.crossEvaluator = NewDefaultCrossPackageEvaluator(idx)
	e.ifaceResolver = NewDefaultInterfaceResolver(idx)
}

// SetCrossEvaluator allows overriding the CrossPackageEvaluator adapter.
func (e *Engine) SetCrossEvaluator(ce CrossPackageEvaluator) {
	e.crossEvaluator = ce
}

// SetInterfaceResolver allows overriding the InterfaceResolver adapter.
func (e *Engine) SetInterfaceResolver(ir InterfaceResolver) {
	e.ifaceResolver = ir
}

// AuditCandidates audits a list of CandidateEdge candidates against Go and Python AST specifications.
func (e *Engine) AuditCandidates(ctx context.Context, candidates []CandidateEdge) ([]AuditedEdge, error) {
	ctx, cancel := context.WithTimeout(ctx, e.cfg.AuditorTimeout)
	defer cancel()

	results := make([]AuditedEdge, 0, len(candidates))
	for _, cand := range candidates {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if strings.HasSuffix(cand.SourceFile, ".py") {
			if e.pyBridge != nil {
				pyResults, err := e.pyBridge.AuditPythonCandidates(ctx, cand.SourceFile, []CandidateEdge{cand})
				if err == nil && len(pyResults) > 0 {
					results = append(results, pyResults...)
					continue
				}
			}
			results = append(results, AuditedEdge{
				CandidateEdge:    cand,
				ProvenanceStatus: string(ProvenanceInferredHeuristic),
				Confidence:       0.5,
			})
			continue
		}

		absSource := cand.SourceFile
		if !filepath.IsAbs(absSource) && e.cfg.WorkspaceRootPath != "" {
			absSource = filepath.Join(e.cfg.WorkspaceRootPath, cand.SourceFile)
		}

		fileAST, _, err := e.parser.ParseFile(ctx, absSource)
		if err != nil {
			// Fail-safe: mark unparseable files with zero confidence without breaking pipeline
			results = append(results, AuditedEdge{
				CandidateEdge:    cand,
				ProvenanceStatus: string(ProvenancePrunedPhantom),
				Confidence:       0.0,
			})
			continue
		}

		// 1. Interface Implementation Evaluation
		if e.ifaceResolver != nil {
			ifacePkg := filepath.Dir(cand.TargetSymbol)
			ifaceName := filepath.Base(cand.TargetSymbol)
			concretePkg := filepath.Dir(cand.SourceFile)
			concreteType := cand.SourceSymbol

			ifImpl, _, ifErr := e.ifaceResolver.CheckImplementation(ctx, concretePkg, concreteType, ifacePkg, ifaceName)
			if ifErr == nil && ifImpl {
				results = append(results, AuditedEdge{
					CandidateEdge:    cand,
					ProvenanceStatus: string(ProvenanceResolvedInterfaceImpl),
					Confidence:       1.0,
					ASTNodePattern:   fmt.Sprintf("interface_impl: %s satisfies %s", concreteType, ifaceName),
				})
				continue
			}
		}

		// 2. Multi-Package Cross-Call Evaluation
		if e.crossEvaluator != nil {
			targetPkg := filepath.Dir(cand.TargetSymbol)
			targetSym := filepath.Base(cand.TargetSymbol)
			if targetPkg == "." || targetPkg == "" {
				targetPkg = ""
				targetSym = cand.TargetSymbol
			}

			crossMatched, provStatus, crossPattern, crossErr := e.crossEvaluator.EvaluateCrossPackageCall(ctx, absSource, cand.SourceSymbol, targetPkg, targetSym)
			if crossErr == nil && crossMatched {
				results = append(results, AuditedEdge{
					CandidateEdge:    cand,
					ProvenanceStatus: string(provStatus),
					Confidence:       1.0,
					ASTNodePattern:   crossPattern,
				})
				continue
			}
		}

		// 3. Local Single-File AST Evaluation
		matched, pattern, err := e.evaluator.EvaluateSelector(fileAST, cand.SourceSymbol, cand.TargetSymbol)
		if err == nil && matched {
			results = append(results, AuditedEdge{
				CandidateEdge:    cand,
				ProvenanceStatus: string(ProvenanceExtractedAST),
				Confidence:       1.0,
				ASTNodePattern:   pattern,
			})
			continue
		}

		// 4. Default: Pruned Phantom
		results = append(results, AuditedEdge{
			CandidateEdge:    cand,
			ProvenanceStatus: string(ProvenancePrunedPhantom),
			Confidence:       0.0,
		})
	}

	return results, nil
}

// AuditAndSave performs candidate auditing and persists audited edges into targetStorePath.
func (e *Engine) AuditAndSave(ctx context.Context, targetStorePath string, candidates []CandidateEdge) ([]AuditedEdge, error) {
	audited, err := e.AuditCandidates(ctx, candidates)
	if err != nil {
		return nil, fmt.Errorf("audit candidates failed: %w", err)
	}

	if e.store != nil && targetStorePath != "" {
		if saveErr := e.store.SaveAuditedEdges(ctx, targetStorePath, audited); saveErr != nil {
			return audited, fmt.Errorf("saving audited edges failed: %w", saveErr)
		}
	}

	return audited, nil
}

// AuditGraphFile scans graphify-out/graph.json, evaluates Go and Python code edges, prunes phantoms, and updates the file atomically.
func (e *Engine) AuditGraphFile(ctx context.Context, graphPath string, verbose bool) (*ASTAuditReport, error) {
	startTime := time.Now()
	ctx, cancel := context.WithTimeout(ctx, e.cfg.AuditorTimeout)
	defer cancel()

	cleanGraphPath := filepath.Clean(graphPath)
	data, err := os.ReadFile(cleanGraphPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read graph file at %s: %w", cleanGraphPath, err)
	}

	var graph GraphData
	if err := json.Unmarshal(data, &graph); err != nil {
		return nil, fmt.Errorf("failed to unmarshal graph JSON: %w", err)
	}

	report := &ASTAuditReport{
		Timestamp:    time.Now(),
		TotalEdges:   len(graph.Links),
		AuditedEdges: make([]AuditedEdge, 0),
	}

	// Step 1: Pre-index workspace packages for deep cross-package AST analysis
	if e.indexer != nil && e.cfg.WorkspaceRootPath != "" {
		_, _ = e.indexer.IndexWorkspace(ctx, e.cfg.WorkspaceRootPath)
	}

	// Build node lookup map for symbol labels and file paths
	nodeMap := make(map[string]map[string]interface{})
	for _, node := range graph.Nodes {
		if id, ok := node["id"].(string); ok && id != "" {
			nodeMap[id] = node
		}
	}

	// Separate candidate edges by source file for efficient batch analysis
	goCandidates := make(map[string][]CandidateEdge)
	pyCandidates := make(map[string][]CandidateEdge)

	for _, link := range graph.Links {
		sourceID, _ := link["source"].(string)
		targetID, _ := link["target"].(string)
		relation, _ := link["relation"].(string)

		var sourceFile string
		if sf, ok := link["source_file"].(string); ok && sf != "" {
			sourceFile = sf
		} else if srcNode, exists := nodeMap[sourceID]; exists {
			if sf, ok := srcNode["source_file"].(string); ok {
				sourceFile = sf
			} else if sf, ok := srcNode["file"].(string); ok {
				sourceFile = sf
			}
		}

		sourceSymbol := sourceID
		if srcNode, exists := nodeMap[sourceID]; exists {
			if lbl, ok := srcNode["label"].(string); ok && lbl != "" {
				sourceSymbol = lbl
			}
		}

		targetSymbol := targetID
		if tgtNode, exists := nodeMap[targetID]; exists {
			if lbl, ok := tgtNode["label"].(string); ok && lbl != "" {
				targetSymbol = lbl
			}
		}

		cand := CandidateEdge{
			ID:              fmt.Sprintf("%s->%s", sourceID, targetID),
			SourceFile:      sourceFile,
			SourceSymbol:    sourceSymbol,
			TargetSymbol:    targetSymbol,
			InitialRelation: relation,
		}

		if strings.HasSuffix(sourceFile, ".go") {
			goCandidates[sourceFile] = append(goCandidates[sourceFile], cand)
		} else if strings.HasSuffix(sourceFile, ".py") {
			pyCandidates[sourceFile] = append(pyCandidates[sourceFile], cand)
		} else {
			report.HeuristicCount++
		}
	}

	auditedMap := make(map[string]AuditedEdge)

	// Audit Go Candidates with tiered evaluation (interface -> cross-package -> single-file)
	for srcFile, cands := range goCandidates {
		absPath := srcFile
		if !filepath.IsAbs(absPath) && e.cfg.WorkspaceRootPath != "" {
			absPath = filepath.Join(e.cfg.WorkspaceRootPath, srcFile)
		}

		fileAST, _, parseErr := e.parser.ParseFile(ctx, absPath)
		for _, c := range cands {
			if parseErr != nil {
				audited := AuditedEdge{
					CandidateEdge:    c,
					ProvenanceStatus: string(ProvenancePrunedPhantom),
					Confidence:       0.0,
				}
				auditedMap[c.ID] = audited
				report.PrunedPhantomCount++
				report.AuditedEdges = append(report.AuditedEdges, audited)
				continue
			}

			// 1. Interface Implementation Evaluation
			if e.ifaceResolver != nil {
				ifacePkg := filepath.Dir(c.TargetSymbol)
				ifaceName := filepath.Base(c.TargetSymbol)
				concretePkg := filepath.Dir(c.SourceFile)
				concreteType := c.SourceSymbol

				ifImpl, _, ifErr := e.ifaceResolver.CheckImplementation(ctx, concretePkg, concreteType, ifacePkg, ifaceName)
				if ifErr == nil && ifImpl {
					audited := AuditedEdge{
						CandidateEdge:    c,
						ProvenanceStatus: string(ProvenanceResolvedInterfaceImpl),
						Confidence:       1.0,
						ASTNodePattern:   fmt.Sprintf("interface_impl: %s satisfies %s", concreteType, ifaceName),
					}
					auditedMap[c.ID] = audited
					report.VerifiedASTCount++
					report.AuditedEdges = append(report.AuditedEdges, audited)
					continue
				}
			}

			// 2. Multi-Package Cross-Call Evaluation
			if e.crossEvaluator != nil {
				targetPkg := filepath.Dir(c.TargetSymbol)
				targetSym := filepath.Base(c.TargetSymbol)
				if targetPkg == "." || targetPkg == "" {
					targetPkg = ""
					targetSym = c.TargetSymbol
				}

				crossMatched, provStatus, crossPattern, crossErr := e.crossEvaluator.EvaluateCrossPackageCall(ctx, absPath, c.SourceSymbol, targetPkg, targetSym)
				if crossErr == nil && crossMatched {
					audited := AuditedEdge{
						CandidateEdge:    c,
						ProvenanceStatus: string(provStatus),
						Confidence:       1.0,
						ASTNodePattern:   crossPattern,
					}
					auditedMap[c.ID] = audited
					report.VerifiedASTCount++
					report.AuditedEdges = append(report.AuditedEdges, audited)
					continue
				}
			}

			// 3. Local Single-File Selector Check
			matched, pattern, evalErr := e.evaluator.EvaluateSelector(fileAST, c.SourceSymbol, c.TargetSymbol)
			if evalErr == nil && matched {
				audited := AuditedEdge{
					CandidateEdge:    c,
					ProvenanceStatus: string(ProvenanceExtractedAST),
					Confidence:       1.0,
					ASTNodePattern:   pattern,
				}
				auditedMap[c.ID] = audited
				report.VerifiedASTCount++
				report.AuditedEdges = append(report.AuditedEdges, audited)
				continue
			}

			// 4. Default: Pruned Phantom
			audited := AuditedEdge{
				CandidateEdge:    c,
				ProvenanceStatus: string(ProvenancePrunedPhantom),
				Confidence:       0.0,
			}
			auditedMap[c.ID] = audited
			report.PrunedPhantomCount++
			report.AuditedEdges = append(report.AuditedEdges, audited)
		}
	}

	// Audit Python Candidates via Subprocess Bridge
	for srcFile, cands := range pyCandidates {
		absPath := srcFile
		if !filepath.IsAbs(absPath) && e.cfg.WorkspaceRootPath != "" {
			absPath = filepath.Join(e.cfg.WorkspaceRootPath, srcFile)
		}

		var pyAudited []AuditedEdge
		var pyErr error
		if e.pyBridge != nil {
			pyAudited, pyErr = e.pyBridge.AuditPythonCandidates(ctx, absPath, cands)
		} else {
			pyErr = fmt.Errorf("python bridge not configured")
		}

		if pyErr != nil {
			for _, c := range cands {
				audited := AuditedEdge{
					CandidateEdge:    c,
					ProvenanceStatus: string(ProvenanceInferredHeuristic),
					Confidence:       0.5,
				}
				auditedMap[c.ID] = audited
				report.HeuristicCount++
				report.AuditedEdges = append(report.AuditedEdges, audited)
			}
			continue
		}

		for _, a := range pyAudited {
			auditedMap[a.ID] = a
			report.AuditedEdges = append(report.AuditedEdges, a)
			switch a.ProvenanceStatus {
			case string(ProvenanceExtractedAST):
				report.VerifiedASTCount++
			case string(ProvenancePrunedPhantom):
				report.PrunedPhantomCount++
			default:
				report.HeuristicCount++
			}
		}
	}

	// Update links in GraphData and prune phantom edges (confidence == 0.0)
	updatedLinks := make([]map[string]interface{}, 0, len(graph.Links))
	for _, link := range graph.Links {
		sourceID, _ := link["source"].(string)
		targetID, _ := link["target"].(string)
		linkID := fmt.Sprintf("%s->%s", sourceID, targetID)

		if audited, exists := auditedMap[linkID]; exists {
			link["provenance"] = audited.ProvenanceStatus
			link["confidence"] = audited.ProvenanceStatus
			link["confidence_score"] = audited.Confidence
			if audited.ASTNodePattern != "" {
				link["ast_pattern"] = audited.ASTNodePattern
			}

			// Only keep edges that were verified or heuristic; prune phantoms (confidence 0.0)
			if audited.Confidence > 0.0 {
				updatedLinks = append(updatedLinks, link)
			}
		} else {
			updatedLinks = append(updatedLinks, link)
		}
	}

	graph.Links = updatedLinks

	// Save updated graph data atomically
	if jsStore, ok := e.store.(*JSONGraphStore); ok {
		if err := jsStore.SaveGraphData(ctx, cleanGraphPath, &graph); err != nil {
			return nil, fmt.Errorf("failed to save audited graph: %w", err)
		}
	} else if e.store != nil {
		if err := e.store.SaveAuditedEdges(ctx, cleanGraphPath, report.AuditedEdges); err != nil {
			return nil, fmt.Errorf("failed to save audited edges: %w", err)
		}
	}

	report.Duration = time.Since(startTime)
	return report, nil
}
