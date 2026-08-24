package auditor

import (
	"context"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// DefaultPackageSymbolIndexer traverses the workspace and indexes Go package compilation units.
type DefaultPackageSymbolIndexer struct {
	mu           sync.RWMutex
	fset         *token.FileSet
	packageIndex map[string]*PackageSymbolTable
}

// NewDefaultPackageSymbolIndexer initializes a new DefaultPackageSymbolIndexer instance.
func NewDefaultPackageSymbolIndexer() *DefaultPackageSymbolIndexer {
	return &DefaultPackageSymbolIndexer{
		fset:         token.NewFileSet(),
		packageIndex: make(map[string]*PackageSymbolTable),
	}
}

// IndexWorkspace traverses workspaceRoot, parses all Go packages, and indexes exported symbols.
func (idx *DefaultPackageSymbolIndexer) IndexWorkspace(ctx context.Context, workspaceRoot string) (map[string]*PackageSymbolTable, error) {
	cleanRoot := filepath.Clean(workspaceRoot)
	if cleanRoot == "" {
		cleanRoot = "."
	}

	idx.mu.Lock()
	defer idx.mu.Unlock()

	idx.fset = token.NewFileSet()
	idx.packageIndex = make(map[string]*PackageSymbolTable)

	visitedDirs := make(map[string]bool)

	walkErr := filepath.WalkDir(cleanRoot, func(path string, d fs.DirEntry, err error) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if err != nil {
			return nil // Skip unreadable paths gracefully
		}

		if !d.IsDir() {
			return nil
		}

		name := d.Name()
		if strings.HasPrefix(name, ".") && name != "." {
			return filepath.SkipDir
		}
		if name == "vendor" || name == "node_modules" || name == "graphify-out" {
			return filepath.SkipDir
		}

		cleanPath := filepath.Clean(path)
		if visitedDirs[cleanPath] {
			return nil
		}
		visitedDirs[cleanPath] = true

		if _, err := ValidatePathBoundary(cleanRoot, cleanPath); err != nil {
			return filepath.SkipDir
		}

		// Parse all Go files in this directory
		pkgs, parseErr := parser.ParseDir(idx.fset, cleanPath, func(fi os.FileInfo) bool {
			return !strings.HasSuffix(fi.Name(), "_test.go")
		}, parser.ParseComments)

		if parseErr != nil || len(pkgs) == 0 {
			return nil
		}

		for pkgName, pkgAST := range pkgs {
			relDir, _ := filepath.Rel(cleanRoot, cleanPath)
			if relDir == "" {
				relDir = "."
			}
			relDir = filepath.ToSlash(relDir)

			table := &PackageSymbolTable{
				PackageName: pkgName,
				PackagePath: relDir,
				Directory:   cleanPath,
				Symbols:     make(map[string]ExportedSymbol),
				MethodSets:  make(map[string][]string),
				FileSet:     idx.fset,
				PackageAST:  pkgAST,
			}

			idx.extractPackageSymbols(table, pkgAST, cleanRoot)

			// Index by relative directory path and package name
			idx.packageIndex[relDir] = table
			idx.packageIndex[pkgName] = table
			if strings.HasPrefix(relDir, "internal/") || strings.HasPrefix(relDir, "cmd/") {
				idx.packageIndex[filepath.Base(relDir)] = table
			}
		}

		return nil
	})

	if walkErr != nil && walkErr != context.Canceled {
		return nil, fmt.Errorf("failed indexing workspace at %s: %w", cleanRoot, walkErr)
	}

	return idx.packageIndex, nil
}

func (idx *DefaultPackageSymbolIndexer) extractPackageSymbols(table *PackageSymbolTable, pkg *ast.Package, workspaceRoot string) {
	for fileName, file := range pkg.Files {
		relFile, _ := filepath.Rel(workspaceRoot, fileName)
		relFile = filepath.ToSlash(relFile)

		for _, decl := range file.Decls {
			switch d := decl.(type) {
			case *ast.FuncDecl:
				funcName := d.Name.Name
				if !ast.IsExported(funcName) {
					continue
				}

				pos := idx.fset.Position(d.Pos())
				docSummary := ""
				if d.Doc != nil {
					docSummary = strings.TrimSpace(d.Doc.Text())
				}

				if d.Recv == nil {
					// Standalone function
					table.Symbols[funcName] = ExportedSymbol{
						Name:        funcName,
						Kind:        KindFunction,
						PackagePath: table.PackagePath,
						FilePath:    relFile,
						LineNumber:  pos.Line,
						DocSummary:  docSummary,
					}
				} else {
					// Method on struct receiver
					recvType := extractReceiverTypeName(d.Recv)
					if recvType != "" {
						table.MethodSets[recvType] = append(table.MethodSets[recvType], funcName)
						table.Symbols[funcName] = ExportedSymbol{
							Name:        funcName,
							Kind:        KindMethod,
							Receiver:    recvType,
							PackagePath: table.PackagePath,
							FilePath:    relFile,
							LineNumber:  pos.Line,
							DocSummary:  docSummary,
						}
					}
				}

			case *ast.GenDecl:
				for _, spec := range d.Specs {
					switch s := spec.(type) {
					case *ast.TypeSpec:
						typeName := s.Name.Name
						if !ast.IsExported(typeName) {
							continue
						}

						pos := idx.fset.Position(s.Pos())
						docSummary := ""
						if d.Doc != nil {
							docSummary = strings.TrimSpace(d.Doc.Text())
						}

						kind := KindStruct
						switch s.Type.(type) {
						case *ast.InterfaceType:
							kind = KindInterface
							// Extract interface method requirements
							if iface, ok := s.Type.(*ast.InterfaceType); ok && iface.Methods != nil {
								var ifaceMethods []string
								for _, method := range iface.Methods.List {
									for _, mName := range method.Names {
										if ast.IsExported(mName.Name) {
											ifaceMethods = append(ifaceMethods, mName.Name)
										}
									}
								}
								table.MethodSets[typeName] = ifaceMethods
							}
						case *ast.Ident, *ast.SelectorExpr:
							kind = KindTypeAlias
						}

						table.Symbols[typeName] = ExportedSymbol{
							Name:        typeName,
							Kind:        kind,
							PackagePath: table.PackagePath,
							FilePath:    relFile,
							LineNumber:  pos.Line,
							DocSummary:  docSummary,
						}

					case *ast.ValueSpec:
						for _, ident := range s.Names {
							if ast.IsExported(ident.Name) {
								pos := idx.fset.Position(ident.Pos())
								kind := KindVariable
								if d.Tok == token.CONST {
									kind = KindConstant
								}
								table.Symbols[ident.Name] = ExportedSymbol{
									Name:        ident.Name,
									Kind:        kind,
									PackagePath: table.PackagePath,
									FilePath:    relFile,
									LineNumber:  pos.Line,
								}
							}
						}
					}
				}
			}
		}
	}
}

func extractReceiverTypeName(recv *ast.FieldList) string {
	if recv == nil || len(recv.List) == 0 {
		return ""
	}
	typeExpr := recv.List[0].Type
	if star, ok := typeExpr.(*ast.StarExpr); ok {
		typeExpr = star.X
	}
	if ident, ok := typeExpr.(*ast.Ident); ok {
		return ident.Name
	}
	return ""
}

// GetPackageTable returns the cached symbol table for a package path or package name.
func (idx *DefaultPackageSymbolIndexer) GetPackageTable(packagePath string) (*PackageSymbolTable, bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	cleanPath := filepath.ToSlash(filepath.Clean(packagePath))
	if table, found := idx.packageIndex[cleanPath]; found {
		return table, true
	}

	// Try lookup by base name
	base := filepath.Base(cleanPath)
	if table, found := idx.packageIndex[base]; found {
		return table, true
	}

	return nil, false
}

// GetAllPackages returns a shallow copy of all indexed package symbol tables.
func (idx *DefaultPackageSymbolIndexer) GetAllPackages() map[string]*PackageSymbolTable {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	copyMap := make(map[string]*PackageSymbolTable, len(idx.packageIndex))
	for k, v := range idx.packageIndex {
		copyMap[k] = v
	}
	return copyMap
}
