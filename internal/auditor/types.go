package auditor

import (
	"context"
	"go/ast"
	"go/token"
	"go/types"
	"time"
)

// ProvenanceStatus defines the verified authenticity status of a graph relationship edge.
type ProvenanceStatus string

const (
	// ProvenanceExtractedAST indicates direct single-file AST call extraction.
	ProvenanceExtractedAST ProvenanceStatus = "EXTRACTED_AST"
	// ProvenanceInferredHeuristic indicates unverified or heuristic extraction.
	ProvenanceInferredHeuristic ProvenanceStatus = "INFERRED_HEURISTIC"
	// ProvenancePrunedPhantom indicates an edge proven false by AST analysis.
	ProvenancePrunedPhantom ProvenanceStatus = "PRUNED_PHANTOM"
	// ProvenanceResolvedCrossPackage indicates a verified cross-package call expression.
	ProvenanceResolvedCrossPackage ProvenanceStatus = "RESOLVED_CROSS_PACKAGE_CALL"
	// ProvenanceResolvedInterfaceImpl indicates verified implicit interface implementation.
	ProvenanceResolvedInterfaceImpl ProvenanceStatus = "RESOLVED_INTERFACE_IMPL"
)

// ExportedKind categorizes Go declaration types.
type ExportedKind string

const (
	KindFunction  ExportedKind = "FUNCTION"
	KindMethod    ExportedKind = "METHOD"
	KindStruct    ExportedKind = "STRUCT"
	KindInterface ExportedKind = "INTERFACE"
	KindTypeAlias ExportedKind = "TYPE_ALIAS"
	KindConstant  ExportedKind = "CONSTANT"
	KindVariable  ExportedKind = "VARIABLE"
)

// ExportedSymbol represents an exported declaration in a package.
type ExportedSymbol struct {
	Name        string       `json:"name"`
	Kind        ExportedKind `json:"kind"`
	Receiver    string       `json:"receiver,omitempty"` // For methods: struct receiver type
	PackagePath string       `json:"package_path"`
	FilePath    string       `json:"file_path"`
	LineNumber  int          `json:"line_number"`
	DocSummary  string       `json:"doc_summary,omitempty"`
}

// PackageSymbolTable indexes all declarations within a single Go package compilation unit.
type PackageSymbolTable struct {
	PackageName string                    `json:"package_name"`
	PackagePath string                    `json:"package_path"`
	Directory   string                    `json:"directory"`
	Symbols     map[string]ExportedSymbol `json:"symbols"`     // Key: Symbol Name
	MethodSets  map[string][]string       `json:"method_sets"` // Key: Type Name -> Method Names
	FileSet     *token.FileSet            `json:"-"`
	PackageAST  *ast.Package              `json:"-"`
	TypePackage *types.Package            `json:"-"`
}

// InterfaceBinding maps a concrete struct implementation to an interface definition.
type InterfaceBinding struct {
	InterfacePackage string   `json:"interface_package"`
	InterfaceName    string   `json:"interface_name"`
	ConcretePackage  string   `json:"concrete_package"`
	ConcreteTypeName string   `json:"concrete_type_name"`
	MatchedMethods   []string `json:"matched_methods"`
}

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
	ProvenanceStatus string  `json:"provenance_status"` // e.g. "EXTRACTED_AST", "RESOLVED_CROSS_PACKAGE_CALL", "RESOLVED_INTERFACE_IMPL", "PRUNED_PHANTOM"
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

// SelectorEvaluator inspects Go AST nodes for explicit call and selector expressions within a file.
type SelectorEvaluator interface {
	EvaluateSelector(file *ast.File, callerIdent, selectorIdent string) (bool, string, error)
}

// PackageSymbolIndexer traverses the workspace and indexes all Go package declarations.
type PackageSymbolIndexer interface {
	IndexWorkspace(ctx context.Context, workspaceRoot string) (map[string]*PackageSymbolTable, error)
	GetPackageTable(packagePath string) (*PackageSymbolTable, bool)
	GetAllPackages() map[string]*PackageSymbolTable
}

// CrossPackageEvaluator evaluates selector and call expressions across imported packages.
type CrossPackageEvaluator interface {
	EvaluateCrossPackageCall(ctx context.Context, sourceFile, callerSymbol, targetPkg, targetSymbol string) (bool, ProvenanceStatus, string, error)
	ResolveFileImports(file *ast.File) map[string]string // alias/ident -> packagePath
}

// InterfaceResolver computes and validates interface satisfaction across packages.
type InterfaceResolver interface {
	CheckImplementation(ctx context.Context, concretePkg, concreteType, ifacePkg, ifaceName string) (bool, *InterfaceBinding, error)
	FindImplementations(ctx context.Context, ifacePkg, ifaceName string) ([]InterfaceBinding, error)
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
