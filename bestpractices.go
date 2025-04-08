package bestpractices

import (
	"go/ast"
	"go/token"

	"github.com/user/code-review-assistant/internal/models"
)

// BestPractice represents a Go best practice to check
type BestPractice struct {
	Name        string
	Description string
	Category    string
	Severity    string
	Detector    func(fset *token.FileSet, node ast.Node) *models.Issue
}

// GetGoBestPractices returns a list of Go-specific best practices to check
func GetGoBestPractices() []*BestPractice {
	return []*BestPractice{
		// Error handling best practice
		{
			Name:        "error-handling",
			Description: "Proper error handling",
			Category:    "best-practice",
			Severity:    "high",
			Detector:    detectImproperErrorHandling,
		},
		// Context propagation
		{
			Name:        "context-propagation",
			Description: "Proper context propagation",
			Category:    "best-practice",
			Severity:    "high",
			Detector:    detectMissingContextPropagation,
		},
		// Interface segregation
		{
			Name:        "interface-segregation",
			Description: "Interface segregation principle",
			Category:    "best-practice",
			Severity:    "medium",
			Detector:    detectInterfaceSegregation,
		},
		// Defer usage
		{
			Name:        "defer-usage",
			Description: "Proper use of defer",
			Category:    "best-practice",
			Severity:    "medium",
			Detector:    detectImproperDeferUsage,
		},
		// Named return values
		{
			Name:        "named-returns",
			Description: "Proper use of named return values",
			Category:    "best-practice",
			Severity:    "low",
			Detector:    detectImproperNamedReturns,
		},
		// Package naming
		{
			Name:        "package-naming",
			Description: "Proper package naming",
			Category:    "best-practice",
			Severity:    "low",
			Detector:    detectImproperPackageNaming,
		},
		// Function naming
		{
			Name:        "function-naming",
			Description: "Proper function naming",
			Category:    "best-practice",
			Severity:    "low",
			Detector:    detectImproperFunctionNaming,
		},
		// Variable naming
		{
			Name:        "variable-naming",
			Description: "Proper variable naming",
			Category:    "best-practice",
			Severity:    "low",
			Detector:    detectImproperVariableNaming,
		},
	}
}

// detectImproperErrorHandling detects improper error handling
func detectImproperErrorHandling(fset *token.FileSet, node ast.Node) *models.Issue {
	// Look for ignored errors in assignment statements
	assignStmt, ok := node.(*ast.AssignStmt)
	if !ok {
		return nil
	}

	// Check for assignments where the right side is a function call
	// and the left side doesn't capture all return values
	if len(assignStmt.Rhs) == 1 {
		callExpr, ok := assignStmt.Rhs[0].(*ast.CallExpr)
		if !ok {
			return nil
		}

		// Try to determine if the function returns an error
		// This is a simplified check and would need type information for accuracy
		funcName := ""
		switch fun := callExpr.Fun.(type) {
		case *ast.Ident:
			funcName = fun.Name
		case *ast.SelectorExpr:
			if ident, ok := fun.X.(*ast.Ident); ok {
				funcName = ident.Name + "." + fun.Sel.Name
			}
		}

		// Common functions that return errors
		errorReturningFuncs := map[string]bool{
			"os.Open":       true,
			"ioutil.ReadFile": true,
			"json.Unmarshal": true,
			"io.Copy":       true,
			"http.Get":      true,
		}

		if errorReturningFuncs[funcName] && len(assignStmt.Lhs) < 2 {
			pos := fset.Position(assignStmt.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Error not handled from call to '" + funcName + "'",
				Category:   "best-practice",
				Severity:   "high",
				Confidence: "medium",
				Suggestion: "Capture and handle the error return value",
				Rule:       "error-handling",
			}
		}
	}

	return nil
}

// detectMissingContextPropagation detects missing context propagation
func detectMissingContextPropagation(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectInterfaceSegregation detects violations of interface segregation principle
func detectInterfaceSegregation(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectImproperDeferUsage detects improper use of defer
func detectImproperDeferUsage(fset *token.FileSet, node ast.Node) *models.Issue {
	deferStmt, ok := node.(*ast.DeferStmt)
	if !ok {
		return nil
	}

	// Check for deferred function calls that don't close resources
	callExpr, ok := deferStmt.Call.(*ast.CallExpr)
	if !ok {
		return nil
	}

	// Check if the deferred call is a method call
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	// Check if the method is Close()
	if selectorExpr.Sel.Name != "Close" {
		pos := fset.Position(deferStmt.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Defer used for function other than resource closing",
			Category:   "best-practice",
			Severity:   "low",
			Confidence: "low",
			Suggestion: "Defer is most commonly used for closing resources. Consider if this is the appropriate use case.",
			Rule:       "defer-usage",
		}
	}

	return nil
}

// detectImproperNamedReturns detects improper use of named return values
func detectImproperNamedReturns(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectImproperPackageNaming detects improper package naming
func detectImproperPackageNaming(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectImproperFunctionNaming detects improper function naming
func detectImproperFunctionNaming(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	name := funcDecl.Name.Name
	
	// Check for mixed case in unexported functions
	if !funcDecl.Name.IsExported() {
		for i, c := range name {
			if i > 0 && c >= 'A' && c <= 'Z' {
				pos := fset.Position(funcDecl.Pos())
				return &models.Issue{
					File:       pos.Filename,
					Line:       pos.Line,
					Column:     pos.Column,
					Message:    "Unexported function '" + name + "' uses mixed case",
					Category:   "best-practice",
					Severity:   "low",
					Confidence: "high",
					Suggestion: "Use camelCase for unexported functions",
					Rule:       "function-naming",
				}
			}
		}
	}

	return nil
}

// detectImproperVariableNaming detects improper variable naming
func detectImproperVariableNaming(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}
