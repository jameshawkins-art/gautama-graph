package scaffold

import (
	"context"
	"io/fs"
	"time"
)

// ScaffoldActionType defines the nature of an individual file or directory mutation.
type ScaffoldActionType string

const (
	ActionCreateFile      ScaffoldActionType = "CREATE_FILE"
	ActionSkipExists      ScaffoldActionType = "SKIP_EXISTS"
	ActionOverwriteForce  ScaffoldActionType = "OVERWRITE_FORCE"
	ActionAppendMerge     ScaffoldActionType = "APPEND_MERGE"
	ActionCreateDirectory ScaffoldActionType = "CREATE_DIRECTORY"
)

// ScaffoldAction represents a discrete file or directory scaffolding operation.
type ScaffoldAction struct {
	RelativePath string             `json:"relative_path"`
	AbsolutePath string             `json:"absolute_path"`
	ActionType   ScaffoldActionType `json:"action_type"`
	FileMode     fs.FileMode        `json:"file_mode"`
	Content      string             `json:"content,omitempty"`
	Reason       string             `json:"reason"`
	ExistedPrior bool               `json:"existed_prior"`
	Applied      bool               `json:"applied"`
}

// ScaffoldOptions specifies configuration flags for the scaffolding process.
type ScaffoldOptions struct {
	WorkspaceRoot string `json:"workspace_root"`
	Force         bool   `json:"force"`
	DryRun        bool   `json:"dry_run"`
	Minimal       bool   `json:"minimal"`
	WithScripts   bool   `json:"with_scripts"`
	WithGitIgnore bool   `json:"with_gitignore"`
	Verbose       bool   `json:"verbose"`
}

// ScaffoldPlan aggregates all planned actions before execution.
type ScaffoldPlan struct {
	Timestamp       time.Time        `json:"timestamp"`
	WorkspaceRoot   string           `json:"workspace_root"`
	DryRun          bool             `json:"dry_run"`
	TotalActions    int              `json:"total_actions"`
	CreatedFiles    int              `json:"created_files"`
	ModifiedFiles   int              `json:"modified_files"`
	SkippedFiles    int              `json:"skipped_files"`
	Actions         []ScaffoldAction `json:"actions"`
	ExecutionTimeMs float64          `json:"execution_time_ms"`
}

// ScaffoldVerificationReport summarizes post-installation health checks.
type ScaffoldVerificationReport struct {
	Timestamp     time.Time `json:"timestamp"`
	WorkspaceRoot string    `json:"workspace_root"`
	AllValid      bool      `json:"all_valid"`
	RulesFound    bool      `json:"rules_found"`
	WorkflowFound bool      `json:"workflow_found"`
	ScriptFound   bool      `json:"script_found"`
	Errors        []string  `json:"errors,omitempty"`
}

// ScaffolderService orchestrates workspace inspection, plan formulation, template rendering, and file generation.
type ScaffolderService interface {
	Plan(ctx context.Context, opts ScaffoldOptions) (*ScaffoldPlan, error)
	Execute(ctx context.Context, plan *ScaffoldPlan) error
	Verify(ctx context.Context, workspaceRoot string) (*ScaffoldVerificationReport, error)
}
