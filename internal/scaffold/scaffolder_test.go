package scaffold_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jameshawkins-art/gautama-graph/internal/scaffold"
)

func TestScaffolder_NominalInstallation(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	opts := scaffold.ScaffoldOptions{
		WorkspaceRoot: tempDir,
		Force:         false,
		DryRun:        false,
		Minimal:       false,
		WithScripts:   true,
		WithGitIgnore: true,
	}

	plan, err := svc.Plan(ctx, opts)
	if err != nil {
		t.Fatalf("Plan returned unexpected error: %v", err)
	}

	if plan.TotalActions != 5 {
		t.Errorf("expected 5 planned actions, got %d", plan.TotalActions)
	}
	if plan.CreatedFiles != 5 {
		t.Errorf("expected 5 created files, got %d", plan.CreatedFiles)
	}

	if err := svc.Execute(ctx, plan); err != nil {
		t.Fatalf("Execute returned unexpected error: %v", err)
	}

	// Verify files on disk
	rulesPath := filepath.Join(tempDir, ".agents", "rules", "graphify.md")
	if _, err := os.Stat(rulesPath); err != nil {
		t.Errorf("expected %s to exist on disk", rulesPath)
	}

	workflowPath := filepath.Join(tempDir, ".agents", "workflows", "graphify.md")
	if _, err := os.Stat(workflowPath); err != nil {
		t.Errorf("expected %s to exist on disk", workflowPath)
	}

	scriptPath := filepath.Join(tempDir, "scripts", "graphify_sync.sh")
	stat, err := os.Stat(scriptPath)
	if err != nil {
		t.Errorf("expected %s to exist on disk", scriptPath)
	} else if stat.Mode()&0111 == 0 {
		t.Errorf("expected %s to be executable (0755), got mode %v", scriptPath, stat.Mode())
	}

	report, err := svc.Verify(ctx, tempDir)
	if err != nil {
		t.Fatalf("Verify failed: %v", err)
	}
	if !report.AllValid {
		t.Errorf("expected verification report to be AllValid, got errors: %v", report.Errors)
	}
}

func TestScaffolder_IdempotentSkip(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	opts := scaffold.ScaffoldOptions{
		WorkspaceRoot: tempDir,
		Force:         false,
	}

	// First run
	plan1, err := svc.Plan(ctx, opts)
	if err != nil {
		t.Fatalf("first Plan failed: %v", err)
	}
	if err := svc.Execute(ctx, plan1); err != nil {
		t.Fatalf("first Execute failed: %v", err)
	}

	// Second run without --force
	plan2, err := svc.Plan(ctx, opts)
	if err != nil {
		t.Fatalf("second Plan failed: %v", err)
	}

	if plan2.CreatedFiles != 0 {
		t.Errorf("expected 0 created files on second run, got %d", plan2.CreatedFiles)
	}
	if plan2.ModifiedFiles != 0 {
		t.Errorf("expected 0 modified files on second run, got %d", plan2.ModifiedFiles)
	}
	if plan2.SkippedFiles != 5 {
		t.Errorf("expected 5 skipped files on second run, got %d", plan2.SkippedFiles)
	}

	// Ensure execution does nothing
	if err := svc.Execute(ctx, plan2); err != nil {
		t.Fatalf("second Execute failed: %v", err)
	}
}

func TestScaffolder_ForceOverwrite(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	// First run
	plan1, _ := svc.Plan(ctx, scaffold.ScaffoldOptions{WorkspaceRoot: tempDir})
	_ = svc.Execute(ctx, plan1)

	// Mutate a rule file
	rulesPath := filepath.Join(tempDir, ".agents", "rules", "graphify.md")
	_ = os.WriteFile(rulesPath, []byte("CUSTOM_OUTDATED_RULE"), 0644)

	// Second run with --force
	plan2, err := svc.Plan(ctx, scaffold.ScaffoldOptions{WorkspaceRoot: tempDir, Force: true})
	if err != nil {
		t.Fatalf("Plan with force failed: %v", err)
	}

	foundOverwrite := false
	for _, action := range plan2.Actions {
		if action.RelativePath == filepath.Join(".agents", "rules", "graphify.md") {
			if action.ActionType == scaffold.ActionOverwriteForce {
				foundOverwrite = true
			}
		}
	}
	if !foundOverwrite {
		t.Errorf("expected OVERWRITE_FORCE on .agents/rules/graphify.md with --force")
	}

	if err := svc.Execute(ctx, plan2); err != nil {
		t.Fatalf("Execute with force failed: %v", err)
	}

	content, _ := os.ReadFile(rulesPath)
	if strings.Contains(string(content), "CUSTOM_OUTDATED_RULE") {
		t.Errorf("expected rules file to be overwritten with template, but found outdated custom string")
	}
}

func TestScaffolder_DryRun(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	opts := scaffold.ScaffoldOptions{
		WorkspaceRoot: tempDir,
		DryRun:        true,
	}

	plan, err := svc.Plan(ctx, opts)
	if err != nil {
		t.Fatalf("Plan returned error: %v", err)
	}

	if !plan.DryRun {
		t.Errorf("expected plan.DryRun to be true")
	}

	if err := svc.Execute(ctx, plan); err != nil {
		t.Fatalf("Execute returned error on dry-run: %v", err)
	}

	// Verify nothing was written
	rulesPath := filepath.Join(tempDir, ".agents", "rules", "graphify.md")
	if _, err := os.Stat(rulesPath); err == nil {
		t.Errorf("expected %s NOT to exist on disk during dry-run", rulesPath)
	}
}

func TestScaffolder_MergeFiles(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	// Pre-create .gitignore and AGENTS.md with user content
	gitIgnorePath := filepath.Join(tempDir, ".gitignore")
	_ = os.WriteFile(gitIgnorePath, []byte("node_modules/\n.env\n"), 0644)

	agentsDir := filepath.Join(tempDir, ".agents")
	_ = os.MkdirAll(agentsDir, 0755)
	agentsPath := filepath.Join(agentsDir, "AGENTS.md")
	_ = os.WriteFile(agentsPath, []byte("# My Custom Agents\n- [architect.md](./personas/architect.md)\n"), 0644)

	opts := scaffold.ScaffoldOptions{
		WorkspaceRoot: tempDir,
	}

	plan, err := svc.Plan(ctx, opts)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if err := svc.Execute(ctx, plan); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Verify merged .gitignore
	gitIgnoreContent, _ := os.ReadFile(gitIgnorePath)
	if !strings.Contains(string(gitIgnoreContent), "node_modules/") {
		t.Errorf("expected original gitignore content to be preserved")
	}
	if !strings.Contains(string(gitIgnoreContent), "graphify-out/") {
		t.Errorf("expected graphify-out/ to be appended to gitignore")
	}

	// Verify merged AGENTS.md
	agentsContent, _ := os.ReadFile(agentsPath)
	if !strings.Contains(string(agentsContent), "architect.md") {
		t.Errorf("expected original AGENTS.md content to be preserved")
	}
	if !strings.Contains(string(agentsContent), "rules/graphify.md") {
		t.Errorf("expected rules/graphify.md to be appended to AGENTS.md")
	}
}

func TestScaffolder_MinimalMode(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	opts := scaffold.ScaffoldOptions{
		WorkspaceRoot: tempDir,
		Minimal:       true,
	}

	plan, err := svc.Plan(ctx, opts)
	if err != nil {
		t.Fatalf("Plan failed: %v", err)
	}

	if plan.TotalActions != 3 {
		t.Errorf("expected 3 actions in minimal mode, got %d", plan.TotalActions)
	}

	if err := svc.Execute(ctx, plan); err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	// Scripts should NOT exist
	scriptPath := filepath.Join(tempDir, "scripts", "graphify_sync.sh")
	if _, err := os.Stat(scriptPath); err == nil {
		t.Errorf("expected scripts/graphify_sync.sh NOT to exist in minimal mode")
	}
}

func TestScaffolder_EmptyWorkspaceRoot_And_EmptyExecute(t *testing.T) {
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	// Empty execute
	if err := svc.Execute(ctx, nil); err != nil {
		t.Errorf("Execute(nil) returned error: %v", err)
	}
	emptyPlan := &scaffold.ScaffoldPlan{Actions: []scaffold.ScaffoldAction{}}
	if err := svc.Execute(ctx, emptyPlan); err != nil {
		t.Errorf("Execute(empty) returned error: %v", err)
	}

	// Empty workspace root should default to CWD without error on Plan
	plan, err := svc.Plan(ctx, scaffold.ScaffoldOptions{WorkspaceRoot: "", DryRun: true})
	if err != nil {
		t.Fatalf("Plan with empty workspace failed: %v", err)
	}
	if plan.WorkspaceRoot == "" {
		t.Errorf("expected non-empty WorkspaceRoot")
	}
}

func TestScaffolder_ContextCancellation(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := svc.Plan(ctx, scaffold.ScaffoldOptions{WorkspaceRoot: tempDir})
	if err == nil {
		t.Errorf("expected context cancellation error on Plan")
	}

	validPlan, _ := svc.Plan(context.Background(), scaffold.ScaffoldOptions{WorkspaceRoot: tempDir})
	err = svc.Execute(ctx, validPlan)
	if err == nil {
		t.Errorf("expected context cancellation error on Execute")
	}
}

func TestScaffolder_Verify_Failures(t *testing.T) {
	tempDir := t.TempDir()
	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	// Initially empty directory
	report, err := svc.Verify(ctx, tempDir)
	if err != nil {
		t.Fatalf("Verify returned error: %v", err)
	}
	if report.AllValid {
		t.Errorf("expected AllValid to be false on empty workspace")
	}
	if len(report.Errors) != 3 {
		t.Errorf("expected 3 errors on empty workspace, got %d", len(report.Errors))
	}

	// Create non-executable script
	scriptsDir := filepath.Join(tempDir, "scripts")
	_ = os.MkdirAll(scriptsDir, 0755)
	scriptPath := filepath.Join(scriptsDir, "graphify_sync.sh")
	_ = os.WriteFile(scriptPath, []byte("#!/bin/sh\n"), 0644) // not 0755

	report2, _ := svc.Verify(ctx, tempDir)
	foundScriptErr := false
	for _, e := range report2.Errors {
		if strings.Contains(e, "not executable") {
			foundScriptErr = true
		}
	}
	if !foundScriptErr {
		t.Errorf("expected non-executable script error in report, got: %v", report2.Errors)
	}
}
