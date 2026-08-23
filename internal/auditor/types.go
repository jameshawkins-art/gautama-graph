package auditor

import (
	"context"
	"go/ast"
	"go/token"
	"time"
)

// CandidateEdge defines an unverified heuristic or inferred relationship candidate.
type CandidateEdge struct {
	ID              string `json:"id"`
	SourceFile      string `json:"source_file"`
	SourceSymbol    string `json:"source_symbol"`
	TargetSymbol    string `json:"target_symbol"`
	InitialRelation string `json:"initial_relation"`
}

// AuditedEdge defines a validated graph relationship annotated with AST provenance.
type AuditedEdge struct {
	CandidateEdge
	ProvenanceStatus string  `json:"provenance_status"` // "EXTRACTED_AST", "INFERRED_HEURISTIC", "PRUNED_PHANTOM"
	Confidence       float64 `json:"confidence"`        // 0.0 to 1.0
	ASTNodePattern   string  `json:"ast_node_pattern,omitempty"`
}

// ASTAuditReport summarizes an AST code relationship audit pass.
type ASTAuditReport struct {
	Timestamp          time.Time     `json:"timestamp"`
	TotalEdges         int           `json:"total_edges"`
	VerifiedASTCount   int           `json:"verified_ast_count"`
	PrunedPhantomCount int           `json:"pruned_phantom_count"`
	HeuristicCount     int           `json:"heuristic_count"`
	Duration           time.Duration `json:"duration"`
	AuditedEdges       []AuditedEdge `json:"audited_edges"`
}

// GraphData models the complete graphify-out/graph.json document preserving all node/link attributes.
type GraphData struct {
	Nodes []map[string]interface{} `json:"nodes"`
	Links []map[string]interface{} `json:"links"`
}

// Config configures auditor execution boundaries and depth thresholds.
type Config struct {
	MaxASTDepth       int           `json:"max_ast_depth"`
	MinConfidence     float64       `json:"min_confidence"`
	AuditorTimeout    time.Duration `json:"auditor_timeout"`
	WorkspaceRootPath string        `json:"workspace_root_path"`
}

// ASTParser abstracts Go file AST parsing with position tracking.
type ASTParser interface {
	ParseFile(ctx context.Context, filePath string) (*ast.File, *token.FileSet, error)
}

// SelectorEvaluator inspects Go AST nodes for explicit call and selector expressions.
type SelectorEvaluator interface {
	EvaluateSelector(file *ast.File, callerIdent, selectorIdent string) (bool, string, error)
}

// PythonASTBridge delegates Python file AST call analysis to python/ast_auditor_bridge.py.
type PythonASTBridge interface {
	AuditPythonCandidates(ctx context.Context, targetFile string, candidates []CandidateEdge) ([]AuditedEdge, error)
}

// GraphStore handles atomic persistence of audited edge metadata into graphify-out/graph.json.
type GraphStore interface {
	SaveAuditedEdges(ctx context.Context, targetPath string, edges []AuditedEdge) error
}

// ASTGraphAuditorService defines the primary entry point for auditing code relationships.
type ASTGraphAuditorService interface {
	AuditGraphFile(ctx context.Context, graphPath string, verbose bool) (*ASTAuditReport, error)
}
