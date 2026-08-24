package auditor

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// DocNodeResult represents a single markdown document node audit status.
type DocNodeResult struct {
	FilePath      string   `json:"file_path"`
	InDegree      int      `json:"in_degree"`
	OutDegree     int      `json:"out_degree"`
	IsOrphan      bool     `json:"is_orphan"`
	OutboundLinks []string `json:"outbound_links"`
}

// BrokenLinkResult represents a dead markdown link reference.
type BrokenLinkResult struct {
	SourceFile   string `json:"source_file"`
	LinkTarget   string `json:"link_target"`
	ResolvedPath string `json:"resolved_path"`
	ErrorReason  string `json:"error_reason"`
}

// DocAuditReport represents the complete diagnostic summary written to graphify-out/doc_graph_audit.json.
type DocAuditReport struct {
	Timestamp       time.Time          `json:"timestamp"`
	TotalDocNodes   int                `json:"total_doc_nodes"`
	TotalDocEdges   int                `json:"total_doc_edges"`
	OrphanCount     int                `json:"orphan_count"`
	BrokenLinkCount int                `json:"broken_link_count"`
	OrphanNodes     []string           `json:"orphan_nodes"`
	BrokenLinks     []BrokenLinkResult `json:"broken_links"`
	NodeDetails     []DocNodeResult    `json:"node_details"`
}

// DocGraphParser defines workspace markdown file scanning and link resolution.
type DocGraphParser interface {
	ParseWorkspaceDocs(ctx context.Context, workspaceRoot string) ([]DocNodeResult, []BrokenLinkResult, error)
}

// DocGraphStore handles atomic persistence of diagnostic audit reports.
type DocGraphStore interface {
	SaveDocAuditReport(ctx context.Context, outputPath string, report *DocAuditReport) error
}

// DocGraphAuditorService orchestrates doc graph parsing, connectivity analysis, and diagnostic output.
type DocGraphAuditorService interface {
	AuditDocGraph(ctx context.Context) (*DocAuditReport, error)
}

// DefaultDocGraphParser implements DocGraphParser using Go stdlib path/filepath and regexp.
type DefaultDocGraphParser struct {
	workspaceRoot string
	linkRegex     *regexp.Regexp
}

// NewDefaultDocGraphParser initializes a DefaultDocGraphParser instance.
func NewDefaultDocGraphParser(workspaceRoot string) *DefaultDocGraphParser {
	return &DefaultDocGraphParser{
		workspaceRoot: filepath.Clean(workspaceRoot),
		linkRegex:     regexp.MustCompile(`\[.*?\]\(((?:file:///|[./])?[^)]+)\)`),
	}
}

// ValidatePathBoundary checks if targetPath stays strictly inside workspaceRoot to prevent path traversal.
func ValidatePathBoundary(workspaceRoot, targetPath string) (string, error) {
	cleanRoot := filepath.Clean(workspaceRoot)
	cleanTarget := filepath.Clean(targetPath)

	rel, err := filepath.Rel(cleanRoot, cleanTarget)
	if err != nil || strings.HasPrefix(rel, "..") || rel == ".." {
		return "", fmt.Errorf("SECURITY_PATH_TRAVERSAL: target path %s escapes workspace root %s", cleanTarget, cleanRoot)
	}

	return cleanTarget, nil
}

// ParseWorkspaceDocs scans all .md files under workspaceRoot and extracts link references.
func (p *DefaultDocGraphParser) ParseWorkspaceDocs(ctx context.Context, workspaceRoot string) ([]DocNodeResult, []BrokenLinkResult, error) {
	cleanRoot := filepath.Clean(workspaceRoot)
	var brokenLinks []BrokenLinkResult

	inDegreeMap := make(map[string]int)
	outboundMap := make(map[string][]string)
	allDocPaths := make(map[string]bool)

	err := filepath.Walk(cleanRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip hidden dirs, node_modules, vendor, .git, graphify-out, node- binaries
		if info.IsDir() {
			base := info.Name()
			if (strings.HasPrefix(base, ".") && base != ".") || base == "node_modules" || base == "vendor" || base == "graphify-out" || strings.HasPrefix(base, "node-") {
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) != ".md" {
			return nil
		}

		relPath, _ := filepath.Rel(cleanRoot, path)
		allDocPaths[relPath] = true

		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil // Skip unreadable file gracefully
		}

		// Strip fenced code blocks and inline code snippets to avoid false-positive link parsing
		codeBlockRegex := regexp.MustCompile("(?s)```.*?```|`[^`\\n]+`")
		cleanedContent := codeBlockRegex.ReplaceAllString(string(content), "")

		matches := p.linkRegex.FindAllStringSubmatch(cleanedContent, -1)
		var outbound []string

		for _, match := range matches {
			if len(match) < 2 {
				continue
			}
			rawTarget := match[1]

			// Strip file:// prefix if present
			cleanedTarget := strings.TrimPrefix(rawTarget, "file://")

			// Ignore HTTP/HTTPS external links and anchor fragments
			if strings.HasPrefix(cleanedTarget, "http://") || strings.HasPrefix(cleanedTarget, "https://") || strings.HasPrefix(cleanedTarget, "#") {
				continue
			}

			// Strip anchor query fragments e.g. #L10-L20 or :L10-L20
			if hashIdx := strings.Index(cleanedTarget, "#"); hashIdx != -1 {
				cleanedTarget = cleanedTarget[:hashIdx]
			}
			if colonIdx := strings.Index(cleanedTarget, ":L"); colonIdx != -1 {
				cleanedTarget = cleanedTarget[:colonIdx]
			}
			if cleanedTarget == "" {
				continue
			}

			// Resolve relative path against source file directory
			sourceDir := filepath.Dir(path)
			var targetFullPath string
			if filepath.IsAbs(cleanedTarget) {
				targetFullPath = cleanedTarget
			} else {
				targetFullPath = filepath.Join(sourceDir, cleanedTarget)
			}

			// Security boundary check
			validatedPath, boundaryErr := ValidatePathBoundary(cleanRoot, targetFullPath)
			if boundaryErr != nil {
				brokenLinks = append(brokenLinks, BrokenLinkResult{
					SourceFile:   relPath,
					LinkTarget:   rawTarget,
					ResolvedPath: targetFullPath,
					ErrorReason:  boundaryErr.Error(),
				})
				continue
			}

			// Check target existence on disk
			targetInfo, statErr := os.Stat(validatedPath)
			if statErr != nil {
				brokenLinks = append(brokenLinks, BrokenLinkResult{
					SourceFile:   relPath,
					LinkTarget:   rawTarget,
					ResolvedPath: validatedPath,
					ErrorReason:  "File not found on disk",
				})
			} else if !targetInfo.IsDir() {
				targetRel, _ := filepath.Rel(cleanRoot, validatedPath)
				inDegreeMap[targetRel]++
				outbound = append(outbound, targetRel)
			}
		}

		outboundMap[relPath] = outbound
		return nil
	})

	if err != nil {
		return nil, nil, fmt.Errorf("failed walking workspace docs: %w", err)
	}

	var nodeResults []DocNodeResult
	for relPath := range allDocPaths {
		outbound := outboundMap[relPath]
		inDeg := inDegreeMap[relPath]
		outDeg := len(outbound)
		isOrphan := (inDeg == 0 && outDeg == 0)

		nodeResults = append(nodeResults, DocNodeResult{
			FilePath:      relPath,
			InDegree:      inDeg,
			OutDegree:     outDeg,
			IsOrphan:      isOrphan,
			OutboundLinks: outbound,
		})
	}

	return nodeResults, brokenLinks, nil
}

// DefaultDocGraphStore implements DocGraphStore saving output to graphify-out/doc_graph_audit.json.
type DefaultDocGraphStore struct{}

// NewDefaultDocGraphStore initializes a DefaultDocGraphStore instance.
func NewDefaultDocGraphStore() *DefaultDocGraphStore {
	return &DefaultDocGraphStore{}
}

// SaveDocAuditReport saves the audited report atomically to targetPath using two-phase commit.
func (s *DefaultDocGraphStore) SaveDocAuditReport(ctx context.Context, outputPath string, report *DocAuditReport) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	cleanOutput := filepath.Clean(outputPath)
	dir := filepath.Dir(cleanOutput)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed creating output directory: %w", err)
	}

	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return fmt.Errorf("failed marshaling doc audit report: %w", err)
	}

	tmpFile := cleanOutput + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0644); err != nil {
		return fmt.Errorf("failed writing tmp doc audit report at %s: %w", tmpFile, err)
	}

	if err := os.Rename(tmpFile, cleanOutput); err != nil {
		_ = os.Remove(tmpFile)
		return fmt.Errorf("atomic rename failed for %s: %w", cleanOutput, err)
	}

	return nil
}

// DocGraphAuditor implements DocGraphAuditorService.
type DocGraphAuditor struct {
	workspaceRoot string
	parser        DocGraphParser
	store         DocGraphStore
	outputPath    string
}

// NewDocGraphAuditor initializes a new DocGraphAuditor compositor.
func NewDocGraphAuditor(workspaceRoot string) *DocGraphAuditor {
	cleanRoot := filepath.Clean(workspaceRoot)
	return &DocGraphAuditor{
		workspaceRoot: cleanRoot,
		parser:        NewDefaultDocGraphParser(cleanRoot),
		store:         NewDefaultDocGraphStore(),
		outputPath:    filepath.Join(cleanRoot, "graphify-out", "doc_graph_audit.json"),
	}
}

// AuditDocGraph executes full workspace documentation audit pass.
func (a *DocGraphAuditor) AuditDocGraph(ctx context.Context) (*DocAuditReport, error) {
	nodeResults, brokenLinks, err := a.parser.ParseWorkspaceDocs(ctx, a.workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("doc graph parsing failed: %w", err)
	}

	var orphanNodes []string
	totalEdges := 0

	for _, node := range nodeResults {
		totalEdges += node.OutDegree
		if node.IsOrphan {
			orphanNodes = append(orphanNodes, node.FilePath)
		}
	}

	report := &DocAuditReport{
		Timestamp:       time.Now().UTC(),
		TotalDocNodes:   len(nodeResults),
		TotalDocEdges:   totalEdges,
		OrphanCount:     len(orphanNodes),
		BrokenLinkCount: len(brokenLinks),
		OrphanNodes:     orphanNodes,
		BrokenLinks:     brokenLinks,
		NodeDetails:     nodeResults,
	}

	if saveErr := a.store.SaveDocAuditReport(ctx, a.outputPath, report); saveErr != nil {
		return report, fmt.Errorf("saving doc audit report failed: %w", saveErr)
	}

	return report, nil
}
