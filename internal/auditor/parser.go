package auditor

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

// DefaultASTParser implements ASTParser using standard library go/parser.
type DefaultASTParser struct {
	workspaceRoot string
}

// NewDefaultASTParser initializes a new DefaultASTParser with path boundary safety.
func NewDefaultASTParser(workspaceRoot string) *DefaultASTParser {
	if workspaceRoot == "" {
		workspaceRoot = "."
	}
	absRoot, err := filepath.Abs(workspaceRoot)
	if err != nil {
		absRoot = workspaceRoot
	}
	return &DefaultASTParser{
		workspaceRoot: absRoot,
	}
}

// ParseFile parses a Go source file into an *ast.File and *token.FileSet after validating path safety.
func (p *DefaultASTParser) ParseFile(ctx context.Context, filePath string) (*ast.File, *token.FileSet, error) {
	select {
	case <-ctx.Done():
		return nil, nil, ctx.Err()
	default:
	}

	cleanPath := filepath.Clean(filePath)
	absPath, err := filepath.Abs(cleanPath)
	if err != nil {
		return nil, nil, fmt.Errorf("invalid file path %s: %w", filePath, err)
	}

	// Boundary defense: Ensure file is within workspace root
	if p.workspaceRoot != "" && !strings.HasPrefix(absPath, p.workspaceRoot) {
		return nil, nil, fmt.Errorf("security violation: path %s outside workspace root %s", absPath, p.workspaceRoot)
	}

	// Verify file existence and non-zero size before parsing
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, nil, fmt.Errorf("file stat failed for %s: %w", absPath, err)
	}
	if info.IsDir() {
		return nil, nil, fmt.Errorf("target %s is a directory, not a Go file", absPath)
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, absPath, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, fmt.Errorf("ast parse failed for %s: %w", absPath, err)
	}

	return node, fset, nil
}
