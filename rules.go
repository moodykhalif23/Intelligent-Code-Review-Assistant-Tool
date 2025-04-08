package optimization

import (
	"go/ast"
	"go/token"
	"strings"

	"github.com/user/code-review-assistant/internal/models"
)

// OptimizationRule represents a rule for detecting optimization opportunities
type OptimizationRule struct {
	ID          string
	Name        string
	Description string
	Detector    func(fset *token.FileSet, node ast.Node) *models.Optimization
}

// GetOptimizationRules returns a list of optimization rules
func GetOptimizationRules() []*OptimizationRule {
	return []*OptimizationRule{
		// Inefficient string concatenation in loops
		{
			ID:          "OPT001",
			Name:        "inefficient-string-concat",
			Description: "Inefficient string concatenation in loops",
			Detector:    detectInefficientStringConcat,
		},
		// Unnecessary memory allocations
		{
			ID:          "OPT002",
			Name:        "unnecessary-allocation",
			Description: "Unnecessary memory allocations",
			Detector:    detectUnnecessaryAllocation,
		},
		// Suboptimal slice capacity
		{
			ID:          "OPT003",
			Name:        "suboptimal-slice-capacity",
			Description: "Suboptimal slice capacity",
			Detector:    detectSuboptimalSliceCapacity,
		},
		// Inefficient map initialization
		{
			ID:          "OPT004",
			Name:        "inefficient-map-init",
			Description: "Inefficient map initialization",
			Detector:    detectInefficientMapInit,
		},
		// Redundant type conversions
		{
			ID:          "OPT005",
			Name:        "redundant-type-conversion",
			Description: "Redundant type conversions",
			Detector:    detectRedundantTypeConversion,
		},
		// Inefficient regular expression usage
		{
			ID:          "OPT006",
			Name:        "inefficient-regex",
			Description: "Inefficient regular expression usage",
			Detector:    detectInefficientRegex,
		},
		// Inefficient error handling
		{
			ID:          "OPT007",
			Name:        "inefficient-error-handling",
			Description: "Inefficient error handling",
			Detector:    detectInefficientErrorHandling,
		},
		// Inefficient JSON marshaling
		{
			ID:          "OPT008",
			Name:        "inefficient-json",
			Description: "Inefficient JSON marshaling/unmarshaling",
			Detector:    detectInefficientJSON,
		},
	}
}

// detectInefficientStringConcat detects inefficient string concatenation in loops
func detectInefficientStringConcat(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for string concatenation in loops
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		rangeStmt, ok := node.(*ast.RangeStmt)
		if !ok {
			return nil
		}
		
		// Check for string concatenation in range loop body
		return checkStringConcatInBody(fset, rangeStmt.Body)
	}
	
	// Check for string concatenation in for loop body
	return checkStringConcatInBody(fset, forStmt.Body)
}

// checkStringConcatInBody checks for string concatenation in a block statement
func checkStringConcatInBody(fset *token.FileSet, body *ast.BlockStmt) *models.Optimization {
	if body == nil {
		return nil
	}
	
	// Look for string concatenation using += operator
	for _, stmt := range body.List {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok || assignStmt.Tok != token.ADD_ASSIGN {
			continue
		}
		
		// Check if left side is a string variable
		// This is a simplified check and would need type information for accuracy
		if len(assignStmt.Lhs) > 0 {
			pos := fset.Position(assignStmt.Pos())
			return &models.Optimization{
				File:        pos.Filename,
				Line:        pos.Line,
				Description: "Inefficient string concatenation in loop",
				Benefit:     "Reduced memory allocations and improved performance",
				Example:     "// Instead of:\nvar result string\nfor _, s := range strings {\n    result += s\n}\n\n// Use strings.Builder:\nvar builder strings.Builder\nfor _, s := range strings {\n    builder.WriteString(s)\n}\nresult := builder.String()",
			}
		}
	}
	
	return nil
}

// detectUnnecessaryAllocation detects unnecessary memory allocations
func detectUnnecessaryAllocation(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for unnecessary use of new() or make() for small structs
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	
	// Check if it's a call to new()
	if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "new" {
		pos := fset.Position(callExpr.Pos())
		return &models.Optimization{
			File:        pos.Filename,
			Line:        pos.Line,
			Description: "Unnecessary use of new() for small struct",
			Benefit:     "Reduced heap allocations and improved performance",
			Example:     "// Instead of:\nuser := new(User)\nuser.Name = \"John\"\n\n// Use struct literal:\nuser := User{Name: \"John\"}",
		}
	}
	
	return nil
}

// detectSuboptimalSliceCapacity detects suboptimal slice capacity
func detectSuboptimalSliceCapacity(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for slice creation followed by append in a loop
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		rangeStmt, ok := node.(*ast.RangeStmt)
		if !ok {
			return nil
		}
		
		// Check for slice append in range loop body
		return checkSliceAppendInBody(fset, rangeStmt.Body)
	}
	
	// Check for slice append in for loop body
	return checkSliceAppendInBody(fset, forStmt.Body)
}

// checkSliceAppendInBody checks for slice append in a block statement
func checkSliceAppendInBody(fset *token.FileSet, body *ast.BlockStmt) *models.Optimization {
	if body == nil {
		return nil
	}
	
	// Look for append calls
	for _, stmt := range body.List {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		
		for _, rhs := range assignStmt.Rhs {
			callExpr, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}
			
			if ident, ok := callExpr.Fun.(*ast.Ident); ok && ident.Name == "append" {
				pos := fset.Position(callExpr.Pos())
				return &models.Optimization{
					File:        pos.Filename,
					Line:        pos.Line,
					Description: "Slice being repeatedly appended to in a loop without pre-allocation",
					Benefit:     "Reduced memory allocations and improved performance",
					Example:     "// Instead of:\nvar items []Item\nfor i := 0; i < n; i++ {\n    items = append(items, Item{})\n}\n\n// Pre-allocate the slice:\nitems := make([]Item, 0, n)\nfor i := 0; i < n; i++ {\n    items = append(items, Item{})\n}",
				}
			}
		}
	}
	
	return nil
}

// detectInefficientMapInit detects inefficient map initialization
func detectInefficientMapInit(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for map creation without capacity hint followed by multiple insertions
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		rangeStmt, ok := node.(*ast.RangeStmt)
		if !ok {
			return nil
		}
		
		// Check for map assignments in range loop body
		return checkMapAssignInBody(fset, rangeStmt.Body)
	}
	
	// Check for map assignments in for loop body
	return checkMapAssignInBody(fset, forStmt.Body)
}

// checkMapAssignInBody checks for map assignments in a block statement
func checkMapAssignInBody(fset *token.FileSet, body *ast.BlockStmt) *models.Optimization {
	if body == nil {
		return nil
	}
	
	// Look for map assignments
	for _, stmt := range body.List {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		
		for _, lhs := range assignStmt.Lhs {
			indexExpr, ok := lhs.(*ast.IndexExpr)
			if !ok {
				continue
			}
			
			// This is a simplified check and would need type information for accuracy
			pos := fset.Position(indexExpr.Pos())
			return &models.Optimization{
				File:        pos.Filename,
				Line:        pos.Line,
				Description: "Map being populated in a loop without capacity hint",
				Benefit:     "Reduced memory allocations and improved performance",
				Example:     "// Instead of:\nm := make(map[string]int)\nfor i := 0; i < n; i++ {\n    m[fmt.Sprintf(\"key%d\", i)] = i\n}\n\n// Provide capacity hint:\nm := make(map[string]int, n)\nfor i := 0; i < n; i++ {\n    m[fmt.Sprintf(\"key%d\", i)] = i\n}",
			}
		}
	}
	
	return nil
}

// detectRedundantTypeConversion detects redundant type conversions
func detectRedundantTypeConversion(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for redundant type conversions
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	
	// Check if it's a type conversion
	if typeIdent, ok := callExpr.Fun.(*ast.Ident); ok {
		// Check if argument is of the same type (simplified check)
		if len(callExpr.Args) == 1 {
			if argIdent, ok := callExpr.Args[0].(*ast.Ident); ok {
				// This is a simplified check and would need type information for accuracy
				if typeIdent.Name == "string" && argIdent.Name == "str" {
					pos := fset.Position(callExpr.Pos())
					return &models.Optimization{
						File:        pos.Filename,
						Line:        pos.Line,
						Description: "Potentially redundant type conversion",
						Benefit:     "Cleaner code and potentially improved performance",
						Example:     "// Instead of:\nresult := string(str)\n\n// If str is already a string, simply use:\nresult := str",
					}
				}
			}
		}
	}
	
	return nil
}

// detectInefficientRegex detects inefficient regular expression usage
func detectInefficientRegex(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for regexp.Compile or regexp.MustCompile in loops
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		rangeStmt, ok := node.(*ast.RangeStmt)
		if !ok {
			return nil
		}
		
		// Check for regex compilation in range loop body
		return checkRegexCompileInBody(fset, rangeStmt.Body)
	}
	
	// Check for regex compilation in for loop body
	return checkRegexCompileInBody(fset, forStmt.Body)
}

// checkRegexCompileInBody checks for regex compilation in a block statement
func checkRegexCompileInBody(fset *token.FileSet, body *ast.BlockStmt) *models.Optimization {
	if body == nil {
		return nil
	}
	
	// Look for regexp.Compile or regexp.MustCompile calls
	for _, stmt := range body.List {
		exprStmt, ok := stmt.(*ast.ExprStmt)
		if !ok {
			continue
		}
		
		callExpr, ok := exprStmt.X.(*ast.CallExpr)
		if !ok {
			continue
		}
		
		selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			continue
		}
		
		if xIdent, ok := selectorExpr.X.(*ast.Ident); ok && xIdent.Name == "regexp" {
			if selectorExpr.Sel.Name == "Compile" || selectorExpr.Sel.Name == "MustCompile" {
				pos := fset.Position(callExpr.Pos())
				return &models.Optimization{
					File:        pos.Filename,
					Line:        pos.Line,
					Description: "Regular expression compiled inside a loop",
					Benefit:     "Significantly improved performance by avoiding repeated regex compilation",
					Example:     "// Instead of:\nfor _, s := range strings {\n    re := regexp.MustCompile(`pattern`)\n    matches := re.FindAllString(s, -1)\n}\n\n// Compile the regex once, outside the loop:\nre := regexp.MustCompile(`pattern`)\nfor _, s := range strings {\n    matches := re.FindAllString(s, -1)\n}",
				}
			}
		}
	}
	
	return nil
}

// detectInefficientErrorHandling detects inefficient error handling
func detectInefficientErrorHandling(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for fmt.Errorf with string concatenation
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}
	
	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}
	
	if xIdent, ok := selectorExpr.X.(*ast.Ident); ok && xIdent.Name == "fmt" && selectorExpr.Sel.Name == "Errorf" {
		if len(callExpr.Args) > 0 {
			// Check if the format string contains %w
			if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok && basicLit.Kind == token.STRING {
				if !strings.Contains(basicLit.Value, "%w") && len(callExpr.Args) > 1 {
					// Check if any argument is an error
					for i := 1; i < len(callExpr.Args); i++ {
						if ident, ok := callExpr.Args[i].(*ast.Ident); ok && ident.Name == "err" {
							pos := fset.Position(callExpr.Pos())
							return &models.Optimization{
								File:        pos.Filename,
								Line:        pos.Line,
								Description: "Error wrapping without using %w verb",
								Benefit:     "Proper error wrapping allows for error unwrapping and inspection",
								Example:     "// Instead of:\nreturn fmt.Errorf(\"failed to process: \" + err.Error())\n// Or:\nreturn fmt.Errorf(\"failed to process: %v\", err)\n\n// Use %w for proper error wrapping:\nreturn fmt.Errorf(\"failed to process: %w\", err)",
							}
						}
					}
				}
			}
		}
	}
	
	return nil
}

// detectInefficientJSON detects inefficient JSON marshaling/unmarshaling
func detectInefficientJSON(fset *token.FileSet, node ast.Node) *models.Optimization {
	// Look for json.Marshal or json.Unmarshal in loops
	forStmt, ok := node.(*ast.ForStmt)
	if !ok {
		rangeStmt, ok := node.(*ast.RangeStmt)
		if !ok {
			return nil
		}
		
		// Check for JSON operations in range loop body
		return checkJSONInBody(fset, rangeStmt.Body)
	}
	
	// Check for JSON operations in for loop body
	return checkJSONInBody(fset, forStmt.Body)
}

// checkJSONInBody checks for JSON operations in a block statement
func checkJSONInBody(fset *token.FileSet, body *ast.BlockStmt) *models.Optimization {
	if body == nil {
		return nil
	}
	
	// Look for json.Marshal or json.Unmarshal calls
	for _, stmt := range body.List {
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}
		
		for _, rhs := range assignStmt.Rhs {
			callExpr, ok := rhs.(*ast.CallExpr)
			if !ok {
				continue
			}
			
			selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			
			if xIdent, ok := selectorExpr.X.(*ast.Ident); ok && xIdent.Name == "json" {
				if selectorExpr.Sel.Name == "Marshal" || selectorExpr.Sel.Name == "Unmarshal" {
					pos := fset.Position(callExpr.Pos())
					return &models.Optimization{
						File:        pos.Filename,
						Line:        pos.Line,
						Description: "JSON marshaling/unmarshaling inside a loop",
						Benefit:     "Improved performance by reducing repeated encoding/decoding operations",
						Example:     "// For multiple JSON operations on the same structure, consider:\n// 1. Using a JSON encoder/decoder with io.Pipe for streaming\n// 2. Processing data in batches\n// 3. Using a more efficient encoding like gob or protobuf for internal operations",
					}
				}
			}
		}
	}
	
	return nil
}
