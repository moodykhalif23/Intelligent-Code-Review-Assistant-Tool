package models

import (
	"time"
)

// File represents a source code file to be analyzed
type File struct {
	Path     string    // Absolute path to the file
	RelPath  string    // Path relative to the repository root
	Size     int64     // File size in bytes
	ModTime  time.Time // Last modification time
	IsVendor bool      // Whether the file is in a vendor directory
}

// Repository represents a code repository
type Repository struct {
	Path  string // Path to the repository
	VCS   string // Version control system (git, svn, etc.)
	IsGit bool   // Whether it's a Git repository
}

// Issue represents a code issue found during analysis
type Issue struct {
	File       string // File path where the issue was found
	Line       int    // Line number where the issue was found
	Column     int    // Column number where the issue was found
	Message    string // Description of the issue
	Category   string // Category of the issue (e.g., "code-smell", "security", "performance")
	Severity   string // Severity of the issue (e.g., "critical", "high", "medium", "low")
	Confidence string // Confidence level of the issue (e.g., "high", "medium", "low")
	Suggestion string // Suggested fix for the issue
	Code       string // The problematic code snippet
	Rule       string // The rule that triggered the issue
}

// Function represents a function or method in the code
type Function struct {
	Name       string // Function name
	File       string // File path where the function is defined
	StartLine  int    // Start line of the function
	EndLine    int    // End line of the function
	Complexity int    // Cyclomatic complexity of the function
	Parameters int    // Number of parameters
	Returns    int    // Number of return values
}

// Package represents a Go package
type Package struct {
	Name       string     // Package name
	ImportPath string     // Import path
	Files      []*File    // Files in the package
	Functions  []*Function // Functions in the package
}

// PRSummary represents a summary of a pull request
type PRSummary struct {
	FilesChanged    int      // Number of files changed
	Additions       int      // Number of lines added
	Deletions       int      // Number of lines deleted
	KeyChanges      []string // Key changes in the PR
	AffectedAreas   []string // Areas of the codebase affected
	PotentialIssues []*Issue // Potential issues in the PR
}

// Optimization represents a suggested optimization
type Optimization struct {
	File        string // File path where the optimization can be applied
	Line        int    // Line number where the optimization can be applied
	Description string // Description of the optimization
	Benefit     string // Expected benefit of the optimization
	Example     string // Example code with the optimization applied
}

// LearningData represents data used for machine learning
type LearningData struct {
	Issue      *Issue // The issue that was detected
	Accepted   bool   // Whether the suggestion was accepted
	Repository string // Repository where the issue was found
	Timestamp  time.Time // When the data was collected
}
