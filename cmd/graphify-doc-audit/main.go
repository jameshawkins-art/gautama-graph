package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gautama-graph/internal/auditor"
)

func main() {
	strictFlag := flag.Bool("strict", false, "Fail with non-zero exit code if orphans or broken links exist")
	workspaceFlag := flag.String("workspace", "", "Path to workspace root (defaults to CWD)")
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
	report, err := docAuditor.AuditDocGraph(context.Background())
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

	fmt.Println("--------------------------------------------------")

	if *strictFlag && (report.OrphanCount > 0 || report.BrokenLinkCount > 0) {
		fmt.Println("🔴 STRICT AUDIT FAILED: Fix orphans or broken relative links before committing.")
		os.Exit(1)
	}

	fmt.Println("✅ Documentation Graph Audit Completed Successfully.")
}
