package security

import (
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/user/code-review-assistant/internal/models"
)

// CustomSecurityRule represents a custom security rule
type CustomSecurityRule struct {
	ID          string
	Name        string
	Description string
	Severity    string
	Detector    func(fset *token.FileSet, node ast.Node) *models.Issue
}

// GetCustomSecurityRules returns a list of custom security rules
func GetCustomSecurityRules() []*CustomSecurityRule {
	return []*CustomSecurityRule{
		// Hardcoded secrets in string literals
		{
			ID:          "CS001",
			Name:        "hardcoded-secret",
			Description: "Hardcoded secret or credential",
			Severity:    "critical",
			Detector:    detectHardcodedSecrets,
		},
		// Insecure random number generation
		{
			ID:          "CS002",
			Name:        "insecure-random",
			Description: "Insecure random number generation",
			Severity:    "high",
			Detector:    detectInsecureRandom,
		},
		// Missing content type in HTTP responses
		{
			ID:          "CS003",
			Name:        "missing-content-type",
			Description: "Missing Content-Type header in HTTP response",
			Severity:    "medium",
			Detector:    detectMissingContentType,
		},
		// Insecure cookie settings
		{
			ID:          "CS004",
			Name:        "insecure-cookie",
			Description: "Insecure cookie settings",
			Severity:    "high",
			Detector:    detectInsecureCookie,
		},
		// Weak cryptographic key size
		{
			ID:          "CS005",
			Name:        "weak-crypto-key",
			Description: "Weak cryptographic key size",
			Severity:    "high",
			Detector:    detectWeakCryptoKey,
		},
		// Unvalidated redirect
		{
			ID:          "CS006",
			Name:        "unvalidated-redirect",
			Description: "Unvalidated redirect",
			Severity:    "medium",
			Detector:    detectUnvalidatedRedirect,
		},
		// Logging sensitive information
		{
			ID:          "CS007",
			Name:        "sensitive-log",
			Description: "Logging sensitive information",
			Severity:    "medium",
			Detector:    detectSensitiveLogging,
		},
	}
}

// detectHardcodedSecrets detects hardcoded secrets in string literals
func detectHardcodedSecrets(fset *token.FileSet, node ast.Node) *models.Issue {
	basicLit, ok := node.(*ast.BasicLit)
	if !ok || basicLit.Kind != token.STRING {
		return nil
	}

	// Remove quotes from string literal
	value := strings.Trim(basicLit.Value, `"'`)

	// Skip short strings
	if len(value) < 8 {
		return nil
	}

	// Patterns for potential secrets
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`(?i)password\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)passwd\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)pwd\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)secret\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)api[_-]?key\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)access[_-]?token\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)auth[_-]?token\s*=\s*['"](.+?)['"]`),
		regexp.MustCompile(`(?i)credentials\s*=\s*['"](.+?)['"]`),
	}

	for _, pattern := range patterns {
		if pattern.MatchString(value) {
			pos := fset.Position(basicLit.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Hardcoded secret or credential detected",
				Category:   "security",
				Severity:   "critical",
				Confidence: "medium",
				Suggestion: "Store secrets in environment variables or a secure vault, not in source code",
				Rule:       "CS001",
			}
		}
	}

	return nil
}

// detectInsecureRandom detects insecure random number generation
func detectInsecureRandom(fset *token.FileSet, node ast.Node) *models.Issue {
	// Look for imports of math/rand
	importSpec, ok := node.(*ast.ImportSpec)
	if ok {
		path := strings.Trim(importSpec.Path.Value, `"`)
		if path == "math/rand" {
			pos := fset.Position(importSpec.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Import of math/rand package which is not cryptographically secure",
				Category:   "security",
				Severity:   "high",
				Confidence: "high",
				Suggestion: "Use crypto/rand for security-sensitive operations",
				Rule:       "CS002",
			}
		}
		return nil
	}

	// Look for calls to math/rand functions
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil
	}

	if ident, ok := selectorExpr.X.(*ast.Ident); ok {
		if ident.Name == "rand" && (selectorExpr.Sel.Name == "Int" || 
			selectorExpr.Sel.Name == "Intn" || 
			selectorExpr.Sel.Name == "Float64") {
			pos := fset.Position(callExpr.Pos())
			return &models.Issue{
				File:       pos.Filename,
				Line:       pos.Line,
				Column:     pos.Column,
				Message:    "Use of math/rand." + selectorExpr.Sel.Name + " which is not cryptographically secure",
				Category:   "security",
				Severity:   "high",
				Confidence: "medium",
				Suggestion: "Use crypto/rand for security-sensitive operations",
				Rule:       "CS002",
			}
		}
	}

	return nil
}

// detectMissingContentType detects missing Content-Type header in HTTP responses
func detectMissingContentType(fset *token.FileSet, node ast.Node) *models.Issue {
	// Look for http.ResponseWriter.Write calls without setting Content-Type
	callExpr, ok := node.(*ast.CallExpr)
	if !ok {
		return nil
	}

	selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr)
	if !ok || selectorExpr.Sel.Name != "Write" {
		return nil
	}

	// This is a simplified check and would need more context analysis
	// to determine if it's an HTTP response writer and if Content-Type is set
	pos := fset.Position(callExpr.Pos())
	return &models.Issue{
		File:       pos.Filename,
		Line:       pos.Line,
		Column:     pos.Column,
		Message:    "Potential missing Content-Type header in HTTP response",
		Category:   "security",
		Severity:   "medium",
		Confidence: "low",
		Suggestion: "Set Content-Type header before writing to the response",
		Rule:       "CS003",
	}
}

// detectInsecureCookie detects insecure cookie settings
func detectInsecureCookie(fset *token.FileSet, node ast.Node) *models.Issue {
	// Look for http.Cookie creation without Secure and HttpOnly flags
	compositeLit, ok := node.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	// Check if it's an http.Cookie
	if typeExpr, ok := compositeLit.Type.(*ast.SelectorExpr); ok {
		if ident, ok := typeExpr.X.(*ast.Ident); ok && ident.Name == "http" && typeExpr.Sel.Name == "Cookie" {
			// Check if Secure and HttpOnly are set to true
			secureSet := false
			httpOnlySet := false

			for _, elt := range compositeLit.Elts {
				if kv, ok := elt.(*ast.KeyValueExpr); ok {
					if key, ok := kv.Key.(*ast.Ident); ok {
						if key.Name == "Secure" {
							if lit, ok := kv.Value.(*ast.Ident); ok && lit.Name == "true" {
								secureSet = true
							}
						} else if key.Name == "HttpOnly" {
							if lit, ok := kv.Value.(*ast.Ident); ok && lit.Name == "true" {
								httpOnlySet = true
							}
						}
					}
				}
			}

			if !secureSet || !httpOnlySet {
				pos := fset.Position(compositeLit.Pos())
				message := "Cookie created without "
				if !secureSet && !httpOnlySet {
					message += "Secure and HttpOnly flags"
				} else if !secureSet {
					message += "Secure flag"
				} else {
					message += "HttpOnly flag"
				}

				return &models.Issue{
					File:       pos.Filename,
					Line:       pos.Line,
					Column:     pos.Column,
					Message:    message,
					Category:   "security",
					Severity:   "high",
					Confidence: "high",
					Suggestion: "Set both Secure and HttpOnly flags to true for cookies",
					Rule:       "CS004",
				}
			}
		}
	}

	return nil
}

// detectWeakCryptoKey detects weak cryptographic key sizes
func detectWeakCryptoKey(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectUnvalidatedRedirect detects unvalidated redirects
func detectUnvalidatedRedirect(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}

// detectSensitiveLogging detects logging of sensitive information
func detectSensitiveLogging(fset *token.FileSet, node ast.Node) *models.Issue {
	// Implementation will be added
	return nil
}
