package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// TestIntegration runs integration tests for the code review assistant
func TestIntegration(t *testing.T) {
	// Create test repository
	testDir := filepath.Join("testdata", "integration")
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("Error creating test directory: %v", err)
	}

	// Initialize Git repository
	if err := initGitRepo(testDir); err != nil {
		t.Fatalf("Error initializing Git repository: %v", err)
	}

	// Create test files
	if err := createTestFiles(testDir); err != nil {
		t.Fatalf("Error creating test files: %v", err)
	}

	// Commit files
	if err := commitFiles(testDir); err != nil {
		t.Fatalf("Error committing files: %v", err)
	}

	// Run integration tests
	t.Run("FullAnalysis", func(t *testing.T) { testFullAnalysis(t, testDir) })
	t.Run("PRSummary", func(t *testing.T) { testPRSummary(t, testDir) })
	t.Run("Optimization", func(t *testing.T) { testOptimization(t, testDir) })
}

// initGitRepo initializes a Git repository
func initGitRepo(dir string) error {
	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize Git repository: %w", err)
	}

	// Configure Git user
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure Git user name: %w", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure Git user email: %w", err)
	}

	return nil
}

// createTestFiles creates test files for integration testing
func createTestFiles(dir string) error {
	// Create main.go with some issues
	mainFile := `package main

import (
	"fmt"
	"math/rand"
	"strings"
)

func main() {
	// Inefficient string concatenation
	var result string
	for i := 0; i < 100; i++ {
		result += fmt.Sprintf("%d", i)
	}

	// Security issue: using math/rand
	randomValue := rand.Intn(100)
	fmt.Println("Random value:", randomValue)

	// Call function with too many parameters
	complexFunction("param1", "param2", "param3", "param4", "param5", "param6")

	fmt.Println(result)
}

// Function with too many parameters
func complexFunction(param1, param2, param3, param4, param5, param6 string) {
	// Empty function
}

// Exported function without documentation
func ExportedFunction() {
	fmt.Println("This function should have documentation")
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(mainFile), 0644); err != nil {
		return err
	}

	// Create utils.go with more issues
	utilsFile := `package main

import (
	"fmt"
	"os/exec"
)

// Function with security vulnerability
func runCommand(userInput string) {
	cmd := exec.Command("bash", "-c", userInput) // Command injection
	cmd.Run()
}

// Function with optimization opportunity
func createSlice(n int) []int {
	var slice []int // No capacity hint
	for i := 0; i < n; i++ {
		slice = append(slice, i)
	}
	return slice
}

// Function with hardcoded credentials
func connectToDatabase() {
	username := "admin"
	password := "password123" // Hardcoded credential
	fmt.Printf("Connecting with %s:%s\n", username, password)
}
`
	if err := os.WriteFile(filepath.Join(dir, "utils.go"), []byte(utilsFile), 0644); err != nil {
		return err
	}

	return nil
}

// commitFiles commits the test files to the Git repository
func commitFiles(dir string) error {
	// Add files
	cmd := exec.Command("git", "add", ".")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to add files: %w", err)
	}

	// Commit files
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to commit files: %w", err)
	}

	return nil
}

// testFullAnalysis tests the full analysis functionality
func testFullAnalysis(t *testing.T, dir string) {
	t.Log("Testing full analysis functionality...")

	// TODO: Run the code review assistant on the test repository
	// This would require building and running the executable

	t.Log("Full analysis tests completed")
}

// testPRSummary tests the PR summary generation functionality
func testPRSummary(t *testing.T, dir string) {
	t.Log("Testing PR summary generation functionality...")

	// Create a new branch
	cmd := exec.Command("git", "checkout", "-b", "feature-branch")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to create feature branch: %v", err)
	}

	// Modify a file
	modifiedFile := `package main

import (
	"fmt"
	"crypto/rand"
	"strings"
	"encoding/hex"
)

func main() {
	// Fixed: Use strings.Builder instead of string concatenation
	var builder strings.Builder
	for i := 0; i < 100; i++ {
		builder.WriteString(fmt.Sprintf("%d", i))
	}
	result := builder.String()

	// Fixed: Use crypto/rand instead of math/rand
	randomBytes := make([]byte, 8)
	rand.Read(randomBytes)
	randomValue := hex.EncodeToString(randomBytes)
	fmt.Println("Random value:", randomValue)

	// Call improved function with fewer parameters
	improvedFunction("param1", "param2")

	fmt.Println(result)
}

// Improved function with fewer parameters
func improvedFunction(param1, param2 string) {
	fmt.Println("Using improved function with fewer parameters")
}

// Exported function with documentation
// ExportedFunction demonstrates proper documentation for exported functions
func ExportedFunction() {
	fmt.Println("This function now has documentation")
}
`
	if err := os.WriteFile(filepath.Join(dir, "main.go"), []byte(modifiedFile), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	// Commit the changes
	cmd = exec.Command("git", "add", "main.go")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add modified file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Improve code quality")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
	}

	// TODO: Run the PR summary generator on the test repository
	// This would require building and running the executable

	t.Log("PR summary tests completed")
}

// testOptimization tests the optimization suggestion functionality
func testOptimization(t *testing.T, dir string) {
	t.Log("Testing optimization suggestion functionality...")

	// TODO: Run the optimization suggestion system on the test repository
	// This would require building and running the executable

	t.Log("Optimization tests completed")
}

func main() {
	// Run the tests
	testing.Main(func(pat, str string) (bool, error) { return true, nil }, []testing.InternalTest{
		{Name: "TestIntegration", F: TestIntegration},
	}, nil, nil)
}
