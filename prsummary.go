package prsummary

import (
	"bytes"
	"fmt"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
)

// PRSummaryGenerator generates summaries of pull requests
type PRSummaryGenerator struct {
	config *config.Config
}

// NewPRSummaryGenerator creates a new PR summary generator
func NewPRSummaryGenerator(cfg *config.Config) *PRSummaryGenerator {
	return &PRSummaryGenerator{
		config: cfg,
	}
}

// GenerateSummary generates a summary of changes between two Git references
func (g *PRSummaryGenerator) GenerateSummary(repoPath, baseRef, headRef string) (*models.PRSummary, error) {
	// Ensure we're in a Git repository
	if !isGitRepository(repoPath) {
		return nil, fmt.Errorf("not a git repository: %s", repoPath)
	}

	// Get diff stats
	stats, err := g.getDiffStats(repoPath, baseRef, headRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get diff stats: %w", err)
	}

	// Get changed files
	changedFiles, err := g.getChangedFiles(repoPath, baseRef, headRef)
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	// Analyze affected areas
	affectedAreas := g.analyzeAffectedAreas(changedFiles)

	// Identify key changes
	keyChanges := g.identifyKeyChanges(repoPath, baseRef, headRef, changedFiles)

	// Create summary
	summary := &models.PRSummary{
		FilesChanged:  stats.filesChanged,
		Additions:     stats.additions,
		Deletions:     stats.deletions,
		KeyChanges:    keyChanges,
		AffectedAreas: affectedAreas,
	}

	return summary, nil
}

// diffStats represents statistics about a diff
type diffStats struct {
	filesChanged int
	additions    int
	deletions    int
}

// getDiffStats gets statistics about changes between two Git references
func (g *PRSummaryGenerator) getDiffStats(repoPath, baseRef, headRef string) (*diffStats, error) {
	cmd := exec.Command("git", "diff", "--shortstat", fmt.Sprintf("%s...%s", baseRef, headRef))
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output
	stats := &diffStats{}
	outputStr := string(output)

	// Example output: " 2 files changed, 3 insertions(+), 1 deletion(-)"
	if match := regexp.MustCompile(`(\d+) files? changed`).FindStringSubmatch(outputStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &stats.filesChanged)
	}
	if match := regexp.MustCompile(`(\d+) insertions?\(\+\)`).FindStringSubmatch(outputStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &stats.additions)
	}
	if match := regexp.MustCompile(`(\d+) deletions?\(-\)`).FindStringSubmatch(outputStr); len(match) > 1 {
		fmt.Sscanf(match[1], "%d", &stats.deletions)
	}

	return stats, nil
}

// getChangedFiles gets a list of files changed between two Git references
func (g *PRSummaryGenerator) getChangedFiles(repoPath, baseRef, headRef string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", fmt.Sprintf("%s...%s", baseRef, headRef))
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output by newline
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// Filter out empty strings
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}

	return result, nil
}

// analyzeAffectedAreas analyzes which areas of the codebase are affected by changes
func (g *PRSummaryGenerator) analyzeAffectedAreas(changedFiles []string) []string {
	// Map to track unique areas
	areas := make(map[string]bool)

	for _, file := range changedFiles {
		// Get directory path
		dir := filepath.Dir(file)
		if dir == "." {
			continue
		}

		// Split path into components
		components := strings.Split(dir, "/")
		
		// Add first component as an area
		if len(components) > 0 {
			areas[components[0]] = true
		}
		
		// Add first two components as a more specific area
		if len(components) > 1 {
			areas[components[0]+"/"+components[1]] = true
		}
	}

	// Convert map to slice
	var result []string
	for area := range areas {
		result = append(result, area)
	}

	// Sort areas
	sort.Strings(result)

	return result
}

// identifyKeyChanges identifies key changes between two Git references
func (g *PRSummaryGenerator) identifyKeyChanges(repoPath, baseRef, headRef string, changedFiles []string) []string {
	var keyChanges []string

	// Check for new files
	newFiles, _ := g.getNewFiles(repoPath, baseRef, headRef)
	if len(newFiles) > 0 {
		if len(newFiles) <= 3 {
			keyChanges = append(keyChanges, fmt.Sprintf("Added new files: %s", strings.Join(newFiles, ", ")))
		} else {
			keyChanges = append(keyChanges, fmt.Sprintf("Added %d new files", len(newFiles)))
		}
	}

	// Check for deleted files
	deletedFiles, _ := g.getDeletedFiles(repoPath, baseRef, headRef)
	if len(deletedFiles) > 0 {
		if len(deletedFiles) <= 3 {
			keyChanges = append(keyChanges, fmt.Sprintf("Deleted files: %s", strings.Join(deletedFiles, ", ")))
		} else {
			keyChanges = append(keyChanges, fmt.Sprintf("Deleted %d files", len(deletedFiles)))
		}
	}

	// Check for changes to important files
	importantFiles := g.findChangesToImportantFiles(changedFiles)
	for _, file := range importantFiles {
		keyChanges = append(keyChanges, fmt.Sprintf("Modified important file: %s", file))
	}

	// Check for large changes
	largeChanges, _ := g.findLargeChanges(repoPath, baseRef, headRef)
	for _, change := range largeChanges {
		keyChanges = append(keyChanges, fmt.Sprintf("Large change to %s (%d lines)", change.file, change.lines))
	}

	// Check for changes to interfaces
	interfaceChanges, _ := g.findInterfaceChanges(repoPath, baseRef, headRef)
	for _, change := range interfaceChanges {
		keyChanges = append(keyChanges, fmt.Sprintf("Modified interface: %s", change))
	}

	return keyChanges
}

// getNewFiles gets a list of new files added between two Git references
func (g *PRSummaryGenerator) getNewFiles(repoPath, baseRef, headRef string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=A", fmt.Sprintf("%s...%s", baseRef, headRef))
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output by newline
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// Filter out empty strings
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}

	return result, nil
}

// getDeletedFiles gets a list of files deleted between two Git references
func (g *PRSummaryGenerator) getDeletedFiles(repoPath, baseRef, headRef string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=D", fmt.Sprintf("%s...%s", baseRef, headRef))
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Split output by newline
	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	// Filter out empty strings
	var result []string
	for _, file := range files {
		if file != "" {
			result = append(result, file)
		}
	}

	return result, nil
}

// findChangesToImportantFiles finds changes to important files
func (g *PRSummaryGenerator) findChangesToImportantFiles(changedFiles []string) []string {
	var importantFiles []string

	// Define patterns for important files
	importantPatterns := []string{
		"go.mod",
		"go.sum",
		"Dockerfile",
		"docker-compose.yml",
		"Makefile",
		"README.md",
		"LICENSE",
		".github/workflows/",
		"main.go",
		"config.go",
	}

	// Check each changed file against patterns
	for _, file := range changedFiles {
		for _, pattern := range importantPatterns {
			if strings.HasPrefix(file, pattern) || file == pattern {
				importantFiles = append(importantFiles, file)
				break
			}
		}
	}

	return importantFiles
}

// largeChange represents a large change to a file
type largeChange struct {
	file  string
	lines int
}

// findLargeChanges finds files with large changes
func (g *PRSummaryGenerator) findLargeChanges(repoPath, baseRef, headRef string) ([]largeChange, error) {
	cmd := exec.Command("git", "diff", "--stat", fmt.Sprintf("%s...%s", baseRef, headRef))
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Parse output
	var largeChanges []largeChange
	lines := strings.Split(string(output), "\n")
	
	for _, line := range lines {
		if match := regexp.MustCompile(`(.+?)\s+\|\s+(\d+)\s+`).FindStringSubmatch(line); len(match) > 2 {
			file := strings.TrimSpace(match[1])
			var lineCount int
			fmt.Sscanf(match[2], "%d", &lineCount)
			
			// Consider changes with more than 50 lines as large
			if lineCount > 50 {
				largeChanges = append(largeChanges, largeChange{
					file:  file,
					lines: lineCount,
				})
			}
		}
	}

	return largeChanges, nil
}

// findInterfaceChanges finds changes to interfaces
func (g *PRSummaryGenerator) findInterfaceChanges(repoPath, baseRef, headRef string) ([]string, error) {
	cmd := exec.Command("git", "diff", fmt.Sprintf("%s...%s", baseRef, headRef))
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Look for interface definitions in the diff
	var interfaceChanges []string
	lines := strings.Split(string(output), "\n")
	
	for i, line := range lines {
		if strings.Contains(line, "type") && strings.Contains(line, "interface") {
			// Extract interface name
			if match := regexp.MustCompile(`type\s+(\w+)\s+interface`).FindStringSubmatch(line); len(match) > 1 {
				interfaceChanges = append(interfaceChanges, match[1])
			}
		}
	}

	return interfaceChanges, nil
}

// isGitRepository checks if a directory is a Git repository
func isGitRepository(path string) bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Dir = path
	
	err := cmd.Run()
	return err == nil
}

// FormatSummary formats a PR summary as a string
func (g *PRSummaryGenerator) FormatSummary(summary *models.PRSummary) string {
	var buf bytes.Buffer

	buf.WriteString("# Pull Request Summary\n\n")
	
	buf.WriteString("## Overview\n")
	buf.WriteString(fmt.Sprintf("- **Files Changed:** %d\n", summary.FilesChanged))
	buf.WriteString(fmt.Sprintf("- **Lines Added:** %d\n", summary.Additions))
	buf.WriteString(fmt.Sprintf("- **Lines Deleted:** %d\n", summary.Deletions))
	buf.WriteString("\n")
	
	buf.WriteString("## Key Changes\n")
	if len(summary.KeyChanges) > 0 {
		for _, change := range summary.KeyChanges {
			buf.WriteString(fmt.Sprintf("- %s\n", change))
		}
	} else {
		buf.WriteString("- No significant changes detected\n")
	}
	buf.WriteString("\n")
	
	buf.WriteString("## Affected Areas\n")
	if len(summary.AffectedAreas) > 0 {
		for _, area := range summary.AffectedAreas {
			buf.WriteString(fmt.Sprintf("- %s\n", area))
		}
	} else {
		buf.WriteString("- No specific areas affected\n")
	}
	buf.WriteString("\n")
	
	if len(summary.PotentialIssues) > 0 {
		buf.WriteString("## Potential Issues\n")
		for _, issue := range summary.PotentialIssues {
			buf.WriteString(fmt.Sprintf("- **%s:** %s in `%s:%d`\n", issue.Severity, issue.Message, issue.File, issue.Line))
		}
	}

	return buf.String()
}
