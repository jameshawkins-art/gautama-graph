package runner

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/jameshawkins-art/gautama-graph/internal/auditor"
)

// DefaultOrchestrator coordinates the 4-stage Graphify pipeline.
type DefaultOrchestrator struct {
	binaryManager    BinaryManager
	subprocessRunner SubprocessRunner
	astAuditor       auditor.ASTGraphAuditorService
	docAuditor       auditor.DocGraphAuditorService
}

// NewDefaultOrchestrator initializes a new DefaultOrchestrator instance.
func NewDefaultOrchestrator(bm BinaryManager, sr SubprocessRunner, ast auditor.ASTGraphAuditorService, doc auditor.DocGraphAuditorService) *DefaultOrchestrator {
	return &DefaultOrchestrator{
		binaryManager:    bm,
		subprocessRunner: sr,
		astAuditor:       ast,
		docAuditor:       doc,
	}
}

// NewStandardOrchestrator constructs a DefaultOrchestrator with standard subsystem adapters.
func NewStandardOrchestrator(workspaceRoot string) *DefaultOrchestrator {
	cleanRoot := filepath.Clean(workspaceRoot)
	downloader := NewDefaultReleaseDownloader("")
	binManager := NewDefaultBinaryManager(downloader)
	subRunner := NewDefaultSubprocessRunner()

	astConfig := auditor.Config{
		WorkspaceRootPath: cleanRoot,
		AuditorTimeout:    60 * time.Second,
		MinConfidence:     0.8,
	}
	astEngine := auditor.NewDefaultEngine(astConfig)
	docAuditor := auditor.NewDocGraphAuditor(cleanRoot)

	return &DefaultOrchestrator{
		binaryManager:    binManager,
		subprocessRunner: subRunner,
		astAuditor:       astEngine,
		docAuditor:       docAuditor,
	}
}

// RunPipeline executes Download -> Base Extraction -> AST Code Relationship Audit -> Doc Graph Audit.
func (o *DefaultOrchestrator) RunPipeline(ctx context.Context, cfg RunnerConfig) (*PipelineReport, error) {
	startTime := time.Now()
	cleanWorkspace := filepath.Clean(cfg.WorkspaceRootPath)
	if cleanWorkspace == "" {
		cleanWorkspace = "."
	}

	if cfg.ExecutionTimeout <= 0 {
		cfg.ExecutionTimeout = 120 * time.Second
	}
	ctx, cancel := context.WithTimeout(ctx, cfg.ExecutionTimeout)
	defer cancel()

	report := &PipelineReport{
		Timestamp: startTime,
		Stages:    make([]PipelineStageStatus, 0, 4),
	}

	// --------------------------------------------------
	// Stage 1: Binary Lifecycle & Cache Resolution
	// --------------------------------------------------
	s1Start := time.Now()
	binPath, version, binErr := o.binaryManager.EnsureBinary(ctx, cfg)
	s1Duration := time.Since(s1Start)

	if binErr != nil {
		report.Stages = append(report.Stages, PipelineStageStatus{
			StageName: "Binary Resolution & Download",
			Duration:  s1Duration,
			Success:   false,
			Error:     binErr.Error(),
		})
		report.TotalDuration = time.Since(startTime)
		return report, fmt.Errorf("stage 1 (binary resolution) failed: %w", binErr)
	}

	report.BinarySource = "RESOLVED"
	report.BinaryVersion = version
	report.Stages = append(report.Stages, PipelineStageStatus{
		StageName: "Binary Resolution & Download",
		Duration:  s1Duration,
		Success:   true,
	})

	// --------------------------------------------------
	// Stage 2: Base Graphify Extraction & Clustering
	// --------------------------------------------------
	s2Start := time.Now()
	_, _, execErr := o.subprocessRunner.ExecuteCommand(ctx, binPath, cleanWorkspace, "update", ".")
	s2Duration := time.Since(s2Start)

	if execErr != nil {
		report.Stages = append(report.Stages, PipelineStageStatus{
			StageName: "Base Graphify Extraction",
			Duration:  s2Duration,
			Success:   false,
			Error:     execErr.Error(),
		})
		report.TotalDuration = time.Since(startTime)
		return report, fmt.Errorf("stage 2 (base extraction) failed: %w", execErr)
	}

	report.Stages = append(report.Stages, PipelineStageStatus{
		StageName: "Base Graphify Extraction",
		Duration:  s2Duration,
		Success:   true,
	})

	// --------------------------------------------------
	// Stage 3: Deterministic AST Code Relationship Audit
	// --------------------------------------------------
	s3Start := time.Now()
	graphPath := filepath.Join(cleanWorkspace, "graphify-out", "graph.json")
	astReport, astErr := o.astAuditor.AuditGraphFile(ctx, graphPath, cfg.VerboseLogging)
	s3Duration := time.Since(s3Start)

	if astErr != nil {
		report.Stages = append(report.Stages, PipelineStageStatus{
			StageName: "AST Code Relationship Audit",
			Duration:  s3Duration,
			Success:   false,
			Error:     astErr.Error(),
		})
		report.TotalDuration = time.Since(startTime)
		return report, fmt.Errorf("stage 3 (AST audit) failed: %w", astErr)
	}

	report.GraphEdgeCount = astReport.TotalEdges
	report.PrunedPhantoms = astReport.PrunedPhantomCount
	report.Stages = append(report.Stages, PipelineStageStatus{
		StageName: "AST Code Relationship Audit",
		Duration:  s3Duration,
		Success:   true,
	})

	// --------------------------------------------------
	// Stage 4: Markdown Documentation Graph Audit
	// --------------------------------------------------
	s4Start := time.Now()
	docReport, docErr := o.docAuditor.AuditDocGraph(ctx)
	s4Duration := time.Since(s4Start)

	if docErr != nil {
		report.Stages = append(report.Stages, PipelineStageStatus{
			StageName: "Markdown Doc Graph Audit",
			Duration:  s4Duration,
			Success:   false,
			Error:     docErr.Error(),
		})
		report.TotalDuration = time.Since(startTime)
		return report, fmt.Errorf("stage 4 (doc audit) failed: %w", docErr)
	}

	report.DocOrphanCount = docReport.OrphanCount
	report.BrokenDocLinks = docReport.BrokenLinkCount
	report.GraphNodeCount = docReport.TotalDocNodes
	report.Stages = append(report.Stages, PipelineStageStatus{
		StageName: "Markdown Doc Graph Audit",
		Duration:  s4Duration,
		Success:   true,
	})

	report.TotalDuration = time.Since(startTime)
	return report, nil
}
