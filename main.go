package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/code-review-assistant/internal/analyzer"
	"github.com/user/code-review-assistant/internal/cmd"
	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/prsummary"
	"github.com/user/code-review-assistant/internal/scanner"
)

const (
	version = "1.0.0"
)

func main() {
	// Define command-line flags
	var (
		// Common flags
		repoPath      = flag.String("repo", ".", "Path to the repository to analyze")
		configFile    = flag.String("config", "", "Path to configuration file")
		verbose       = flag.Bool("verbose", false, "Enable verbose output")
		outputFormat  = flag.String("format", "text", "Output format (text, json, html)")
		showVersion   = flag.Bool("version", false, "Show version information")
		
		// Analysis flags
		includeTests  = flag.Bool("include-tests", true, "Include test files in analysis")
		excludeDirs   = flag.String("exclude-dirs", ".git,vendor,node_modules", "Comma-separated list of directories to exclude")
		excludeFiles  = flag.String("exclude-files", "", "Comma-separated list of files to exclude")
		
		// PR summary flags
		baseRef       = flag.String("base", "main", "Base reference for PR summary")
		headRef       = flag.String("head", "HEAD", "Head reference for PR summary")
		
		// Command flags
		analyzeCmd    = flag.Bool("analyze", false, "Run code analysis")
		summaryCmd    = flag.Bool("summary", false, "Generate PR summary")
		optimizeCmd   = flag.Bool("optimize", false, "Suggest optimizations")
		learnCmd      = flag.Bool("learn", false, "Enable machine learning")
		feedbackCmd   = flag.Bool("feedback", false, "Provide feedback for an issue")
		issueID       = flag.String("issue-id", "", "Issue ID for feedback")
		accepted      = flag.Bool("accepted", false, "Whether the issue was accepted")
	)
	
	// Custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Intelligent Code Review Assistant v%s\n\n", version)
		fmt.Fprintf(os.Stderr, "Usage: %s [options] [command]\n\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "Commands:\n")
		fmt.Fprintf(os.Stderr, "  -analyze              Run code analysis\n")
		fmt.Fprintf(os.Stderr, "  -summary              Generate PR summary\n")
		fmt.Fprintf(os.Stderr, "  -optimize             Suggest optimizations\n")
		fmt.Fprintf(os.Stderr, "  -feedback             Provide feedback for an issue\n")
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s -analyze -repo /path/to/repo\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -summary -base main -head feature-branch\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -optimize -repo /path/to/repo\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s -feedback -issue-id \"file.go:10:Error not handled\" -accepted\n", os.Args[0])
	}
	
	flag.Parse()
	
	// Show version and exit
	if *showVersion {
		fmt.Printf("Intelligent Code Review Assistant v%s\n", version)
		os.Exit(0)
	}
	
	// Resolve absolute path for repository
	absPath, err := filepath.Abs(*repoPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving repository path: %v\n", err)
		os.Exit(1)
	}
	
	// Load configuration
	cfg, err := loadConfig(*configFile, *verbose, *includeTests, *excludeDirs, *excludeFiles, *learnCmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading configuration: %v\n", err)
		os.Exit(1)
	}
	
	// Check if at least one command is specified
	if !*analyzeCmd && !*summaryCmd && !*optimizeCmd && !*feedbackCmd {
		// Default to analyze if no command is specified
		*analyzeCmd = true
	}
	
	// Handle feedback command
	if *feedbackCmd {
		if *issueID == "" {
			fmt.Fprintf(os.Stderr, "Error: issue-id is required for feedback command\n")
			os.Exit(1)
		}
		
		if err := cmd.RecordFeedback(*issueID, *accepted, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error recording feedback: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Printf("Feedback recorded for issue: %s\n", *issueID)
		os.Exit(0)
	}
	
	// Handle PR summary command
	if *summaryCmd {
		if err := generatePRSummary(absPath, *baseRef, *headRef, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating PR summary: %v\n", err)
			os.Exit(1)
		}
		
		// Exit if only summary command is specified
		if !*analyzeCmd && !*optimizeCmd {
			os.Exit(0)
		}
	}
	
	// Initialize repository scanner
	repoScanner := scanner.NewScanner(absPath, cfg)
	
	// Scan repository for Go files
	files, err := repoScanner.Scan()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning repository: %v\n", err)
		os.Exit(1)
	}
	
	if cfg.Verbose {
		fmt.Printf("Found %d files to analyze\n", len(files))
	}
	
	// Handle analyze command
	if *analyzeCmd {
		if err := analyzeCode(files, absPath, *outputFormat, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing code: %v\n", err)
			os.Exit(1)
		}
	}
	
	// Handle optimize command
	if *optimizeCmd {
		if err := suggestOptimizations(files, cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error suggesting optimizations: %v\n", err)
			os.Exit(1)
		}
	}
}

// loadConfig loads configuration from a file or creates a default configuration
func loadConfig(configFile string, verbose, includeTests bool, excludeDirs, excludeFiles string, enableLearning bool) (*config.Config, error) {
	var cfg *config.Config
	var err error
	
	if configFile != "" {
		// Load from file
		cfg, err = config.LoadConfig(configFile)
		if err != nil {
			return nil, err
		}
	} else {
		// Create default config
		cfg = config.DefaultConfig()
	}
	
	// Override with command-line flags
	cfg.Verbose = verbose
	cfg.IncludeTests = includeTests
	cfg.EnableLearning = enableLearning
	
	// Parse exclude dirs
	if excludeDirs != "" {
		cfg.ExcludeDirs = filepath.SplitList(excludeDirs)
	}
	
	// Parse exclude files
	if excludeFiles != "" {
		cfg.ExcludeFiles = filepath.SplitList(excludeFiles)
	}
	
	return cfg, nil
}

// analyzeCode analyzes code and prints results
func analyzeCode(files []*models.File, repoPath, outputFormat string, cfg *config.Config) error {
	// Initialize code analyzer
	codeAnalyzer := analyzer.NewAnalyzer(cfg)
	
	// Analyze files
	results, err := codeAnalyzer.Analyze(files)
	if err != nil {
		return fmt.Errorf("failed to analyze code: %w", err)
	}
	
	// Apply machine learning if enabled
	if cfg.EnableLearning {
		sortedIssues, insights, err := cmd.ApplyLearning(results.Issues, repoPath, cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to apply machine learning: %v\n", err)
		} else {
			results.Issues = sortedIssues
			
			// Print insights
			if len(insights) > 0 {
				fmt.Println("\nProject Insights:")
				for _, insight := range insights {
					fmt.Printf("- %s\n", insight)
				}
				fmt.Println()
			}
		}
	}
	
	// Record issues for learning
	if cfg.EnableLearning {
		for _, issue := range results.Issues {
			if err := cmd.RecordIssue(issue, repoPath, cfg); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to record issue for learning: %v\n", err)
			}
		}
	}
	
	// Output results based on format
	switch outputFormat {
	case "text":
		printTextResults(results)
	case "json":
		printJSONResults(results)
	case "html":
		printHTMLResults(results)
	default:
		return fmt.Errorf("unsupported output format: %s", outputFormat)
	}
	
	return nil
}

// generatePRSummary generates a PR summary
func generatePRSummary(repoPath, baseRef, headRef string, cfg *config.Config) error {
	// Create PR summary generator
	generator := prsummary.NewPRSummaryGenerator(cfg)
	
	// Generate summary
	summary, err := generator.GenerateSummary(repoPath, baseRef, headRef)
	if err != nil {
		return fmt.Errorf("failed to generate PR summary: %w", err)
	}
	
	// Format summary
	formattedSummary := generator.FormatSummary(summary)
	
	// Print summary
	fmt.Println(formattedSummary)
	
	return nil
}

// suggestOptimizations suggests code optimizations
func suggestOptimizations(files []*models.File, cfg *config.Config) error {
	return cmd.AnalyzeOptimizations(files, cfg)
}

// printTextResults prints analysis results in text format
func printTextResults(results *analyzer.Results) {
	fmt.Println("Code Review Results:")
	fmt.Println("====================")
	
	if len(results.Issues) == 0 {
		fmt.Println("No issues found!")
		return
	}
	
	for _, issue := range results.Issues {
		fmt.Printf("[%s] %s: %s\n", issue.Severity, issue.Category, issue.Message)
		fmt.Printf("  File: %s:%d\n", issue.File, issue.Line)
		if issue.Suggestion != "" {
			fmt.Printf("  Suggestion: %s\n", issue.Suggestion)
		}
		fmt.Println()
	}
	
	fmt.Printf("Total issues: %d (Critical: %d, High: %d, Medium: %d, Low: %d)\n",
		results.TotalIssues,
		results.CriticalIssues,
		results.HighIssues,
		results.MediumIssues,
		results.LowIssues,
	)
}

// printJSONResults prints analysis results in JSON format
func printJSONResults(results *analyzer.Results) {
	// Placeholder for JSON output
	fmt.Println("JSON output not yet implemented")
}

// printHTMLResults prints analysis results in HTML format
func printHTMLResults(results *analyzer.Results) {
	// Placeholder for HTML output
	fmt.Println("HTML output not yet implemented")
}
