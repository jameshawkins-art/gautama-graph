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

// IPCCommand defines the discrete action requested of the Python daemon worker.
type IPCCommand string

const (
	// CmdPing verifies daemon liveness and measures IPC round-trip latency.
	CmdPing IPCCommand = "PING"
	// CmdAuditPythonCandidates instructs the worker to evaluate candidate edges against Python source AST.
	CmdAuditPythonCandidates IPCCommand = "AUDIT_CANDIDATES"
	// CmdShutdown signals the worker process to terminate cleanly with exit code 0.
	CmdShutdown IPCCommand = "SHUTDOWN"
)

// WorkerState models the lifecycle state machine of an IPC worker daemon.
type WorkerState string

const (
	WorkerStateStarting   WorkerState = "STARTING"
	WorkerStateIdle       WorkerState = "IDLE"
	WorkerStateBusy       WorkerState = "BUSY"
	WorkerStateCrashed    WorkerState = "CRASHED"
	WorkerStateTerminated WorkerState = "TERMINATED"
)

// IPCRequest represents a single NDJSON message sent over stdin to a Python worker.
type IPCRequest struct {
	ID            string          `json:"id"`
	Command       IPCCommand      `json:"command"`
	WorkspaceRoot string          `json:"workspace_root"`
	SourceFile    string          `json:"source_file,omitempty"`
	Candidates    []CandidateEdge `json:"candidates,omitempty"`
	Timestamp     time.Time       `json:"timestamp"`
}

// IPCResponse represents the single NDJSON message received over stdout from a Python worker.
type IPCResponse struct {
	ID           string        `json:"id"`
	Success      bool          `json:"success"`
	Error        string        `json:"error,omitempty"`
	AuditedEdges []AuditedEdge `json:"audited_edges,omitempty"`
	DurationMs   float64       `json:"duration_ms"`
}

// WorkerStats captures telemetry for an individual daemon process.
type WorkerStats struct {
	WorkerID      int         `json:"worker_id"`
	PID           int         `json:"pid"`
	State         WorkerState `json:"state"`
	RequestsTotal int64       `json:"requests_total"`
	ErrorsTotal   int64       `json:"errors_total"`
	RestartsTotal int64       `json:"restarts_total"`
	LastActive    time.Time   `json:"last_active"`
}

// PoolStats provides aggregated operational health metrics across the worker pool.
type PoolStats struct {
	TotalWorkers int           `json:"total_workers"`
	IdleWorkers  int           `json:"idle_workers"`
	BusyWorkers  int           `json:"busy_workers"`
	Workers      []WorkerStats `json:"workers"`
}

// IPCSession encapsulates an active bi-directional streaming pipe session to a child daemon.
type IPCSession interface {
	Send(ctx context.Context, req *IPCRequest) (*IPCResponse, error)
	Ping(ctx context.Context) (time.Duration, error)
	Status() WorkerState
	PID() int
	Close() error
}

// IPCWorkerPool supervises, load-balances, and auto-recovers a pool of persistent Python worker daemons.
type IPCWorkerPool interface {
	AuditPython(ctx context.Context, sourceFile string, candidates []CandidateEdge) ([]AuditedEdge, error)
	SpawnWorkers(ctx context.Context, poolSize int) error
	Stats() PoolStats
	Close() error
}

// GraphStore handles atomic persistence of audited edge metadata into graphify-out/graph.json.
type GraphStore interface {
	SaveAuditedEdges(ctx context.Context, targetPath string, edges []AuditedEdge) error
}

// ASTGraphAuditorService defines the primary entry point for auditing code relationships.
type ASTGraphAuditorService interface {
	AuditGraphFile(ctx context.Context, graphPath string, verbose bool) (*ASTAuditReport, error)
}

// RemediationRule defines the specific heuristic applied to correct a markdown link.
type RemediationRule string

const (
	// RuleFixRelativePath corrects inaccurate relative directory stepping (e.g. ../../docs -> ../docs).
	RuleFixRelativePath RemediationRule = "FIX_RELATIVE_PATH"
	// RuleResolveFuzzyBasename resolves a moved or renamed target file based on unique basename match.
	RuleResolveFuzzyBasename RemediationRule = "RESOLVE_FUZZY_BASENAME"
	// RuleStripInvalidScheme removes unsupported URI schemes (e.g. file:/// -> relative path).
	RuleStripInvalidScheme RemediationRule = "STRIP_INVALID_SCHEME"
	// RuleFixHeadingAnchor corrects or normalizes heading fragment anchors.
	RuleFixHeadingAnchor RemediationRule = "FIX_HEADING_ANCHOR"
)

// RemediationAction captures an individual link rewrite within a markdown document.
type RemediationAction struct {
	SourceFile       string          `json:"source_file"`
	LineNumber       int             `json:"line_number"`
	OriginalLinkText string          `json:"original_link_text"`
	OriginalTarget   string          `json:"original_target"`
	ResolvedTarget   string          `json:"resolved_target"`
	CanonicalRelPath string          `json:"canonical_rel_path"`
	Rule             RemediationRule `json:"rule"`
	Applied          bool            `json:"applied"`
}

// DocRemediationPlan aggregates all planned link rewrites across the workspace.
type DocRemediationPlan struct {
	WorkspaceRoot   string              `json:"workspace_root"`
	Timestamp       time.Time           `json:"timestamp"`
	TotalDocuments  int                 `json:"total_documents"`
	ModifiedDocs    int                 `json:"modified_docs"`
	TotalActions    int                 `json:"total_actions"`
	Actions         []RemediationAction `json:"actions"`
	DryRun          bool                `json:"dry_run"`
	ExecutionTimeMs float64             `json:"execution_time_ms"`
}

// HeadingAnchorTable maps GitHub-Flavored Markdown (GFM) heading anchor slugs to source headings.
type HeadingAnchorTable struct {
	FilePath string            `json:"file_path"`
	Anchors  map[string]string `json:"anchors"` // anchor_slug -> heading_text
}

// DocNode represents a document node in the topology.
type DocNode struct {
	ID       string `json:"id"`
	FilePath string `json:"file_path"`
}

// DocEdge represents a directed link between documents.
type DocEdge struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
}

// DocGraph encapsulates the directed graph of markdown documents and outbound links.
type DocGraph struct {
	Nodes map[string]DocNode `json:"nodes"`
	Edges []DocEdge          `json:"edges"`
}

// CircularCycle represents a closed directed reference loop in the documentation graph.
type CircularCycle struct {
	CycleID  string   `json:"cycle_id"`
	Length   int      `json:"length"`
	DocChain []string `json:"doc_chain"` // e.g. ["docA.md", "docB.md", "docA.md"]
}

// CycleReport aggregates all circular reference cycles detected in the documentation topology.
type CycleReport struct {
	TotalCycles int             `json:"total_cycles"`
	Cycles      []CircularCycle `json:"cycles"`
}

// DocRemediatorService provides automated link calculation, fuzzy resolution, and atomic in-place rewriting.
type DocRemediatorService interface {
	PlanRemediation(ctx context.Context, workspaceRoot string, dryRun bool) (*DocRemediationPlan, error)
	ApplyRemediation(ctx context.Context, plan *DocRemediationPlan) error
	DetectCycles(ctx context.Context, workspaceRoot string) (*CycleReport, error)
	IndexHeadingAnchors(ctx context.Context, workspaceRoot string) (map[string]*HeadingAnchorTable, error)
}

// CycleDetector computes strongly connected components and closed cycles on the doc link graph.
type CycleDetector interface {
	FindCycles(graph *DocGraph) *CycleReport
}
