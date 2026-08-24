package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jameshawkins-art/gautama-graph/internal/auditor"
)

func main() {
	strictFlag := flag.Bool("strict", false, "Fail with non-zero exit code if orphans or broken links exist")
	workspaceFlag := flag.String("workspace", "", "Path to workspace root (defaults to CWD)")
	fixFlag := flag.Bool("fix", false, "Automatically remediate broken relative links in-place across markdown files")
	dryRunFlag := flag.Bool("dry-run", false, "Preview planned link remediation actions without modifying disk")
	detectCyclesFlag := flag.Bool("detect-cycles", false, "Run Tarjan's SCC cycle detector on the documentation graph")
	checkAnchorsFlag := flag.Bool("check-anchors", false, "Verify GitHub-Flavored Markdown heading anchors on relative links")
	flag.Parse()

	workspaceRoot := *workspaceFlag
	if workspaceRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("Failed to get working directory: %v", err)
		}
		workspaceRoot = cwd
	}

	workspaceRoot = filepath.Clean(workspaceRoot)
	docAuditor := auditor.NewDocGraphAuditor(workspaceRoot)
	ctx := context.Background()

	report, err := docAuditor.AuditDocGraph(ctx)
	if err != nil {
		log.Fatalf("Doc graph audit error: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("📊 Gautama Social Markdown Doc Graph Audit Summary")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("  Workspace Root  : %s\n", workspaceRoot)
	fmt.Printf("  Total Doc Nodes : %d\n", report.TotalDocNodes)
	fmt.Printf("  Total Doc Edges : %d\n", report.TotalDocEdges)
	fmt.Printf("  Orphan Nodes    : %d\n", report.OrphanCount)
	fmt.Printf("  Broken Links    : %d\n", report.BrokenLinkCount)
	fmt.Printf("  Diagnostic File : graphify-out/doc_graph_audit.json\n")

	remediator := auditor.NewDefaultDocRemediatorService()

	// 1. Cycle Detection
	if *detectCyclesFlag {
		cycleReport, cycleErr := remediator.DetectCycles(ctx, workspaceRoot)
		if cycleErr == nil && cycleReport != nil {
			fmt.Printf("  Circular Cycles : %d\n", cycleReport.TotalCycles)
			if cycleReport.TotalCycles > 0 {
				fmt.Println("\n🔄 Circular Documentation Cycles Detected:")
				for _, c := range cycleReport.Cycles {
					fmt.Printf("  - Cycle [%s]: %v\n", c.CycleID, c.DocChain)
				}
			}
		}
	}

	// 2. Heading Anchor Indexing
	if *checkAnchorsFlag {
		anchorTables, anchorErr := remediator.IndexHeadingAnchors(ctx, workspaceRoot)
		if anchorErr == nil {
			totalAnchors := 0
			for _, t := range anchorTables {
				totalAnchors += len(t.Anchors)
			}
			fmt.Printf("  Heading Anchors : %d indexed across %d documents\n", totalAnchors, len(anchorTables))
		}
	}

	if report.OrphanCount > 0 {
		fmt.Println("\n⚠️  Orphan Markdown Documents (Degree == 0):")
		for _, orphan := range report.OrphanNodes {
			fmt.Printf("  - %s\n", orphan)
		}
	}

	if report.BrokenLinkCount > 0 {
		fmt.Println("\n❌ Broken Relative Markdown Links:")
		for _, broken := range report.BrokenLinks {
			fmt.Printf("  - Source: %s -> Target: %s (%s)\n", broken.SourceFile, broken.LinkTarget, broken.ErrorReason)
		}
	}

	// 3. Auto-Remediation (--fix / --dry-run)
	if *fixFlag || *dryRunFlag {
		plan, planErr := remediator.PlanRemediation(ctx, workspaceRoot, *dryRunFlag)
		if planErr != nil {
			log.Fatalf("Failed planning link remediation: %v", planErr)
		}

		fmt.Println("\n🛠️  Documentation Link Auto-Remediation Plan:")
		fmt.Printf("  Planned Actions : %d modifications across %d files (execution time: %.2fms)\n", plan.TotalActions, plan.ModifiedDocs, plan.ExecutionTimeMs)

		for _, act := range plan.Actions {
			fmt.Printf("  - [%s] %s:%d -> %s (was: %s)\n", act.Rule, act.SourceFile, act.LineNumber, act.CanonicalRelPath, act.OriginalTarget)
		}

		if *fixFlag {
			if applyErr := remediator.ApplyRemediation(ctx, plan); applyErr != nil {
				log.Fatalf("Failed applying remediation plan: %v", applyErr)
			}
			fmt.Println("✅ Successfully applied all remediation actions in-place.")
		}
	}

	fmt.Println("--------------------------------------------------")

	if *strictFlag && (report.OrphanCount > 0 || report.BrokenLinkCount > 0) {
		fmt.Println("🔴 STRICT AUDIT FAILED: Fix orphans or broken relative links before committing.")
		os.Exit(1)
	}

	fmt.Println("✅ Documentation Graph Audit Completed Successfully.")
}
