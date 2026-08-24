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

// DefaultCrossPackageEvaluator implements CrossPackageEvaluator.
type DefaultCrossPackageEvaluator struct {
	indexer PackageSymbolIndexer
	fset    *token.FileSet
}

// NewDefaultCrossPackageEvaluator initializes a new DefaultCrossPackageEvaluator instance.
func NewDefaultCrossPackageEvaluator(indexer PackageSymbolIndexer) *DefaultCrossPackageEvaluator {
	return &DefaultCrossPackageEvaluator{
		indexer: indexer,
		fset:    token.NewFileSet(),
	}
}

// ResolveFileImports maps all import identifiers and aliases in a Go file to their package paths.
func (e *DefaultCrossPackageEvaluator) ResolveFileImports(file *ast.File) map[string]string {
	imports := make(map[string]string)
	if file == nil {
		return imports
	}

	for _, imp := range file.Imports {
		if imp.Path == nil {
			continue
		}
		importPath := strings.Trim(imp.Path.Value, "\"")

		if imp.Name != nil {
			alias := imp.Name.Name
			if alias != "_" {
				imports[alias] = importPath
			}
		} else {
			defaultAlias := filepath.Base(importPath)
			imports[defaultAlias] = importPath
		}
	}

	return imports
}

// EvaluateCrossPackageCall checks if a caller function in sourceFile invokes targetSymbol in targetPkg.
func (e *DefaultCrossPackageEvaluator) EvaluateCrossPackageCall(ctx context.Context, sourceFile, callerSymbol, targetPkg, targetSymbol string) (bool, ProvenanceStatus, string, error) {
	select {
	case <-ctx.Done():
		return false, ProvenancePrunedPhantom, "", ctx.Err()
	default:
	}

	cleanSource := filepath.Clean(sourceFile)
	srcData, err := os.ReadFile(cleanSource)
	if err != nil {
		return false, ProvenancePrunedPhantom, "", fmt.Errorf("failed reading source file %s: %w", cleanSource, err)
	}

	fset := token.NewFileSet()
	fileAST, err := parser.ParseFile(fset, cleanSource, srcData, parser.ParseComments)
	if err != nil {
		return false, ProvenancePrunedPhantom, "", fmt.Errorf("failed parsing source file AST %s: %w", cleanSource, err)
	}

	imports := e.ResolveFileImports(fileAST)

	// Check if target package is imported
	targetClean := filepath.ToSlash(filepath.Clean(targetPkg))
	targetBase := filepath.Base(targetClean)

	var matchingAliases []string
	hasDotImport := false

	if targetPkg != "" && targetPkg != "." {
		for alias, path := range imports {
			cleanImp := filepath.ToSlash(filepath.Clean(path))
			if cleanImp == targetClean || filepath.Base(cleanImp) == targetBase || strings.HasSuffix(cleanImp, "/"+targetBase) {
				if alias == "." {
					hasDotImport = true
				} else {
					matchingAliases = append(matchingAliases, alias)
				}
			}
		}
	} else if e.indexer != nil {
		// When targetPkg is unspecified, search across all imports for targetSymbol
		for alias, path := range imports {
			if table, found := e.indexer.GetPackageTable(path); found && table != nil {
				if _, exists := table.Symbols[targetSymbol]; exists {
					if alias == "." {
						hasDotImport = true
					} else {
						matchingAliases = append(matchingAliases, alias)
					}
				}
			}
		}
	}

	// If no imported packages matched
	if len(matchingAliases) == 0 && !hasDotImport {
		return false, ProvenancePrunedPhantom, "", nil
	}

	// Locate caller function/method in source file
	var callerDecl *ast.FuncDecl
	for _, decl := range fileAST.Decls {
		if fn, ok := decl.(*ast.FuncDecl); ok {
			if fn.Name.Name == callerSymbol {
				callerDecl = fn
				break
			}
		}
	}

	if callerDecl == nil {
		// Caller symbol not found in file
		return false, ProvenancePrunedPhantom, "", nil
	}

	// Walk caller AST body for selector expressions
	var matchedPattern string
	isMatch := false

	ast.Inspect(callerDecl, func(n ast.Node) bool {
		if isMatch || n == nil {
			return !isMatch
		}

		switch expr := n.(type) {
		case *ast.SelectorExpr:
			if ident, ok := expr.X.(*ast.Ident); ok {
				for _, alias := range matchingAliases {
					if ident.Name == alias && expr.Sel.Name == targetSymbol {
						isMatch = true
						matchedPattern = fmt.Sprintf("cross_pkg_selector: %s.%s", ident.Name, expr.Sel.Name)
						return false
					}
				}
			}
		case *ast.CallExpr:
			if ident, ok := expr.Fun.(*ast.Ident); ok && hasDotImport {
				if ident.Name == targetSymbol {
					isMatch = true
					matchedPattern = fmt.Sprintf("dot_import_call: %s()", ident.Name)
					return false
				}
			}
		}

		return true
	})

	if isMatch {
		return true, ProvenanceResolvedCrossPackage, matchedPattern, nil
	}

	return false, ProvenancePrunedPhantom, "", nil
}
