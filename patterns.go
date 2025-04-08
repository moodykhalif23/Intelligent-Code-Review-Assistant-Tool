package patterns

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/user/code-review-assistant/internal/models"
)

// Pattern represents a code pattern to detect
type Pattern struct {
	Name        string
	Description string
	Category    string
	Severity    string
	Detector    func(fset *token.FileSet, node ast.Node) *models.Issue
}

// GetGoPatterns returns a list of Go-specific code patterns to detect
func GetGoPatterns() []*Pattern {
	return []*Pattern{
		// Empty function pattern
		{
			Name:        "empty-function",
			Description: "Function with an empty body",
			Category:    "code-smell",
			Severity:    "low",
			Detector:    detectEmptyFunction,
		},
		// Too many parameters pattern
		{
			Name:        "too-many-params",
			Description: "Function with too many parameters",
			Category:    "code-smell",
			Severity:    "medium",
			Detector:    detectTooManyParams,
		},
		// Long function pattern
		{
			Name:        "long-function",
			Description: "Function that is too long",
			Category:    "code-smell",
			Severity:    "medium",
			Detector:    detectLongFunction,
		},
		// Deeply nested code pattern
		{
			Name:        "deep-nesting",
			Description: "Deeply nested control structures",
			Category:    "code-smell",
			Severity:    "medium",
			Detector:    detectDeepNesting,
		},
		// Naked return pattern
		{
			Name:        "naked-return",
			Description: "Naked return in a function with named return values",
			Category:    "code-smell",
			Severity:    "low",
			Detector:    detectNakedReturn,
		},
		// Unused parameter pattern
		{
			Name:        "unused-param",
			Description: "Unused function parameter",
			Category:    "code-smell",
			Severity:    "low",
			Detector:    detectUnusedParam,
		},
		// Boolean parameter pattern
		{
			Name:        "boolean-param",
			Description: "Boolean parameter in function signature",
			Category:    "code-smell",
			Severity:    "low",
			Detector:    detectBooleanParam,
		},
		// Magic number pattern
		{
			Name:        "magic-number",
			Description: "Magic number in code",
			Category:    "code-smell",
			Severity:    "low",
			Detector:    detectMagicNumber,
		},
		// Exported function without comment
		{
			Name:        "undocumented-exported",
			Description: "Exported function without documentation",
			Category:    "documentation",
			Severity:    "medium",
			Detector:    detectUndocumentedExported,
		},
		// Inefficient string concatenation
		{
			Name:        "inefficient-string-concat",
			Description: "Inefficient string concatenation in a loop",
			Category:    "performance",
			Severity:    "medium",
			Detector:    detectInefficientStringConcat,
		},
	}
}

// detectEmptyFunction detects functions with empty bodies
func detectEmptyFunction(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	if funcDecl.Body != nil && len(funcDecl.Body.List) == 0 {
		pos := fset.Position(funcDecl.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Function '" + funcDecl.Name.Name + "' has an empty body",
			Category:   "code-smell",
			Severity:   "low",
			Confidence: "high",
			Suggestion: "Consider implementing the function or removing it if not needed",
			Rule:       "empty-function",
		}
	}

	return nil
}

// detectTooManyParams detects functions with too many parameters
func detectTooManyParams(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 5 {
		pos := fset.Position(funcDecl.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Function '" + funcDecl.Name.Name + "' has too many parameters (" + string(rune('0'+len(funcDecl.Type.Params.List))) + ")",
			Category:   "code-smell",
			Severity:   "medium",
			Confidence: "high",
			Suggestion: "Consider refactoring to use a struct for parameters",
			Rule:       "too-many-params",
		}
	}

	return nil
}

// detectLongFunction detects functions that are too long
func detectLongFunction(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Body == nil {
		return nil
	}

	startPos := fset.Position(funcDecl.Body.Lbrace)
	endPos := fset.Position(funcDecl.Body.Rbrace)
	lineCount := endPos.Line - startPos.Line

	if lineCount > 50 {
		pos := fset.Position(funcDecl.Pos())
		return &models.Issue{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Message:    "Function '" + funcDecl.Name.Name + "' is too long (" + string(rune('0'+lineCount/10)) + string(rune('0'+lineCount%10)) + " lines)",
			Category:   "code-smell",
			Severity:   "medium",
			Confidence: "high",
			Suggestion: "Consider breaking down the function into smaller, more focused functions",
			Rule:       "long-function",
		}
	}

	return nil
}

// detectDeepNesting detects deeply nested control structures
func detectDeepNesting(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectNakedReturn detects naked returns in functions with named return values
func detectNakedReturn(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectUnusedParam detects unused function parameters
func detectUnusedParam(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectBooleanParam detects boolean parameters in function signatures
func detectBooleanParam(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok || funcDecl.Type.Params == nil {
		return nil
	}

	for _, field := range funcDecl.Type.Params.List {
		ident, ok := field.Type.(*ast.Ident)
		if ok && ident.Name == "bool" {
			pos := fset.Position(field.Pos())
			paramName := ""
			if len(field.Names) > 0 {
				paramName = field.Names[0].Name
			}
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Function '" + funcDecl.Name.Name + "' has boolean parameter '" + paramName + "'",
				Category:   "code-smell",
				Severity:   "low",
				Confidence: "medium",
				Suggestion: "Consider using an enum type or constants for better readability and extensibility",
				Rule:       "boolean-param",
			}
		}
	}

	return nil
}

// detectMagicNumber detects magic numbers in code
func detectMagicNumber(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectUndocumentedExported detects exported functions without documentation
func detectUndocumentedExported(fset *token.FileSet, node ast.Node) *models.Issue {
	funcDecl, ok := node.(*ast.FuncDecl)
	if !ok {
		return nil
	}

	// Check if function is exported (starts with uppercase letter)
	if funcDecl.Name.IsExported() {
		// Check if function has a doc comment
		if funcDecl.Doc == nil || len(funcDecl.Doc.List) == 0 {
			pos := fset.Position(funcDecl.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Exported function '" + funcDecl.Name.Name + "' lacks documentation",
				Category:   "documentation",
				Severity:   "medium",
				Confidence: "high",
				Suggestion: "Add documentation comments to describe the function's purpose, parameters, and return values",
				Rule:       "undocumented-exported",
			}
		}
	}

	return nil
}

// detectInefficientStringConcat detects inefficient string concatenation in loops
func detectInefficientStringConcat(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}
