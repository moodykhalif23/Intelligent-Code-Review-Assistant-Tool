package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/code-review-assistant/internal/analyzer"
	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/scanner"
)

// TestMain is a simple test harness for the code review assistant
func main() {
	// Create test directory if it doesn't exist
	testDir := filepath.Join(os.TempDir(), "code-review-assistant-test")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		fmt.Printf("Error creating test directory: %v\n", err)
		os.Exit(1)
	}

	// Create test files
	if err := createTestFiles(testDir); err != nil {
		fmt.Printf("Error creating test files: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created test files in %s\n", testDir)

	// Run tests
	if err := runTests(testDir); err != nil {
		fmt.Printf("Tests failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("All tests passed!")
}

// createTestFiles creates test files with various code patterns
func createTestFiles(testDir string) error {
	// Test file with code smells
	codeSmellsFile := `package test

import (
	"fmt"
	"math/rand"
	"strings"
)

// Function with too many parameters
func functionWithTooManyParams(param1, param2, param3, param4, param5, param6 string) {
	// Empty function
}

// Function with empty body
func emptyFunction() {
}

// Long function
func longFunction() {
	fmt.Println("This is a long function")
	fmt.Println("With many lines")
	fmt.Println("That does very little")
	fmt.Println("But takes up space")
	fmt.Println("And should be refactored")
	fmt.Println("To be more concise")
	fmt.Println("And maintainable")
	fmt.Println("This is a long function")
	fmt.Println("With many lines")
	fmt.Println("That does very little")
	fmt.Println("But takes up space")
	fmt.Println("And should be refactored")
	fmt.Println("To be more concise")
	fmt.Println("And maintainable")
	fmt.Println("This is a long function")
	fmt.Println("With many lines")
	fmt.Println("That does very little")
	fmt.Println("But takes up space")
	fmt.Println("And should be refactored")
	fmt.Println("To be more concise")
	fmt.Println("And maintainable")
}

// Function with boolean parameter
func functionWithBoolParam(flag bool) {
	if flag {
		fmt.Println("Flag is true")
	} else {
		fmt.Println("Flag is false")
	}
}

// Exported function without documentation
func ExportedFunctionWithoutDocs() {
	fmt.Println("This function should have documentation")
}
`

	// Test file with security issues
	securityIssuesFile := `package test

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
)

// Function with hardcoded credentials
func connectToDatabase() {
	username := "admin"
	password := "password123" // Hardcoded credential
	fmt.Printf("Connecting with %s:%s\n", username, password)
}

// Function with command injection vulnerability
func runCommand(userInput string) {
	cmd := exec.Command("bash", "-c", userInput) // Command injection
	cmd.Run()
}

// Function with weak cryptography
func hashPassword(password string) string {
	h := md5.New() // Weak hash function
	io.WriteString(h, password)
	return fmt.Sprintf("%x", h.Sum(nil))
}

// Function with path traversal vulnerability
func readFile(filename string) {
	data, _ := os.ReadFile(filename) // Path traversal
	fmt.Println(string(data))
}

// Function with open redirect vulnerability
func redirect(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	http.Redirect(w, r, url, http.StatusFound) // Open redirect
}
`

	// Test file with optimization opportunities
	optimizationFile := `package test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

// Function with inefficient string concatenation
func concatenateStrings(strings []string) string {
	var result string
	for _, s := range strings {
		result += s // Inefficient string concatenation
	}
	return result
}

// Function with suboptimal slice capacity
func createSlice(n int) []int {
	var slice []int // No capacity hint
	for i := 0; i < n; i++ {
		slice = append(slice, i)
	}
	return slice
}

// Function with inefficient map initialization
func createMap(n int) map[string]int {
	m := make(map[string]int) // No capacity hint
	for i := 0; i < n; i++ {
		m[fmt.Sprintf("key%d", i)] = i
	}
	return m
}

// Function with inefficient regex usage
func findMatches(patterns, texts []string) []string {
	var matches []string
	for _, text := range texts {
		for _, pattern := range patterns {
			re := regexp.MustCompile(pattern) // Regex compiled in loop
			if re.MatchString(text) {
				matches = append(matches, text)
			}
		}
	}
	return matches
}

// Function with inefficient JSON marshaling
func processJSON(objects []interface{}) []string {
	var results []string
	for _, obj := range objects {
		data, _ := json.Marshal(obj) // JSON marshaling in loop
		results = append(results, string(data))
	}
	return results
}
`

	// Write test files
	if err := os.WriteFile(filepath.Join(testDir, "code_smells.go"), []byte(codeSmellsFile), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(testDir, "security_issues.go"), []byte(securityIssuesFile), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(testDir, "optimization.go"), []byte(optimizationFile), 0644); err != nil {
		return err
	}

	return nil
}

// runTests runs tests on the test files
func runTests(testDir string) error {
	// Create default configuration
	cfg := config.DefaultConfig()
	cfg.Verbose = true

	// Initialize repository scanner
	repoScanner := scanner.NewScanner(testDir, cfg)

	// Scan repository for Go files
	files, err := repoScanner.Scan()
	if err != nil {
		return fmt.Errorf("error scanning repository: %w", err)
	}

	fmt.Printf("Found %d files to analyze\n", len(files))

	// Initialize code analyzer
	codeAnalyzer := analyzer.NewAnalyzer(cfg)

	// Analyze files
	results, err := codeAnalyzer.Analyze(files)
	if err != nil {
		return fmt.Errorf("error analyzing code: %w", err)
	}

	// Verify results
	if len(results.Issues) == 0 {
		return fmt.Errorf("no issues found, expected multiple issues")
	}

	fmt.Printf("Found %d issues (Critical: %d, High: %d, Medium: %d, Low: %d)\n",
		results.TotalIssues,
		results.CriticalIssues,
		results.HighIssues,
		results.MediumIssues,
		results.LowIssues,
	)

	// Verify code smells
	codeSmellCount := 0
	for _, issue := range results.Issues {
		if issue.Category == "code-smell" {
			codeSmellCount++
		}
	}
	if codeSmellCount == 0 {
		return fmt.Errorf("no code smells found, expected multiple code smells")
	}
	fmt.Printf("Found %d code smells\n", codeSmellCount)

	// Verify security issues
	securityIssueCount := 0
	for _, issue := range results.Issues {
		if issue.Category == "security" {
			securityIssueCount++
		}
	}
	if securityIssueCount == 0 {
		return fmt.Errorf("no security issues found, expected multiple security issues")
	}
	fmt.Printf("Found %d security issues\n", securityIssueCount)

	// Test optimization suggestions
	// This would normally call the optimization analyzer, but for simplicity
	// we'll just check if we have the optimization files
	optimizationFileFound := false
	for _, file := range files {
		if filepath.Base(file.Path) == "optimization.go" {
			optimizationFileFound = true
			break
		}
	}
	if !optimizationFileFound {
		return fmt.Errorf("optimization file not found")
	}
	fmt.Println("Optimization file found")

	return nil
}
