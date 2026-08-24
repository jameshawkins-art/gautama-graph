package auditor

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

// DefaultInterfaceResolver implements InterfaceResolver using indexed package method sets.
type DefaultInterfaceResolver struct {
	indexer PackageSymbolIndexer
}

// NewDefaultInterfaceResolver initializes a new DefaultInterfaceResolver instance.
func NewDefaultInterfaceResolver(indexer PackageSymbolIndexer) *DefaultInterfaceResolver {
	return &DefaultInterfaceResolver{
		indexer: indexer,
	}
}

// CheckImplementation determines if a concrete struct satisfies all methods required by an interface.
func (r *DefaultInterfaceResolver) CheckImplementation(ctx context.Context, concretePkg, concreteType, ifacePkg, ifaceName string) (bool, *InterfaceBinding, error) {
	select {
	case <-ctx.Done():
		return false, nil, ctx.Err()
	default:
	}

	cleanConcreteType := strings.TrimPrefix(concreteType, "*")

	// 1. Retrieve Interface Symbol Table
	var ifaceTable *PackageSymbolTable
	var actualIfacePkg string

	if ifacePkg != "" && ifacePkg != "." {
		if table, found := r.indexer.GetPackageTable(ifacePkg); found {
			ifaceTable = table
			actualIfacePkg = ifacePkg
		} else if table, found := r.indexer.GetPackageTable(filepath.Base(ifacePkg)); found {
			ifaceTable = table
			actualIfacePkg = ifacePkg
		} else {
			return false, nil, fmt.Errorf("interface package %s not found in symbol index", ifacePkg)
		}
	} else if r.indexer != nil {
		// Search all packages for ifaceName
		for pkgPath, table := range r.indexer.GetAllPackages() {
			if sym, exists := table.Symbols[ifaceName]; exists && sym.Kind == KindInterface {
				ifaceTable = table
				actualIfacePkg = pkgPath
				break
			}
		}
	}

	if ifaceTable == nil {
		return false, nil, fmt.Errorf("interface %s not found in symbol index", ifaceName)
	}

	sym, ifaceExists := ifaceTable.Symbols[ifaceName]
	if !ifaceExists || sym.Kind != KindInterface {
		return false, nil, fmt.Errorf("interface %s not found in package %s", ifaceName, actualIfacePkg)
	}

	requiredMethods := ifaceTable.MethodSets[ifaceName]

	// 2. Retrieve Concrete Struct Symbol Table
	var concreteTable *PackageSymbolTable
	var actualConcretePkg string

	if concretePkg != "" && concretePkg != "." {
		if table, found := r.indexer.GetPackageTable(concretePkg); found {
			concreteTable = table
			actualConcretePkg = concretePkg
		} else if table, found := r.indexer.GetPackageTable(filepath.Base(concretePkg)); found {
			concreteTable = table
			actualConcretePkg = concretePkg
		} else {
			return false, nil, fmt.Errorf("concrete package %s not found in symbol index", concretePkg)
		}
	} else if r.indexer != nil {
		// Search all packages for cleanConcreteType
		for pkgPath, table := range r.indexer.GetAllPackages() {
			if sym, exists := table.Symbols[cleanConcreteType]; exists && (sym.Kind == KindStruct || sym.Kind == KindTypeAlias) {
				concreteTable = table
				actualConcretePkg = pkgPath
				break
			}
		}
	}

	if concreteTable == nil {
		return false, nil, fmt.Errorf("concrete struct %s not found in symbol index", cleanConcreteType)
	}

	structSym, structExists := concreteTable.Symbols[cleanConcreteType]
	if !structExists || (structSym.Kind != KindStruct && structSym.Kind != KindTypeAlias) {
		return false, nil, fmt.Errorf("concrete struct %s not found in package %s", cleanConcreteType, actualConcretePkg)
	}

	concreteMethods := concreteTable.MethodSets[cleanConcreteType]
	concreteMethodMap := make(map[string]bool)
	for _, m := range concreteMethods {
		concreteMethodMap[m] = true
	}

	// 3. Verify Interface Method Satisfaction
	var matchedMethods []string
	for _, req := range requiredMethods {
		if !concreteMethodMap[req] {
			// Missing required interface method
			return false, nil, nil
		}
		matchedMethods = append(matchedMethods, req)
	}

	binding := &InterfaceBinding{
		InterfacePackage: actualIfacePkg,
		InterfaceName:    ifaceName,
		ConcretePackage:  actualConcretePkg,
		ConcreteTypeName: cleanConcreteType,
		MatchedMethods:   matchedMethods,
	}

	return true, binding, nil
}

// FindImplementations searches all indexed workspace packages for structs satisfying the interface.
func (r *DefaultInterfaceResolver) FindImplementations(ctx context.Context, ifacePkg, ifaceName string) ([]InterfaceBinding, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	var results []InterfaceBinding
	if r.indexer == nil {
		return results, nil
	}

	for pkgPath, table := range r.indexer.GetAllPackages() {
		for symName, sym := range table.Symbols {
			if sym.Kind == KindStruct {
				if isImpl, binding, err := r.CheckImplementation(ctx, pkgPath, symName, ifacePkg, ifaceName); err == nil && isImpl {
					results = append(results, *binding)
				}
			}
		}
	}

	return results, nil
}
