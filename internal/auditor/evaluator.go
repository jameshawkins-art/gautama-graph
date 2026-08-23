package auditor

import (
	"fmt"
	"go/ast"
	"strings"
)

// DefaultSelectorEvaluator implements SelectorEvaluator using go/ast traversal.
type DefaultSelectorEvaluator struct {
	maxASTDepth int
}

// NewDefaultSelectorEvaluator creates a new DefaultSelectorEvaluator with recursion depth limits.
func NewDefaultSelectorEvaluator(maxASTDepth int) *DefaultSelectorEvaluator {
	if maxASTDepth <= 0 {
		maxASTDepth = 50
	}
	return &DefaultSelectorEvaluator{
		maxASTDepth: maxASTDepth,
	}
}

// EvaluateSelector inspects an ast.File for explicit selector or function call expressions.
// It checks whether selectorIdent is called within the enclosing function/method named callerIdent,
// or as a receiver/package selector on callerIdent.
func (e *DefaultSelectorEvaluator) EvaluateSelector(file *ast.File, callerIdent, selectorIdent string) (bool, string, error) {
	if file == nil {
		return false, "", fmt.Errorf("file AST cannot be nil")
	}

	cleanCaller := strings.TrimSpace(callerIdent)
	cleanSelector := strings.TrimSpace(selectorIdent)
	if cleanSelector == "" {
		return false, "", fmt.Errorf("selectorIdent cannot be empty")
	}

	// 1. Check if callerIdent corresponds to a FuncDecl in file
	var targetFuncDecl *ast.FuncDecl
	if cleanCaller != "" {
		for _, decl := range file.Decls {
			if fn, ok := decl.(*ast.FuncDecl); ok {
				if fn.Name != nil && (fn.Name.Name == cleanCaller || strings.EqualFold(fn.Name.Name, cleanCaller)) {
					targetFuncDecl = fn
					break
				}
			}
		}
	}

	var rootNode ast.Node = file
	if targetFuncDecl != nil && targetFuncDecl.Body != nil {
		rootNode = targetFuncDecl.Body
	}

	var matched bool
	var matchedPattern string

	ast.Inspect(rootNode, func(n ast.Node) bool {
		if n == nil || matched {
			return false
		}

		// Look for CallExpr: e.g., t.Run(...), HelperAdd(...), s.Handle(...)
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		// Case A: SelectorExpr (e.g. t.Run, http.HandleFunc, receiver.Method)
		if sel, isSel := call.Fun.(*ast.SelectorExpr); isSel {
			if sel.Sel != nil && (sel.Sel.Name == cleanSelector || strings.EqualFold(sel.Sel.Name, cleanSelector)) {
				if ident, isIdent := sel.X.(*ast.Ident); isIdent {
					if cleanCaller == "" || ident.Name == cleanCaller || targetFuncDecl != nil {
						matched = true
						matchedPattern = fmt.Sprintf("ast.SelectorExpr(%s.%s)", ident.Name, sel.Sel.Name)
						return false
					}
				} else {
					matched = true
					matchedPattern = fmt.Sprintf("ast.SelectorExpr(_.%s)", sel.Sel.Name)
					return false
				}
			}
		}

		// Case B: Direct Ident CallExpr (e.g. HelperAdd(), at(), make())
		if ident, isIdent := call.Fun.(*ast.Ident); isIdent {
			if ident.Name == cleanSelector || strings.EqualFold(ident.Name, cleanSelector) {
				matched = true
				matchedPattern = fmt.Sprintf("ast.CallExpr(%s)", ident.Name)
				return false
			}
		}

		return true
	})

	// If not found in targetFuncDecl, do a fallback scan over the entire file
	if !matched && targetFuncDecl != nil {
		ast.Inspect(file, func(n ast.Node) bool {
			if n == nil || matched {
				return false
			}

			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			if sel, isSel := call.Fun.(*ast.SelectorExpr); isSel {
				if sel.Sel != nil && (sel.Sel.Name == cleanSelector || strings.EqualFold(sel.Sel.Name, cleanSelector)) {
					if ident, isIdent := sel.X.(*ast.Ident); isIdent {
						if ident.Name == cleanCaller || cleanCaller == "" {
							matched = true
							matchedPattern = fmt.Sprintf("ast.SelectorExpr(%s.%s)", ident.Name, sel.Sel.Name)
							return false
						}
					}
				}
			}

			if ident, isIdent := call.Fun.(*ast.Ident); isIdent {
				if ident.Name == cleanSelector || strings.EqualFold(ident.Name, cleanSelector) {
					matched = true
					matchedPattern = fmt.Sprintf("ast.CallExpr(%s)", ident.Name)
					return false
				}
			}

			return true
		})
	}

	return matched, matchedPattern, nil
}
