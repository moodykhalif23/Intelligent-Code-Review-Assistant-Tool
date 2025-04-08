package patterns

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/user/code-review-assistant/internal/models"
)

// AntiPattern represents a code anti-pattern to detect
type AntiPattern struct {
	Name        string
	Description string
	Category    string
	Severity    string
	Detector    func(fset *token.FileSet, node ast.Node) *models.Issue
}

// GetGoAntiPatterns returns a list of Go-specific code anti-patterns to detect
func GetGoAntiPatterns() []*AntiPattern {
	return []*AntiPattern{
		// Singleton pattern (often overused in Go)
		{
			Name:        "singleton-pattern",
			Description: "Singleton pattern usage",
			Category:    "anti-pattern",
			Severity:    "medium",
			Detector:    detectSingleton,
		},
		// Panic in non-main functions
		{
			Name:        "panic-usage",
			Description: "Use of panic in non-main functions",
			Category:    "anti-pattern",
			Severity:    "high",
			Detector:    detectPanic,
		},
		// Returning unexported types from exported functions
		{
			Name:        "unexported-return",
			Description: "Returning unexported types from exported functions",
			Category:    "anti-pattern",
			Severity:    "medium",
			Detector:    detectUnexportedReturn,
		},
		// Large interface anti-pattern
		{
			Name:        "large-interface",
			Description: "Interface with too many methods",
			Category:    "anti-pattern",
			Severity:    "medium",
			Detector:    detectLargeInterface,
		},
		// Empty interface without context
		{
			Name:        "empty-interface",
			Description: "Use of empty interface without clear context",
			Category:    "anti-pattern",
			Severity:    "low",
			Detector:    detectEmptyInterface,
		},
		// Goroutine without context or cancellation
		{
			Name:        "unmanaged-goroutine",
			Description: "Goroutine without context or cancellation mechanism",
			Category:    "anti-pattern",
			Severity:    "high",
			Detector:    detectUnmanagedGoroutine,
		},
		// Misuse of init function
		{
			Name:        "init-misuse",
			Description: "Misuse of init function for complex initialization",
			Category:    "anti-pattern",
			Severity:    "medium",
			Detector:    detectInitMisuse,
		},
	}
}

// detectSingleton detects singleton pattern usage
func detectSingleton(fset *token.FileSet, node ast.Node) *models.Issue {
	// Look for package-level variables with getter functions
	// This is a simplified implementation
	varDecl, ok := node.(*ast.GenDecl)
	if !ok || varDecl.Tok != token.VAR {
		return nil
	}

	for _, spec := range varDecl.Specs {
		valueSpec, ok := spec.(*ast.ValueSpec)
		if !ok || len(valueSpec.Names) == 0 {
			continue
		}

		varName := valueSpec.Names[0].Name
		if !valueSpec.Names[0].IsExported() && strings.HasPrefix(strings.ToLower(varName), "instance") {
			pos := fset.Position(varDecl.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Possible singleton pattern detected with variable '" + varName + "'",
				Category:   "anti-pattern",
				Severity:   "medium",
				Confidence: "medium",
				Suggestion: "Consider using dependency injection instead of singleton pattern",
				Rule:       "singleton-pattern",
			}
		}
	}

	return nil
}

// detectPanic detects use of panic in non-main functions
func detectPanic(fset *token.FileSet, node ast.Node) *models.Issue {
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}

	funcIdent, ok := callExpr.Fun.(*ast.Ident)
	if !ok || funcIdent.Name != "panic" {
		return nil
	}

	// Find the enclosing function
	var funcName string
	ast.Inspect(node, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			funcName = funcDecl.Name.Name
			return false
		}
		return true
	})

	// If not in main or init function, report issue
	if funcName != "main" && funcName != "init" {
		pos := fset.Position(callExpr.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Use of panic in function '" + funcName + "'",
			Category:   "anti-pattern",
			Severity:   "high",
			Confidence: "high",
			Suggestion: "Consider returning errors instead of using panic",
			Rule:       "panic-usage",
		}
	}

	return nil
}

// detectUnexportedReturn detects returning unexported types from exported functions
func detectUnexportedReturn(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectLargeInterface detects interfaces with too many methods
func detectLargeInterface(fset *token.FileSet, node ast.Node) *models.Issue {
	typeSpec, ok := node.(*ast.TypeSpec)
	if !ok {
		return nil
	}

	interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
	if !ok || interfaceType.Methods == nil {
		return nil
	}

	methodCount := len(interfaceType.Methods.List)
	if methodCount > 5 {
		pos := fset.Position(typeSpec.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Interface '" + typeSpec.Name.Name + "' has too many methods (" + string(rune('0'+methodCount)) + ")",
			Category:   "anti-pattern",
			Severity:   "medium",
			Confidence: "high",
			Suggestion: "Consider breaking down the interface into smaller, more focused interfaces",
			Rule:       "large-interface",
		}
	}

	return nil
}

// detectEmptyInterface detects use of empty interface without clear context
func detectEmptyInterface(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectUnmanagedGoroutine detects goroutines without context or cancellation
func detectUnmanagedGoroutine(fset *token.FileSet, node ast.Node) *models.Issue {
	goStmt, ok := node.(*ast.GoStmt)
	if !ok {
		return nil
	}

	// Check if the function or its parent function has a context parameter
	hasContext := false
	ast.Inspect(node, func(n ast.Node) bool {
		if funcDecl, ok := n.(*ast.FuncDecl); ok {
			if funcDecl.Type.Params != nil {
				for _, field := range funcDecl.Type.Params.List {
					if expr, ok := field.Type.(*ast.SelectorExpr); ok {
						if ident, ok := expr.X.(*ast.Ident); ok {
							if ident.Name == "context" {
								hasContext = true
								return false
							}
						}
					}
				}
			}
		}
		return true
	})

	if !hasContext {
		pos := fset.Position(goStmt.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Goroutine without context or cancellation mechanism",
			Category:   "anti-pattern",
			Severity:   "high",
			Confidence: "medium",
			Suggestion: "Use context.Context to manage goroutine lifecycle",
			Rule:       "unmanaged-goroutine",
		}
	}

	return nil
}

// detectInitMisuse detects misuse of init function
func detectInitMisuse(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Name.Name != "init" {
		return nil
	}

	// Check if init function is too complex
	if funcDecl.Body != nil {
		startPos := fset.Position(funcDecl.Body.Lbrace)
		endPos := fset.Position(funcDecl.Body.Rbrace)
		lineCount := endPos.Line - startPos.Line

		if lineCount > 10 {
			pos := fset.Position(funcDecl.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Complex init function with " + string(rune('0'+lineCount/10)) + string(rune('0'+lineCount%10)) + " lines",
				Category:   "anti-pattern",
				Severity:   "medium",
				Confidence: "medium",
				Suggestion: "Move complex initialization to dedicated functions that can be explicitly called",
				Rule:       "init-misuse",
			}
		}
	}

	return nil
}
