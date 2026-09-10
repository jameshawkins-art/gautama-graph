package runner

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ErrGraphNotFound indicates that the target graphify-out/graph.json file does not exist.
var ErrGraphNotFound = errors.New("knowledge graph not found")

// ErrPathOutOfBounds indicates that a path parameter escapes the workspace boundary.
var ErrPathOutOfBounds = errors.New("path escapes workspace boundary")

// ValidatePathBoundary ensures targetPath is strictly contained within rootPath.
func ValidatePathBoundary(rootPath, targetPath string) error {
	cleanRoot := filepath.Clean(rootPath)
	cleanTarget := filepath.Clean(targetPath)

	absRoot, err := filepath.Abs(cleanRoot)
	if err != nil {
		return fmt.Errorf("failed resolving absolute root path: %w", err)
	}
	absTarget, err := filepath.Abs(cleanTarget)
	if err != nil {
		return fmt.Errorf("failed resolving absolute target path: %w", err)
	}

	if !strings.HasPrefix(absTarget, absRoot+string(filepath.Separator)) && absTarget != absRoot {
		return fmt.Errorf("%w: target %s escapes root %s", ErrPathOutOfBounds, cleanTarget, cleanRoot)
	}
	return nil
}

// DefaultQueryService coordinates Graphify binary discovery, validation, and execution.
type DefaultQueryService struct {
	binManager BinaryManager
	runner     SubprocessRunner
	mu         sync.Mutex
}

// NewDefaultQueryService initializes a new DefaultQueryService instance with injected dependencies.
func NewDefaultQueryService(binManager BinaryManager, runner SubprocessRunner) *DefaultQueryService {
	if binManager == nil {
		binManager = NewDefaultBinaryManager(nil)
	}
	if runner == nil {
		runner = NewDefaultSubprocessRunner()
	}
	return &DefaultQueryService{
		binManager: binManager,
		runner:     runner,
	}
}

// Query executes a BFS or DFS question traversal across the knowledge graph.
func (s *DefaultQueryService) Query(ctx context.Context, question string, opts QueryOptions) (*QueryResult, error) {
	trimmedQuestion := strings.TrimSpace(question)
	if trimmedQuestion == "" {
		return nil, errors.New("query question cannot be empty")
	}

	wsRoot := opts.WorkspaceRoot
	if wsRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed resolving current working directory: %w", err)
		}
		wsRoot = cwd
	}
	cleanWsRoot := filepath.Clean(wsRoot)

	graphPath := opts.GraphPath
	if graphPath == "" {
		graphPath = filepath.Join(cleanWsRoot, "graphify-out", "graph.json")
	} else if !filepath.IsAbs(graphPath) {
		graphPath = filepath.Join(cleanWsRoot, graphPath)
	}
	cleanGraphPath := filepath.Clean(graphPath)

	if err := ValidatePathBoundary(cleanWsRoot, cleanGraphPath); err != nil {
		return nil, err
	}

	if _, err := os.Stat(cleanGraphPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w at %s: run 'gautama-graph' or 'make graphify-update' first to generate graphify-out/graph.json", ErrGraphNotFound, cleanGraphPath)
		}
		return nil, fmt.Errorf("failed accessing graph path: %w", err)
	}

	s.mu.Lock()
	cfg := RunnerConfig{WorkspaceRootPath: cleanWsRoot}
	binPath, version, err := s.binManager.EnsureBinary(ctx, cfg)
	s.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("failed ensuring Graphify binary: %w", err)
	}

	args := []string{"query", trimmedQuestion}
	if opts.DFS {
		args = append(args, "--dfs")
	}
	if opts.BudgetTokens > 0 {
		args = append(args, "--budget", strconv.Itoa(opts.BudgetTokens))
	}
	for _, c := range opts.Context {
		args = append(args, "--context", c)
	}
	if opts.GraphPath != "" {
		args = append(args, "--graph", cleanGraphPath)
	}
	if opts.JSONOutput {
		args = append(args, "--json")
	}

	start := time.Now()
	stdout, _, execErr := s.runner.ExecuteCommand(ctx, binPath, cleanWsRoot, args...)
	if execErr != nil {
		return nil, fmt.Errorf("graph query execution failed: %w", execErr)
	}

	return &QueryResult{
		Output:        string(stdout),
		BinarySource:  "resolved",
		BinaryVersion: version,
		Duration:      time.Since(start),
		WorkspaceRoot: cleanWsRoot,
	}, nil
}

// Path finds the shortest topological path between two nodes in the knowledge graph.
func (s *DefaultQueryService) Path(ctx context.Context, nodeA, nodeB string, opts PathOptions) (*QueryResult, error) {
	trimmedA := strings.TrimSpace(nodeA)
	trimmedB := strings.TrimSpace(nodeB)
	if trimmedA == "" || trimmedB == "" {
		return nil, errors.New("both start and end nodes are required for path query")
	}

	wsRoot := opts.WorkspaceRoot
	if wsRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed resolving current working directory: %w", err)
		}
		wsRoot = cwd
	}
	cleanWsRoot := filepath.Clean(wsRoot)

	graphPath := opts.GraphPath
	if graphPath == "" {
		graphPath = filepath.Join(cleanWsRoot, "graphify-out", "graph.json")
	} else if !filepath.IsAbs(graphPath) {
		graphPath = filepath.Join(cleanWsRoot, graphPath)
	}
	cleanGraphPath := filepath.Clean(graphPath)

	if err := ValidatePathBoundary(cleanWsRoot, cleanGraphPath); err != nil {
		return nil, err
	}

	if _, err := os.Stat(cleanGraphPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w at %s: run 'gautama-graph' or 'make graphify-update' first to generate graphify-out/graph.json", ErrGraphNotFound, cleanGraphPath)
		}
		return nil, fmt.Errorf("failed accessing graph path: %w", err)
	}

	s.mu.Lock()
	cfg := RunnerConfig{WorkspaceRootPath: cleanWsRoot}
	binPath, version, err := s.binManager.EnsureBinary(ctx, cfg)
	s.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("failed ensuring Graphify binary: %w", err)
	}

	args := []string{"path", trimmedA, trimmedB}
	if opts.GraphPath != "" {
		args = append(args, "--graph", cleanGraphPath)
	}

	start := time.Now()
	stdout, _, execErr := s.runner.ExecuteCommand(ctx, binPath, cleanWsRoot, args...)
	if execErr != nil {
		return nil, fmt.Errorf("graph path execution failed: %w", execErr)
	}

	return &QueryResult{
		Output:        string(stdout),
		BinarySource:  "resolved",
		BinaryVersion: version,
		Duration:      time.Since(start),
		WorkspaceRoot: cleanWsRoot,
	}, nil
}

// Explain generates a plain-language explanation of a node and its adjacent neighborhood.
func (s *DefaultQueryService) Explain(ctx context.Context, concept string, opts ExplainOptions) (*QueryResult, error) {
	trimmedConcept := strings.TrimSpace(concept)
	if trimmedConcept == "" {
		return nil, errors.New("concept identifier cannot be empty for explain query")
	}

	wsRoot := opts.WorkspaceRoot
	if wsRoot == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("failed resolving current working directory: %w", err)
		}
		wsRoot = cwd
	}
	cleanWsRoot := filepath.Clean(wsRoot)

	graphPath := opts.GraphPath
	if graphPath == "" {
		graphPath = filepath.Join(cleanWsRoot, "graphify-out", "graph.json")
	} else if !filepath.IsAbs(graphPath) {
		graphPath = filepath.Join(cleanWsRoot, graphPath)
	}
	cleanGraphPath := filepath.Clean(graphPath)

	if err := ValidatePathBoundary(cleanWsRoot, cleanGraphPath); err != nil {
		return nil, err
	}

	if _, err := os.Stat(cleanGraphPath); err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("%w at %s: run 'gautama-graph' or 'make graphify-update' first to generate graphify-out/graph.json", ErrGraphNotFound, cleanGraphPath)
		}
		return nil, fmt.Errorf("failed accessing graph path: %w", err)
	}

	s.mu.Lock()
	cfg := RunnerConfig{WorkspaceRootPath: cleanWsRoot}
	binPath, version, err := s.binManager.EnsureBinary(ctx, cfg)
	s.mu.Unlock()
	if err != nil {
		return nil, fmt.Errorf("failed ensuring Graphify binary: %w", err)
	}

	args := []string{"explain", trimmedConcept}
	if opts.GraphPath != "" {
		args = append(args, "--graph", cleanGraphPath)
	}

	start := time.Now()
	stdout, _, execErr := s.runner.ExecuteCommand(ctx, binPath, cleanWsRoot, args...)
	if execErr != nil {
		return nil, fmt.Errorf("graph explain execution failed: %w", execErr)
	}

	return &QueryResult{
		Output:        string(stdout),
		BinarySource:  "resolved",
		BinaryVersion: version,
		Duration:      time.Since(start),
		WorkspaceRoot: cleanWsRoot,
	}, nil
}
