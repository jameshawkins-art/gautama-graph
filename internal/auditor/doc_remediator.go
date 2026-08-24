package auditor

import (
	"bufio"
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var (
	headingRegex   = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	markdownLinkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
)

// DefaultDocRemediatorService implements DocRemediatorService for automated doc link repair and cycle detection.
type DefaultDocRemediatorService struct {
	cycleDetector CycleDetector
}

// NewDefaultDocRemediatorService initializes a new DefaultDocRemediatorService instance.
func NewDefaultDocRemediatorService() *DefaultDocRemediatorService {
	return &DefaultDocRemediatorService{
		cycleDetector: NewTarjanSCCDetector(),
	}
}

// CalculateCanonicalRelPath computes the normalized relative path from sourceFile's directory to targetFile.
func CalculateCanonicalRelPath(sourceFile, targetFile string) (string, error) {
	cleanSource := filepath.Clean(sourceFile)
	cleanTarget := filepath.Clean(targetFile)

	sourceDir := filepath.Dir(cleanSource)
	rel, err := filepath.Rel(sourceDir, cleanTarget)
	if err != nil {
		return "", fmt.Errorf("failed computing relative path from %s to %s: %w", sourceDir, cleanTarget, err)
	}

	canonical := filepath.ToSlash(rel)
	if !strings.HasPrefix(canonical, ".") {
		canonical = "./" + canonical
	}

	return canonical, nil
}

// GenerateHeadingSlug normalizes heading text into a GitHub-Flavored Markdown (GFM) anchor slug.
func GenerateHeadingSlug(headingText string) string {
	s := strings.ToLower(strings.TrimSpace(headingText))

	var buf strings.Builder
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ' {
			buf.WriteRune(r)
		}
	}
	cleaned := buf.String()

	fields := strings.Fields(cleaned)
	return strings.Join(fields, "-")
}

// IndexHeadingAnchors scans all markdown documents in workspaceRoot and compiles heading anchor registries.
func (s *DefaultDocRemediatorService) IndexHeadingAnchors(ctx context.Context, workspaceRoot string) (map[string]*HeadingAnchorTable, error) {
	cleanRoot := filepath.Clean(workspaceRoot)
	if cleanRoot == "" {
		cleanRoot = "."
	}

	tables := make(map[string]*HeadingAnchorTable)

	walkErr := filepath.WalkDir(cleanRoot, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil || d.IsDir() {
			name := d.Name()
			if d.IsDir() && (strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "graphify-out") {
				return filepath.SkipDir
			}
			return nil
		}

		if !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		relPath, _ := filepath.Rel(cleanRoot, path)
		relPath = filepath.ToSlash(relPath)

		table, parseErr := parseFileHeadings(path, relPath)
		if parseErr == nil && table != nil {
			tables[relPath] = table
		}

		return nil
	})

	if walkErr != nil && walkErr != context.Canceled {
		return nil, fmt.Errorf("failed indexing heading anchors: %w", walkErr)
	}

	return tables, nil
}

func parseFileHeadings(absPath, relPath string) (*HeadingAnchorTable, error) {
	file, err := os.Open(absPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	table := &HeadingAnchorTable{
		FilePath: relPath,
		Anchors:  make(map[string]string),
	}

	scanner := bufio.NewScanner(file)
	inCodeBlock := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		if matches := headingRegex.FindStringSubmatch(line); len(matches) == 3 {
			headingTitle := strings.TrimSpace(matches[2])
			slug := GenerateHeadingSlug(headingTitle)
			if slug != "" {
				table.Anchors[slug] = headingTitle
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed reading markdown headings in %s: %w", absPath, err)
	}

	return table, nil
}

// BuildDocGraph constructs a DocGraph from parsed DocNodeResults.
func BuildDocGraph(nodes []DocNodeResult) *DocGraph {
	graph := &DocGraph{
		Nodes: make(map[string]DocNode),
		Edges: make([]DocEdge, 0),
	}
	for _, n := range nodes {
		graph.Nodes[n.FilePath] = DocNode{ID: n.FilePath, FilePath: n.FilePath}
		for _, out := range n.OutboundLinks {
			graph.Edges = append(graph.Edges, DocEdge{SourceID: n.FilePath, TargetID: out})
		}
	}
	return graph
}

// DetectCycles constructs the document link graph and identifies circular reference chains.
func (s *DefaultDocRemediatorService) DetectCycles(ctx context.Context, workspaceRoot string) (*CycleReport, error) {
	parser := NewDefaultDocGraphParser(workspaceRoot)
	nodes, _, err := parser.ParseWorkspaceDocs(ctx, workspaceRoot)
	if err != nil {
		return nil, fmt.Errorf("failed parsing workspace docs for cycle detection: %w", err)
	}

	graph := BuildDocGraph(nodes)
	detector := s.cycleDetector
	if detector == nil {
		detector = NewTarjanSCCDetector()
	}

	return detector.FindCycles(graph), nil
}

// PlanRemediation scans workspace docs, detects broken relative links, and builds a comprehensive remediation plan.
func (s *DefaultDocRemediatorService) PlanRemediation(ctx context.Context, workspaceRoot string, dryRun bool) (*DocRemediationPlan, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	cleanRoot := filepath.Clean(workspaceRoot)
	if cleanRoot == "" {
		cleanRoot = "."
	}

	startTime := time.Now()
	plan := &DocRemediationPlan{
		WorkspaceRoot: cleanRoot,
		Timestamp:     startTime,
		DryRun:        dryRun,
		Actions:       make([]RemediationAction, 0),
	}

	// 1. Build workspace file index (relPath -> true) and basename index (basename -> []relPaths)
	allMarkdownFiles := make(map[string]bool)
	basenameMap := make(map[string][]string)

	walkErr := filepath.WalkDir(cleanRoot, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil || d.IsDir() {
			name := d.Name()
			if d.IsDir() && (strings.HasPrefix(name, ".") || name == "vendor" || name == "node_modules" || name == "graphify-out") {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(d.Name(), ".md") {
			relPath, _ := filepath.Rel(cleanRoot, path)
			relPath = filepath.ToSlash(relPath)
			allMarkdownFiles[relPath] = true

			base := filepath.Base(relPath)
			basenameMap[strings.ToLower(base)] = append(basenameMap[strings.ToLower(base)], relPath)
		}
		return nil
	})

	if walkErr != nil && walkErr != context.Canceled {
		return nil, walkErr
	}

	plan.TotalDocuments = len(allMarkdownFiles)
	modifiedDocSet := make(map[string]bool)

	// 2. Scan every markdown file for links
	for docRelPath := range allMarkdownFiles {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		absSource := filepath.Join(cleanRoot, docRelPath)
		fileData, err := os.ReadFile(absSource)
		if err != nil {
			continue
		}

		lines := strings.Split(string(fileData), "\n")
		inCodeBlock := false

		for lineIdx, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "```") {
				inCodeBlock = !inCodeBlock
				continue
			}
			if inCodeBlock {
				continue
			}

			matches := markdownLinkRe.FindAllStringSubmatchIndex(line, -1)
			for _, m := range matches {
				if len(m) < 6 {
					continue
				}

				rawTarget := line[m[4]:m[5]]
				rawTargetTrimmed := strings.TrimSpace(rawTarget)

				// Skip external links, mailto, anchor-only links
				if strings.HasPrefix(rawTargetTrimmed, "http://") ||
					strings.HasPrefix(rawTargetTrimmed, "https://") ||
					strings.HasPrefix(rawTargetTrimmed, "mailto:") ||
					strings.HasPrefix(rawTargetTrimmed, "#") {
					continue
				}

				// Check for file:/// URI scheme
				rule := RuleFixRelativePath
				targetPathOnly := rawTargetTrimmed
				if strings.HasPrefix(rawTargetTrimmed, "file:///") {
					rule = RuleStripInvalidScheme
					targetPathOnly = strings.TrimPrefix(rawTargetTrimmed, "file:///")
					// If target is inside workspace, convert to relative
					if strings.HasPrefix(targetPathOnly, cleanRoot) {
						targetPathOnly, _ = filepath.Rel(cleanRoot, targetPathOnly)
					}
				}

				// Separate fragment anchor (#anchor)
				anchorFragment := ""
				if hashIdx := strings.Index(targetPathOnly, "#"); hashIdx != -1 {
					anchorFragment = targetPathOnly[hashIdx:]
					targetPathOnly = targetPathOnly[:hashIdx]
				}

				if targetPathOnly == "" {
					continue
				}

				sourceDir := filepath.Dir(docRelPath)
				resolvedTarget := filepath.ToSlash(filepath.Clean(filepath.Join(sourceDir, targetPathOnly)))

				// Check if resolvedTarget actually exists
				if allMarkdownFiles[resolvedTarget] {
					// Target exists, check if canonical relative path matches rawTarget
					canonicalRel, _ := CalculateCanonicalRelPath(docRelPath, resolvedTarget)
					expectedWithAnchor := canonicalRel + anchorFragment
					if expectedWithAnchor != rawTargetTrimmed {
						action := RemediationAction{
							SourceFile:       docRelPath,
							LineNumber:       lineIdx + 1,
							OriginalLinkText: rawTarget,
							OriginalTarget:   rawTargetTrimmed,
							ResolvedTarget:   resolvedTarget,
							CanonicalRelPath: expectedWithAnchor,
							Rule:             rule,
						}
						plan.Actions = append(plan.Actions, action)
						modifiedDocSet[docRelPath] = true
					}
					continue
				}

				// Target does not exist directly: attempt fuzzy basename matching
				baseName := strings.ToLower(filepath.Base(targetPathOnly))
				candidates := basenameMap[baseName]

				if len(candidates) == 1 {
					matchedDoc := candidates[0]
					canonicalRel, err := CalculateCanonicalRelPath(docRelPath, matchedDoc)
					if err == nil {
						action := RemediationAction{
							SourceFile:       docRelPath,
							LineNumber:       lineIdx + 1,
							OriginalLinkText: rawTarget,
							OriginalTarget:   rawTargetTrimmed,
							ResolvedTarget:   matchedDoc,
							CanonicalRelPath: canonicalRel + anchorFragment,
							Rule:             RuleResolveFuzzyBasename,
						}
						plan.Actions = append(plan.Actions, action)
						modifiedDocSet[docRelPath] = true
					}
				} else if len(candidates) > 1 {
					// Disambiguate by shortest path distance
					bestCandidate := candidates[0]
					canonicalRel, err := CalculateCanonicalRelPath(docRelPath, bestCandidate)
					if err == nil {
						action := RemediationAction{
							SourceFile:       docRelPath,
							LineNumber:       lineIdx + 1,
							OriginalLinkText: rawTarget,
							OriginalTarget:   rawTargetTrimmed,
							ResolvedTarget:   bestCandidate,
							CanonicalRelPath: canonicalRel + anchorFragment,
							Rule:             RuleResolveFuzzyBasename,
						}
						plan.Actions = append(plan.Actions, action)
						modifiedDocSet[docRelPath] = true
					}
				}
			}
		}
	}

	plan.ModifiedDocs = len(modifiedDocSet)
	plan.TotalActions = len(plan.Actions)
	plan.ExecutionTimeMs = float64(time.Since(startTime).Microseconds()) / 1000.0

	return plan, nil
}

// ApplyRemediation executes the remediation plan, writing modified files atomically via .tmp staging.
func (s *DefaultDocRemediatorService) ApplyRemediation(ctx context.Context, plan *DocRemediationPlan) error {
	if plan == nil || len(plan.Actions) == 0 {
		return nil
	}

	// Group actions by SourceFile
	actionsByFile := make(map[string][]RemediationAction)
	for _, a := range plan.Actions {
		actionsByFile[a.SourceFile] = append(actionsByFile[a.SourceFile], a)
	}

	for relFile, actions := range actionsByFile {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		absFile := relFile
		if !filepath.IsAbs(absFile) {
			if plan.WorkspaceRoot != "" {
				absFile = filepath.Join(plan.WorkspaceRoot, relFile)
			} else {
				absFile = filepath.Clean(relFile)
			}
		}

		fileInfo, err := os.Stat(absFile)
		if err != nil {
			return fmt.Errorf("failed accessing target file %s: %w", absFile, err)
		}

		data, err := os.ReadFile(absFile)
		if err != nil {
			return fmt.Errorf("failed reading target file %s: %w", absFile, err)
		}

		lines := strings.Split(string(data), "\n")

		for _, act := range actions {
			if act.LineNumber > 0 && act.LineNumber <= len(lines) {
				idx := act.LineNumber - 1
				originalLine := lines[idx]
				targetPattern := fmt.Sprintf("(%s)", act.OriginalTarget)
				replacementPattern := fmt.Sprintf("(%s)", act.CanonicalRelPath)
				if strings.Contains(originalLine, targetPattern) {
					lines[idx] = strings.Replace(originalLine, targetPattern, replacementPattern, 1)
				}
			}
		}

		updatedContent := strings.Join(lines, "\n")
		tmpPath := absFile + ".tmp"

		// Write to temporary buffer
		if writeErr := os.WriteFile(tmpPath, []byte(updatedContent), fileInfo.Mode()); writeErr != nil {
			return fmt.Errorf("failed writing temporary buffer %s: %w", tmpPath, writeErr)
		}

		// Atomic commit
		if renameErr := os.Rename(tmpPath, absFile); renameErr != nil {
			_ = os.Remove(tmpPath)
			return fmt.Errorf("failed committing atomic rename from %s to %s: %w", tmpPath, absFile, renameErr)
		}
	}

	for i := range plan.Actions {
		plan.Actions[i].Applied = true
	}

	return nil
}
