package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gautama-graph/internal/auditor"
)

func main() {
	strictFlag := flag.Bool("strict", false, "Fail with non-zero exit code if phantom edges are detected")
	workspaceFlag := flag.String("workspace", "", "Path to workspace root (defaults to CWD)")
	graphFlag := flag.String("graph", "", "Path to graph.json (defaults to graphify-out/graph.json)")
	verboseFlag := flag.Bool("verbose", false, "Enable verbose candidate logging")
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

	graphPath := *graphFlag
	if graphPath == "" {
		graphPath = filepath.Join(workspaceRoot, "graphify-out", "graph.json")
	}

	if _, err := os.Stat(graphPath); os.IsNotExist(err) {
		log.Fatalf("Graph file not found at %s. Please run 'graphify update .' first.", graphPath)
	}

	cfg := auditor.Config{
		WorkspaceRootPath: workspaceRoot,
		AuditorTimeout:    60 * time.Second,
		MinConfidence:     0.8,
	}

	engine := auditor.NewDefaultEngine(cfg)
	report, err := engine.AuditGraphFile(context.Background(), graphPath, *verboseFlag)
	if err != nil {
		log.Fatalf("AST code relationship audit error: %v", err)
	}

	fmt.Println("--------------------------------------------------")
	fmt.Println("⚡ Gautama Social AST Code Relationship Audit Summary")
	fmt.Println("--------------------------------------------------")
	fmt.Printf("  Workspace Root  : %s\n", workspaceRoot)
	fmt.Printf("  Graph File      : %s\n", graphPath)
	fmt.Printf("  Duration        : %v\n", report.Duration)
	fmt.Printf("  Total Edges     : %d\n", report.TotalEdges)
	fmt.Printf("  Verified AST    : %d\n", report.VerifiedASTCount)
	fmt.Printf("  Pruned Phantoms : %d\n", report.PrunedPhantomCount)
	fmt.Printf("  Heuristic/Other : %d\n", report.HeuristicCount)
	fmt.Println("--------------------------------------------------")

	if *verboseFlag && len(report.AuditedEdges) > 0 {
		fmt.Println("\n🔍 Candidate Edge Provenance Details:")
		for _, edge := range report.AuditedEdges {
			fmt.Printf("  - [%s] %s (confidence: %.1f)\n", edge.ProvenanceStatus, edge.ID, edge.Confidence)
		}
		fmt.Println("--------------------------------------------------")
	}

	if *strictFlag && report.PrunedPhantomCount > 0 {
		fmt.Printf("🔴 STRICT AUDIT FAILED: %d phantom edges pruned from %s.\n", report.PrunedPhantomCount, graphPath)
		os.Exit(1)
	}

	fmt.Println("✅ AST Code Relationship Audit Completed Successfully.")
}
