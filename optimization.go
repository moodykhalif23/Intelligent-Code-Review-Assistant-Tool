package cmd

import (
	"fmt"
	"os"

	"github.com/user/code-review-assistant/internal/optimization"
	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
)

// AnalyzeOptimizations analyzes a repository for optimization opportunities
func AnalyzeOptimizations(files []*models.File, cfg *config.Config) error {
	// Create optimization analyzer
	analyzer := optimization.NewAnalyzer(cfg)

	// Analyze files
	optimizations, err := analyzer.Analyze(files)
	if err != nil {
		return fmt.Errorf("failed to analyze optimizations: %w", err)
	}

	// Format optimizations
	formattedOptimizations := analyzer.FormatOptimizations(optimizations)

	// Print optimizations
	fmt.Println(formattedOptimizations)

	return nil
}
