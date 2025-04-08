package ml

import (
	"math"
	"sort"
	"strings"

	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
)

// LearningEngine is responsible for learning from feedback and improving suggestions
type LearningEngine struct {
	config        *config.Config
	dataCollector *DataCollector
}

// NewLearningEngine creates a new learning engine
func NewLearningEngine(cfg *config.Config) (*LearningEngine, error) {
	// Create data collector
	dataCollector := NewDataCollector(cfg)
	
	// Initialize data collector
	if err := dataCollector.Initialize(); err != nil {
		return nil, err
	}
	
	return &LearningEngine{
		config:        cfg,
		dataCollector: dataCollector,
	}, nil
}

// RecordIssue records an issue for learning
func (e *LearningEngine) RecordIssue(issue *models.Issue, repository string) error {
	return e.dataCollector.RecordIssue(issue, repository)
}

// RecordFeedback records feedback for an issue
func (e *LearningEngine) RecordFeedback(issueID string, accepted bool) error {
	return e.dataCollector.RecordFeedback(issueID, accepted)
}

// AdjustIssueConfidence adjusts the confidence of an issue based on learning data
func (e *LearningEngine) AdjustIssueConfidence(issue *models.Issue) {
	if !e.config.EnableLearning {
		return
	}
	
	// Get rule confidence
	confidence := e.dataCollector.GetRuleConfidence(issue.Rule)
	
	// Adjust issue confidence
	switch issue.Confidence {
	case "high":
		if confidence == "low" {
			issue.Confidence = "medium"
		}
	case "medium":
		if confidence == "high" {
			issue.Confidence = "high"
		} else if confidence == "low" {
			issue.Confidence = "low"
		}
	case "low":
		if confidence == "high" {
			issue.Confidence = "medium"
		}
	}
}

// FilterIssues filters issues based on learning data
func (e *LearningEngine) FilterIssues(issues []*models.Issue) []*models.Issue {
	if !e.config.EnableLearning {
		return issues
	}
	
	// Filter out issues with low acceptance rate
	var filtered []*models.Issue
	for _, issue := range issues {
		rate := e.dataCollector.GetAcceptanceRate(issue.Rule)
		
		// Keep issues with acceptance rate >= 0.3
		if rate >= 0.3 {
			filtered = append(filtered, issue)
		}
	}
	
	return filtered
}

// SortIssues sorts issues based on learning data
func (e *LearningEngine) SortIssues(issues []*models.Issue) []*models.Issue {
	if !e.config.EnableLearning || len(issues) <= 1 {
		return issues
	}
	
	// Create a copy of issues
	sorted := make([]*models.Issue, len(issues))
	copy(sorted, issues)
	
	// Sort by a combination of severity and acceptance rate
	sort.Slice(sorted, func(i, j int) bool {
		// Get severity scores (critical=4, high=3, medium=2, low=1)
		iSeverity := getSeverityScore(sorted[i].Severity)
		jSeverity := getSeverityScore(sorted[j].Severity)
		
		// Get acceptance rates
		iRate := e.dataCollector.GetAcceptanceRate(sorted[i].Rule)
		jRate := e.dataCollector.GetAcceptanceRate(sorted[j].Rule)
		
		// Calculate combined scores
		iScore := float64(iSeverity) * (0.5 + 0.5*iRate)
		jScore := float64(jSeverity) * (0.5 + 0.5*jRate)
		
		// Sort by score (descending)
		return iScore > jScore
	})
	
	return sorted
}

// getSeverityScore returns a numeric score for a severity level
func getSeverityScore(severity string) int {
	switch severity {
	case "critical":
		return 4
	case "high":
		return 3
	case "medium":
		return 2
	case "low":
		return 1
	default:
		return 0
	}
}

// SuggestCustomRules suggests custom rules based on learning data
func (e *LearningEngine) SuggestCustomRules() []string {
	if !e.config.EnableLearning {
		return nil
	}
	
	// Get top rules
	topRules := e.dataCollector.GetTopRules(5)
	
	// Create suggestions
	var suggestions []string
	for _, rule := range topRules {
		suggestions = append(suggestions, "Consider creating custom rules similar to "+rule)
	}
	
	return suggestions
}

// AnalyzeProjectPatterns analyzes patterns specific to a project
func (e *LearningEngine) AnalyzeProjectPatterns(repository string, issues []*models.Issue) []string {
	if !e.config.EnableLearning {
		return nil
	}
	
	// Count occurrences of each rule
	ruleCounts := make(map[string]int)
	for _, issue := range issues {
		ruleCounts[issue.Rule]++
	}
	
	// Find rules that occur frequently
	var frequentRules []string
	for rule, count := range ruleCounts {
		if count >= 3 {
			frequentRules = append(frequentRules, rule)
		}
	}
	
	// Create insights
	var insights []string
	if len(frequentRules) > 0 {
		insights = append(insights, "Common issues in this project:")
		for _, rule := range frequentRules {
			insights = append(insights, "- "+rule+": occurs "+string(rune('0'+ruleCounts[rule]/10))+string(rune('0'+ruleCounts[rule]%10))+" times")
		}
	}
	
	return insights
}

// PredictIssueAcceptance predicts whether an issue will be accepted
func (e *LearningEngine) PredictIssueAcceptance(issue *models.Issue) float64 {
	if !e.config.EnableLearning {
		return 0.5 // Default to 50% if learning is disabled
	}
	
	// Base prediction on rule acceptance rate
	baseRate := e.dataCollector.GetAcceptanceRate(issue.Rule)
	
	// Adjust based on severity
	severityFactor := 0.0
	switch issue.Severity {
	case "critical":
		severityFactor = 0.2
	case "high":
		severityFactor = 0.1
	case "medium":
		severityFactor = 0.0
	case "low":
		severityFactor = -0.1
	}
	
	// Adjust based on confidence
	confidenceFactor := 0.0
	switch issue.Confidence {
	case "high":
		confidenceFactor = 0.1
	case "medium":
		confidenceFactor = 0.0
	case "low":
		confidenceFactor = -0.1
	}
	
	// Calculate final prediction
	prediction := baseRate + severityFactor + confidenceFactor
	
	// Clamp to [0, 1]
	return math.Max(0, math.Min(1, prediction))
}

// GetSimilarIssues finds similar issues to a given issue
func (e *LearningEngine) GetSimilarIssues(issue *models.Issue) []*models.Issue {
	if !e.config.EnableLearning {
		return nil
	}
	
	var similarIssues []*models.Issue
	
	// Get all issues for the same rule
	for _, data := range e.dataCollector.issueData[issue.Rule] {
		// Skip the same issue
		if data.Issue.File == issue.File && data.Issue.Line == issue.Line {
			continue
		}
		
		// Add to similar issues
		similarIssues = append(similarIssues, data.Issue)
	}
	
	// Sort by similarity
	sort.Slice(similarIssues, func(i, j int) bool {
		// Calculate similarity scores
		iScore := calculateSimilarity(issue, similarIssues[i])
		jScore := calculateSimilarity(issue, similarIssues[j])
		
		// Sort by score (descending)
		return iScore > jScore
	})
	
	// Return top 5 similar issues
	if len(similarIssues) > 5 {
		return similarIssues[:5]
	}
	return similarIssues
}

// calculateSimilarity calculates a similarity score between two issues
func calculateSimilarity(a, b *models.Issue) float64 {
	score := 0.0
	
	// Same rule
	if a.Rule == b.Rule {
		score += 1.0
	}
	
	// Same severity
	if a.Severity == b.Severity {
		score += 0.5
	}
	
	// Same confidence
	if a.Confidence == b.Confidence {
		score += 0.3
	}
	
	// Similar message
	if strings.Contains(a.Message, b.Message) || strings.Contains(b.Message, a.Message) {
		score += 0.7
	}
	
	return score
}
