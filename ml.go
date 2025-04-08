package cmd

import (
	"fmt"

	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/ml"
	"github.com/user/code-review-assistant/internal/models"
)

// RecordIssue records an issue for machine learning
func RecordIssue(issue *models.Issue, repository string, cfg *config.Config) error {
	// Create learning engine
	engine, err := ml.NewLearningEngine(cfg)
	if err != nil {
		return fmt.Errorf("failed to create learning engine: %w", err)
	}

	// Record issue
	return engine.RecordIssue(issue, repository)
}

// RecordFeedback records feedback for an issue
func RecordFeedback(issueID string, accepted bool, cfg *config.Config) error {
	// Create learning engine
	engine, err := ml.NewLearningEngine(cfg)
	if err != nil {
		return fmt.Errorf("failed to create learning engine: %w", err)
	}

	// Record feedback
	return engine.RecordFeedback(issueID, accepted)
}

// ApplyLearning applies machine learning to improve analysis results
func ApplyLearning(issues []*models.Issue, repository string, cfg *config.Config) ([]*models.Issue, []string, error) {
	if !cfg.EnableLearning {
		return issues, nil, nil
	}

	// Create learning engine
	engine, err := ml.NewLearningEngine(cfg)
	if err != nil {
		return issues, nil, fmt.Errorf("failed to create learning engine: %w", err)
	}

	// Adjust issue confidence
	for _, issue := range issues {
		engine.AdjustIssueConfidence(issue)
	}

	// Filter issues
	filteredIssues := engine.FilterIssues(issues)

	// Sort issues
	sortedIssues := engine.SortIssues(filteredIssues)

	// Get project insights
	insights := engine.AnalyzeProjectPatterns(repository, issues)

	// Add custom rule suggestions
	customRules := engine.SuggestCustomRules()
	if len(customRules) > 0 {
		insights = append(insights, "Suggested custom rules:")
		insights = append(insights, customRules...)
	}

	return sortedIssues, insights, nil
}
