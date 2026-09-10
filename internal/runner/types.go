package runner

import (
	"context"
	"time"
)

// PlatformTarget represents the resolved host operating system and CPU architecture.
type PlatformTarget struct {
	OS   string `json:"os"`   // "linux", "darwin", "windows"
	Arch string `json:"arch"` // "amd64", "arm64"
}

// ReleaseAsset defines a downloadable release asset from GitHub releases.
type ReleaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
	ChecksumSHA256     string `json:"checksum_sha256,omitempty"`
}

// ReleaseMetadata encapsulates release tag information from GitHub API.
type ReleaseMetadata struct {
	TagName     string                  `json:"tag_name"`
	PublishedAt time.Time               `json:"published_at"`
	Assets      map[string]ReleaseAsset `json:"assets"` // key: "<os>-<arch>" or asset name
}

// RunnerConfig configures binary caching, execution deadlines, and target workspace paths.
type RunnerConfig struct {
	WorkspaceRootPath  string        `json:"workspace_root_path"`
	CacheDirectoryPath string        `json:"cache_directory_path"`
	GitHubRepoOwner    string        `json:"github_repo_owner"`
	GitHubRepoName     string        `json:"github_repo_name"`
	ExecutionTimeout   time.Duration `json:"execution_timeout"`
	StrictAudit        bool          `json:"strict_audit"`
	ForceDownload      bool          `json:"force_download"`
	VerboseLogging     bool          `json:"verbose_logging"`
}

// PipelineStageStatus records execution metrics for an individual pipeline stage.
type PipelineStageStatus struct {
	StageName string        `json:"stage_name"`
	Duration  time.Duration `json:"duration"`
	Success   bool          `json:"success"`
	Error     string        `json:"error,omitempty"`
}

// PipelineReport summarizes the complete multi-stage execution pass.
type PipelineReport struct {
	Timestamp      time.Time             `json:"timestamp"`
	TotalDuration  time.Duration         `json:"total_duration"`
	BinarySource   string                `json:"binary_source"` // "CACHED", "DOWNLOADED_GITHUB"
	BinaryVersion  string                `json:"binary_version"`
	Stages         []PipelineStageStatus `json:"stages"`
	GraphNodeCount int                   `json:"graph_node_count"`
	GraphEdgeCount int                   `json:"graph_edge_count"`
	PrunedPhantoms int                   `json:"pruned_phantoms"`
	DocOrphanCount int                   `json:"doc_orphan_count"`
	BrokenDocLinks int                   `json:"broken_doc_links"`
}

// ReleaseDownloader abstracts fetching release metadata and binaries from GitHub.
type ReleaseDownloader interface {
	// GetLatestRelease fetches metadata for the latest release on GitHub.
	GetLatestRelease(ctx context.Context, owner, repo string) (*ReleaseMetadata, error)

	// DownloadBinary downloads and caches the target binary for the given platform target.
	DownloadBinary(ctx context.Context, asset ReleaseAsset, destinationPath string) error

	// VerifyChecksum validates a downloaded binary file against expected SHA-256 hash.
	VerifyChecksum(filePath, expectedSHA256 string) (bool, error)
}

// BinaryManager manages local binary resolution, cache validation, and executable permissions.
type BinaryManager interface {
	// EnsureBinary ensures a valid Graphify binary exists locally, downloading if necessary.
	EnsureBinary(ctx context.Context, cfg RunnerConfig) (string, string, error)
}

// SubprocessRunner orchestrates headless Graphify subprocess commands with stream isolation.
type SubprocessRunner interface {
	// ExecuteCommand runs a subcommand on the Graphify binary within the target workspace.
	ExecuteCommand(ctx context.Context, binaryPath, workspaceRoot string, args ...string) ([]byte, []byte, error)
}

// OrchestratorService defines the primary end-to-end multi-stage pipeline coordinator.
type OrchestratorService interface {
	// RunPipeline coordinates Download -> Base Extraction -> AST Audit -> Doc Graph Audit.
	RunPipeline(ctx context.Context, cfg RunnerConfig) (*PipelineReport, error)
}

// QueryOptions configures semantic and question-based graph traversal queries.
type QueryOptions struct {
	// WorkspaceRoot specifies the project root directory (defaults to current working directory).
	WorkspaceRoot string `json:"workspace_root"`

	// GraphPath overrides the default graph location (defaults to graphify-out/graph.json).
	GraphPath string `json:"graph_path,omitempty"`

	// DFS enables depth-first search traversal instead of breadth-first search.
	DFS bool `json:"dfs"`

	// Context specifies repeatable edge-context filter tags.
	Context []string `json:"context,omitempty"`

	// BudgetTokens caps the maximum token output volume (defaults to 2000).
	BudgetTokens int `json:"budget_tokens"`

	// JSONOutput instructs the query engine to emit raw JSON data.
	JSONOutput bool `json:"json_output"`
}

// PathOptions configures shortest-path topological graph queries between two nodes.
type PathOptions struct {
	// WorkspaceRoot specifies the project root directory (defaults to current working directory).
	WorkspaceRoot string `json:"workspace_root"`

	// GraphPath overrides the default graph location (defaults to graphify-out/graph.json).
	GraphPath string `json:"graph_path,omitempty"`
}

// ExplainOptions configures plain-language explanation queries for a target node and neighborhood.
type ExplainOptions struct {
	// WorkspaceRoot specifies the project root directory (defaults to current working directory).
	WorkspaceRoot string `json:"workspace_root"`

	// GraphPath overrides the default graph location (defaults to graphify-out/graph.json).
	GraphPath string `json:"graph_path,omitempty"`
}

// QueryResult encapsulates the raw output, execution duration, and metadata for a completed query.
type QueryResult struct {
	// Output contains the stdout string produced by the graph traversal.
	Output string `json:"output"`

	// BinarySource indicates how the runner resolved the binary ("cached", "system-path", "downloaded").
	BinarySource string `json:"binary_source"`

	// BinaryVersion indicates the release tag or version identifier of the executed binary.
	BinaryVersion string `json:"binary_version"`

	// Duration records total execution time.
	Duration time.Duration `json:"duration"`

	// WorkspaceRoot is the canonical workspace path against which the query was resolved.
	WorkspaceRoot string `json:"workspace_root"`
}

// QueryService coordinates knowledge graph querying via the encapsulated Graphify binary.
type QueryService interface {
	// Query executes a BFS or DFS question traversal across the knowledge graph.
	Query(ctx context.Context, question string, opts QueryOptions) (*QueryResult, error)

	// Path finds the shortest topological path between two nodes in the knowledge graph.
	Path(ctx context.Context, nodeA, nodeB string, opts PathOptions) (*QueryResult, error)

	// Explain generates a plain-language explanation of a node and its adjacent neighborhood.
	Explain(ctx context.Context, concept string, opts ExplainOptions) (*QueryResult, error)
}

