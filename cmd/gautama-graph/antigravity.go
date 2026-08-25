package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jameshawkins-art/gautama-graph/internal/scaffold"
)

// handleAntigravityCommand parses and executes the `antigravity` subcommand.
func handleAntigravityCommand(args []string) error {
	fs := flag.NewFlagSet("antigravity", flag.ExitOnError)
	workspaceFlag := fs.String("workspace", "", "Path to workspace root (defaults to CWD)")
	_ = fs.Bool("project", false, "Scaffold full project Antigravity environment")
	dryRunFlag := fs.Bool("dry-run", false, "Preview planned scaffolding actions without modifying disk")
	forceFlag := fs.Bool("force", false, "Force overwrite of existing template files")
	minimalFlag := fs.Bool("minimal", false, "Scaffold minimal rules and workflows only (no sync scripts or gitignore)")
	verboseFlag := fs.Bool("verbose", false, "Enable verbose logging")

	// If the user supplied "install" as the first argument after "antigravity", slice it off
	parseArgs := args
	if len(parseArgs) > 0 && parseArgs[0] == "install" {
		parseArgs = parseArgs[1:]
	}

	if err := fs.Parse(parseArgs); err != nil {
		return err
	}

	workspaceRoot := *workspaceFlag
	if workspaceRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("failed getting working directory: %w", err)
		}
		workspaceRoot = cwd
	}
	workspaceRoot = filepath.Clean(workspaceRoot)

	opts := scaffold.ScaffoldOptions{
		WorkspaceRoot: workspaceRoot,
		Force:         *forceFlag,
		DryRun:        *dryRunFlag,
		Minimal:       *minimalFlag,
		WithScripts:   !*minimalFlag,
		WithGitIgnore: !*minimalFlag,
		WithMakefile:  !*minimalFlag,
		Verbose:       *verboseFlag,
	}

	svc := scaffold.NewDefaultScaffolderService()
	ctx := context.Background()

	plan, err := svc.Plan(ctx, opts)
	if err != nil {
		return fmt.Errorf("failed formulating scaffold plan: %w", err)
	}

	fmt.Println("==================================================")
	fmt.Println("🚀 Gautama Graph Antigravity Environment Setup")
	fmt.Println("==================================================")
	fmt.Printf("  Workspace Root : %s\n", plan.WorkspaceRoot)
	fmt.Printf("  Dry Run Mode   : %t\n", plan.DryRun)
	fmt.Printf("  Force Mode     : %t\n", opts.Force)
	fmt.Println("--------------------------------------------------")
	fmt.Println("📋 Planned Scaffolding Actions:")
	for _, action := range plan.Actions {
		icon := "[+]"
		switch action.ActionType {
		case scaffold.ActionSkipExists:
			icon = "[~]"
		case scaffold.ActionOverwriteForce:
			icon = "[!]"
		}
		fmt.Printf("  %s %-16s %s (%s)\n", icon, action.ActionType, action.RelativePath, action.Reason)
	}
	fmt.Println("--------------------------------------------------")

	if !plan.DryRun {
		if err := svc.Execute(ctx, plan); err != nil {
			return fmt.Errorf("scaffolding execution failed: %w", err)
		}
		fmt.Println("✅ Files successfully written to workspace.")
	} else {
		fmt.Println("🔎 Dry-run completed. Zero disk changes made.")
	}

	fmt.Printf("📊 Scaffolding Summary:\n")
	fmt.Printf("  Total Actions  : %d\n", plan.TotalActions)
	fmt.Printf("  Created Files  : %d\n", plan.CreatedFiles)
	fmt.Printf("  Modified Files : %d\n", plan.ModifiedFiles)
	fmt.Printf("  Skipped Files  : %d\n", plan.SkippedFiles)
	fmt.Printf("  Duration       : %.2fms\n", plan.ExecutionTimeMs)
	fmt.Println("==================================================")

	if !plan.DryRun {
		report, verifyErr := svc.Verify(ctx, workspaceRoot)
		if verifyErr == nil && report.AllValid {
			fmt.Println("🎯 Antigravity Environment Verification: PASS")
		} else if verifyErr != nil || !report.AllValid {
			fmt.Println("⚠️  Antigravity Environment Verification: WARNING (Some items missing or incomplete)")
			for _, itemErr := range report.Errors {
				fmt.Printf("    - %s\n", itemErr)
			}
		}
	}

	return nil
}
