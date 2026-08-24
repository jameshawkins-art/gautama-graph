package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jameshawkins-art/gautama-graph/internal/runner"
)

func main() {
	strictFlag := flag.Bool("strict", false, "Fail with non-zero exit code if phantom edges or doc issues exist")
	workspaceFlag := flag.String("workspace", "", "Path to workspace root (defaults to CWD)")
	forceDownloadFlag := flag.Bool("force-download", false, "Force download of latest release binary from GitHub")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose candidate logging")
	timeoutFlag := flag.Duration("timeout", 180*time.Second, "Overall pipeline execution timeout")
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

	cfg := runner.RunnerConfig{
		WorkspaceRootPath: workspaceRoot,
		ExecutionTimeout:  *timeoutFlag,
		StrictAudit:       *strictFlag,
		ForceDownload:     *forceDownloadFlag,
		VerboseLogging:    *verboseFlag,
	}

	orchestrator := runner.NewStandardOrchestrator(workspaceRoot)
	report, err := orchestrator.RunPipeline(context.Background(), cfg)
	if err != nil {
		log.Fatalf("❌ Gautama Graph Pipeline Error: %v", err)
	}

	fmt.Println("==================================================")
	fmt.Println("🚀 Gautama Graph Master Pipeline Summary")
	fmt.Println("==================================================")
	fmt.Printf("  Workspace Root  : %s\n", workspaceRoot)
	fmt.Printf("  Binary Source   : %s (%s)\n", report.BinarySource, report.BinaryVersion)
	fmt.Printf("  Total Duration  : %v\n", report.TotalDuration)
	fmt.Println("--------------------------------------------------")
	fmt.Println("📋 Stage Breakdown:")
	for i, stage := range report.Stages {
		status := "✅"
		if !stage.Success {
			status = "❌"
		}
		fmt.Printf("  [%d/4] %s %-30s (%v)\n", i+1, status, stage.StageName, stage.Duration)
		if stage.Error != "" {
			fmt.Printf("        Error: %s\n", stage.Error)
		}
	}
	fmt.Println("--------------------------------------------------")
	fmt.Printf("📊 Graph Metrics:\n")
	fmt.Printf("  Total Document Nodes : %d\n", report.GraphNodeCount)
	fmt.Printf("  Total Evaluated Edges: %d\n", report.GraphEdgeCount)
	fmt.Printf("  Pruned Phantom Edges : %d\n", report.PrunedPhantoms)
	fmt.Printf("  Doc Orphan Count     : %d\n", report.DocOrphanCount)
	fmt.Printf("  Broken Doc Links     : %d\n", report.BrokenDocLinks)
	fmt.Println("--------------------------------------------------")
	fmt.Println("📁 Generated Artifacts:")
	fmt.Printf("  - %s\n", filepath.Join(workspaceRoot, "graphify-out", "graph.json"))
	fmt.Printf("  - %s\n", filepath.Join(workspaceRoot, "graphify-out", "GRAPH_REPORT.md"))
	fmt.Printf("  - %s\n", filepath.Join(workspaceRoot, "graphify-out", "graph.html"))
	fmt.Printf("  - %s\n", filepath.Join(workspaceRoot, "graphify-out", "doc_graph_audit.json"))
	fmt.Println("==================================================")

	if *strictFlag && (report.PrunedPhantoms > 0 || report.BrokenDocLinks > 0) {
		fmt.Printf("🔴 STRICT PIPELINE FAILED: %d phantoms or %d broken links detected.\n", report.PrunedPhantoms, report.BrokenDocLinks)
		os.Exit(1)
	}

	fmt.Println("✅ Turnkey Knowledge Graph Pipeline Succeeded.")
}
