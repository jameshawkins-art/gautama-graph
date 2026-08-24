package scaffold

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jameshawkins-art/gautama-graph/internal/auditor"
)

//go:embed templates/*
var templatesFS embed.FS

// DefaultScaffolderService implements ScaffolderService for turnkey environment initialization.
type DefaultScaffolderService struct {
	fs embed.FS
}

// NewDefaultScaffolderService constructs a new DefaultScaffolderService instance with embedded templates.
func NewDefaultScaffolderService() *DefaultScaffolderService {
	return &DefaultScaffolderService{
		fs: templatesFS,
	}
}

// Plan inspects the target workspace and constructs a non-destructive scaffold execution plan.
func (s *DefaultScaffolderService) Plan(ctx context.Context, opts ScaffoldOptions) (*ScaffoldPlan, error) {
	cleanRoot := filepath.Clean(opts.WorkspaceRoot)
	if cleanRoot == "" || cleanRoot == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed getting working directory: %w", err)
		}
		cleanRoot = cwd
	}

	startTime := time.Now()
	plan := &ScaffoldPlan{
		Timestamp:     startTime,
		WorkspaceRoot: cleanRoot,
		DryRun:        opts.DryRun,
		Actions:       make([]ScaffoldAction, 0),
	}

	type scaffoldItem struct {
		relPath      string
		templatePath string
		mode         fs.FileMode
		isMerge      bool
		mergeMarker  string
	}

	items := []scaffoldItem{
		{
			relPath:      filepath.Join(".agents", "rules", "graphify.md"),
			templatePath: "templates/rules/graphify.md",
			mode:         0644,
		},
		{
			relPath:      filepath.Join(".agents", "workflows", "graphify.md"),
			templatePath: "templates/workflows/graphify.md",
			mode:         0644,
		},
	}

	if !opts.Minimal {
		items = append(items, scaffoldItem{
			relPath:      filepath.Join("scripts", "graphify_sync.sh"),
			templatePath: "templates/scripts/graphify_sync.sh",
			mode:         0755,
		})
	}

	// Merge targets
	items = append(items, scaffoldItem{
		relPath:      filepath.Join(".agents", "AGENTS.md"),
		templatePath: "templates/agents/agents_snippet.md",
		mode:         0644,
		isMerge:      true,
		mergeMarker:  "rules/graphify.md",
	})

	if !opts.Minimal {
		items = append(items, scaffoldItem{
			relPath:      ".gitignore",
			templatePath: "templates/gitignore_snippet.txt",
			mode:         0644,
			isMerge:      true,
			mergeMarker:  "graphify-out/",
		})
	}

	for _, item := range items {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		absTarget := filepath.Join(cleanRoot, item.relPath)

		// Security boundary confinement check
		if _, err := auditor.ValidatePathBoundary(cleanRoot, absTarget); err != nil {
			return nil, fmt.Errorf("path boundary violation on %s: %w", item.relPath, err)
		}

		tmplBytes, err := s.fs.ReadFile(item.templatePath)
		if err != nil {
			return nil, fmt.Errorf("failed reading embedded template %s: %w", item.templatePath, err)
		}
		tmplContent := string(tmplBytes)

		_, statErr := os.Stat(absTarget)
		exists := statErr == nil

		var action ScaffoldAction

		if !item.isMerge {
			if !exists {
				action = ScaffoldAction{
					RelativePath: item.relPath,
					AbsolutePath: absTarget,
					ActionType:   ActionCreateFile,
					FileMode:     item.mode,
					Content:      tmplContent,
					Reason:       "New file scaffolded from official Antigravity 2.0 template",
					ExistedPrior: false,
				}
				plan.CreatedFiles++
			} else if opts.Force {
				action = ScaffoldAction{
					RelativePath: item.relPath,
					AbsolutePath: absTarget,
					ActionType:   ActionOverwriteForce,
					FileMode:     item.mode,
					Content:      tmplContent,
					Reason:       "Overwritten with latest official Antigravity 2.0 template (--force)",
					ExistedPrior: true,
				}
				plan.ModifiedFiles++
			} else {
				action = ScaffoldAction{
					RelativePath: item.relPath,
					AbsolutePath: absTarget,
					ActionType:   ActionSkipExists,
					FileMode:     item.mode,
					Reason:       "File already exists (skipped; use --force to overwrite)",
					ExistedPrior: true,
				}
				plan.SkippedFiles++
			}
		} else {
			// Merge item handling (.gitignore / AGENTS.md)
			if !exists {
				action = ScaffoldAction{
					RelativePath: item.relPath,
					AbsolutePath: absTarget,
					ActionType:   ActionCreateFile,
					FileMode:     item.mode,
					Content:      tmplContent,
					Reason:       "Created new manifest/ignore file with Graphify configurations",
					ExistedPrior: false,
				}
				plan.CreatedFiles++
			} else {
				existingData, readErr := os.ReadFile(absTarget)
				if readErr == nil && strings.Contains(string(existingData), item.mergeMarker) {
					action = ScaffoldAction{
						RelativePath: item.relPath,
						AbsolutePath: absTarget,
						ActionType:   ActionSkipExists,
						FileMode:     item.mode,
						Reason:       "Already configured with Graphify rules (marker found)",
						ExistedPrior: true,
					}
					plan.SkippedFiles++
				} else {
					action = ScaffoldAction{
						RelativePath: item.relPath,
						AbsolutePath: absTarget,
						ActionType:   ActionAppendMerge,
						FileMode:     item.mode,
						Content:      tmplContent,
						Reason:       "Appended Graphify declarations to existing file without overriding custom entries",
						ExistedPrior: true,
					}
					plan.ModifiedFiles++
				}
			}
		}

		plan.Actions = append(plan.Actions, action)
	}

	plan.TotalActions = len(plan.Actions)
	plan.ExecutionTimeMs = float64(time.Since(startTime).Microseconds()) / 1000.0
	return plan, nil
}

// Execute applies the planned actions atomically to disk.
func (s *DefaultScaffolderService) Execute(ctx context.Context, plan *ScaffoldPlan) error {
	if plan == nil || plan.DryRun || len(plan.Actions) == 0 {
		return nil
	}

	for i, action := range plan.Actions {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if action.ActionType == ActionSkipExists {
			continue
		}

		// Ensure parent directory exists with 0755 permissions
		parentDir := filepath.Dir(action.AbsolutePath)
		if mkdirErr := os.MkdirAll(parentDir, 0755); mkdirErr != nil {
			return fmt.Errorf("failed creating parent directory %s: %w", parentDir, mkdirErr)
		}

		var writeContent string

		switch action.ActionType {
		case ActionCreateFile, ActionOverwriteForce:
			writeContent = action.Content
		case ActionAppendMerge:
			existingData, err := os.ReadFile(action.AbsolutePath)
			if err != nil {
				writeContent = action.Content
			} else {
				trimmed := strings.TrimRight(string(existingData), "\n")
				writeContent = trimmed + "\n\n" + strings.TrimSpace(action.Content) + "\n"
			}
		}

		tmpPath := action.AbsolutePath + ".tmp"
		if writeErr := os.WriteFile(tmpPath, []byte(writeContent), action.FileMode); writeErr != nil {
			return fmt.Errorf("failed writing temporary staging buffer %s: %w", tmpPath, writeErr)
		}

		if renameErr := os.Rename(tmpPath, action.AbsolutePath); renameErr != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("failed committing atomic rename to %s: %w", action.AbsolutePath, renameErr)
		}

		plan.Actions[i].Applied = true
	}

	return nil
}

// Verify checks that all required Antigravity 2.0 files exist and conform to syntax rules.
func (s *DefaultScaffolderService) Verify(ctx context.Context, workspaceRoot string) (*ScaffoldVerificationReport, error) {
	cleanRoot := filepath.Clean(workspaceRoot)
	report := &ScaffoldVerificationReport{
		Timestamp:     time.Now(),
		WorkspaceRoot: cleanRoot,
		AllValid:      true,
		Errors:        make([]string, 0),
	}

	rulesPath := filepath.Join(cleanRoot, ".agents", "rules", "graphify.md")
	if _, err := os.Stat(rulesPath); err == nil {
		report.RulesFound = true
	} else {
		report.AllValid = false
		report.Errors = append(report.Errors, "missing .agents/rules/graphify.md")
	}

	workflowPath := filepath.Join(cleanRoot, ".agents", "workflows", "graphify.md")
	if _, err := os.Stat(workflowPath); err == nil {
		report.WorkflowFound = true
	} else {
		report.AllValid = false
		report.Errors = append(report.Errors, "missing .agents/workflows/graphify.md")
	}

	scriptPath := filepath.Join(cleanRoot, "scripts", "graphify_sync.sh")
	if info, err := os.Stat(scriptPath); err == nil {
		report.ScriptFound = true
		// Verify executable bit on POSIX
		if info.Mode()&0111 == 0 {
			report.Errors = append(report.Errors, "scripts/graphify_sync.sh is not executable (mode != 0755)")
		}
	} else {
		report.Errors = append(report.Errors, "missing scripts/graphify_sync.sh")
	}

	return report, nil
}
