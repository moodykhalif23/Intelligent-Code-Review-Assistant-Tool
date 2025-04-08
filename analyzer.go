package analyzer

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"sync"

	"github.com/user/code-review-assistant/internal/analyzer/patterns"
	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
	"github.com/user/code-review-assistant/internal/security"
)

// Results represents the results of code analysis
type Results struct {
	Issues         []*models.Issue
	TotalIssues    int
	CriticalIssues int
	HighIssues     int
	MediumIssues   int
	LowIssues      int
}

// Analyzer is responsible for analyzing code and finding issues
type Analyzer struct {
	config         *config.Config
	fset           *token.FileSet
	patterns       []*patterns.Pattern
	antiPatterns   []*patterns.AntiPattern
	bestPractices  []*patterns.BestPractice
	securityRules  []*security.CustomSecurityRule
	securityScanner *security.GosecScanner
}

// NewAnalyzer creates a new code analyzer
func NewAnalyzer(cfg *config.Config) *Analyzer {
	return &Analyzer{
		config:         cfg,
		fset:           token.NewFileSet(),
		patterns:       patterns.GetGoPatterns(),
		antiPatterns:   patterns.GetGoAntiPatterns(),
		bestPractices:  patterns.GetGoBestPractices(),
		securityRules:  security.GetCustomSecurityRules(),
		securityScanner: security.NewGosecScanner(cfg),
	}
}

// Analyze analyzes a list of files and returns the results
func (a *Analyzer) Analyze(files []*models.File) (*Results, error) {
	results := &Results{
		Issues: make([]*models.Issue, 0),
	}

	// Use a wait group to process files concurrently
	var wg sync.WaitGroup
	var mutex sync.Mutex

	// Process each file
	for _, file := range files {
		wg.Add(1)
		go func(f *models.File) {
			defer wg.Done()

			// Analyze the file
			issues, err := a.analyzeFile(f)
			if err != nil {
				if a.config.Verbose {
					println("Error analyzing file", f.Path, ":", err.Error())
				}
				return
			}

			// Add issues to results
			mutex.Lock()
			results.Issues = append(results.Issues, issues...)
			mutex.Unlock()
		}(file)
	}

	// Wait for all files to be processed
	wg.Wait()

	// Run security scanner on the repository
	if len(files) > 0 {
		// Get repository path from the first file
		repoPath := files[0].Path
		for i := 0; i < len(repoPath); i++ {
			if repoPath[i:] == files[0].RelPath {
				repoPath = repoPath[:i]
				break
			}
		}

		// Run security scanner
		securityIssues, err := a.securityScanner.Scan(repoPath)
		if err != nil {
			if a.config.Verbose {
				println("Error running security scanner:", err.Error())
			}
		} else {
			// Add security issues to results
			mutex.Lock()
			results.Issues = append(results.Issues, securityIssues...)
			mutex.Unlock()
		}
	}

	// Count issues by severity
	for _, issue := range results.Issues {
		results.TotalIssues++
		switch issue.Severity {
		case "critical":
			results.CriticalIssues++
		case "high":
			results.HighIssues++
		case "medium":
			results.MediumIssues++
		case "low":
			results.LowIssues++
		}
	}

	return results, nil
}

// analyzeFile analyzes a single file and returns a list of issues
func (a *Analyzer) analyzeFile(file *models.File) ([]*models.Issue, error) {
	issues := make([]*models.Issue, 0)

	// Read file content
	content, err := ioutil.ReadFile(file.Path)
	if err != nil {
		return nil, err
	}

	// Parse the file
	astFile, err := parser.ParseFile(a.fset, file.Path, content, parser.AllErrors|parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Apply all pattern detectors
	ast.Inspect(astFile, func(node ast.Node) bool {
		if node == nil {
			return true
		}

		// Apply code smell patterns
		for _, p := range a.patterns {
			if issue := p.Detector(a.fset, node); issue != nil {
				// Set relative path for consistent reporting
				issue.File = file.RelPath
				issues = append(issues, issue)
			}
		}

		// Apply anti-patterns
		for _, ap := range a.antiPatterns {
			if issue := ap.Detector(a.fset, node); issue != nil {
				// Set relative path for consistent reporting
				issue.File = file.RelPath
				issues = append(issues, issue)
			}
		}

		// Apply best practices
		for _, bp := range a.bestPractices {
			if issue := bp.Detector(a.fset, node); issue != nil {
				// Set relative path for consistent reporting
				issue.File = file.RelPath
				issues = append(issues, issue)
			}
		}

		// Apply custom security rules
		for _, sr := range a.securityRules {
			if issue := sr.Detector(a.fset, node); issue != nil {
				// Set relative path for consistent reporting
				issue.File = file.RelPath
				issues = append(issues, issue)
			}
		}

		return true
	})

	return issues, nil
}
