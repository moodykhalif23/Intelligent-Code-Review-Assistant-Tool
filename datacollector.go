package ml

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
)

// DataCollector is responsible for collecting data for machine learning
type DataCollector struct {
	config    *config.Config
	dataPath  string
	issueData map[string][]models.LearningData // Map of rule ID to learning data
}

// NewDataCollector creates a new data collector
func NewDataCollector(cfg *config.Config) *DataCollector {
	dataPath := "data"
	if cfg.ModelPath != "" {
		dataPath = cfg.ModelPath
	}

	return &DataCollector{
		config:    cfg,
		dataPath:  dataPath,
		issueData: make(map[string][]models.LearningData),
	}
}

// Initialize initializes the data collector
func (c *DataCollector) Initialize() error {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(c.dataPath, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %w", err)
	}

	// Load existing data
	return c.loadData()
}

// loadData loads existing learning data
func (c *DataCollector) loadData() error {
	dataFile := filepath.Join(c.dataPath, "learning_data.json")

	// Check if file exists
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		// File doesn't exist, initialize empty data
		c.issueData = make(map[string][]models.LearningData)
		return nil
	}

	// Read file
	data, err := os.ReadFile(dataFile)
	if err != nil {
		return fmt.Errorf("failed to read learning data: %w", err)
	}

	// Parse JSON
	if err := json.Unmarshal(data, &c.issueData); err != nil {
		return fmt.Errorf("failed to parse learning data: %w", err)
	}

	return nil
}

// saveData saves learning data to disk
func (c *DataCollector) saveData() error {
	dataFile := filepath.Join(c.dataPath, "learning_data.json")

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(c.issueData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal learning data: %w", err)
	}

	// Write to file
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write learning data: %w", err)
	}

	return nil
}

// RecordIssue records an issue for learning
func (c *DataCollector) RecordIssue(issue *models.Issue, repository string) error {
	if !c.config.EnableLearning {
		return nil
	}

	// Create learning data
	learningData := models.LearningData{
		Issue:      issue,
		Accepted:   false, // Will be updated when feedback is received
		Repository: repository,
		Timestamp:  time.Now(),
	}

	// Add to data
	c.issueData[issue.Rule] = append(c.issueData[issue.Rule], learningData)

	// Save data
	return c.saveData()
}

// RecordFeedback records feedback for an issue
func (c *DataCollector) RecordFeedback(issueID string, accepted bool) error {
	if !c.config.EnableLearning {
		return nil
	}

	// Find issue in data
	for ruleID, issues := range c.issueData {
		for i, data := range issues {
			if data.Issue.File+":"+fmt.Sprint(data.Issue.Line)+":"+data.Issue.Message == issueID {
				// Update acceptance
				c.issueData[ruleID][i].Accepted = accepted
				
				// Save data
				return c.saveData()
			}
		}
	}

	return fmt.Errorf("issue not found: %s", issueID)
}

// GetAcceptanceRate returns the acceptance rate for a rule
func (c *DataCollector) GetAcceptanceRate(ruleID string) float64 {
	issues, ok := c.issueData[ruleID]
	if !ok || len(issues) == 0 {
		return 0.5 // Default to 50% if no data
	}

	// Count accepted issues
	accepted := 0
	for _, data := range issues {
		if data.Accepted {
			accepted++
		}
	}

	return float64(accepted) / float64(len(issues))
}

// GetRuleConfidence returns the confidence level for a rule
func (c *DataCollector) GetRuleConfidence(ruleID string) string {
	rate := c.GetAcceptanceRate(ruleID)
	
	if rate >= 0.8 {
		return "high"
	} else if rate >= 0.5 {
		return "medium"
	} else {
		return "low"
	}
}

// GetTopRules returns the top N rules by acceptance rate
func (c *DataCollector) GetTopRules(n int) []string {
	// Create slice of rule IDs and acceptance rates
	type ruleRate struct {
		ID   string
		Rate float64
	}
	
	var rates []ruleRate
	for ruleID := range c.issueData {
		rates = append(rates, ruleRate{
			ID:   ruleID,
			Rate: c.GetAcceptanceRate(ruleID),
		})
	}
	
	// Sort by acceptance rate (descending)
	for i := 0; i < len(rates); i++ {
		for j := i + 1; j < len(rates); j++ {
			if rates[j].Rate > rates[i].Rate {
				rates[i], rates[j] = rates[j], rates[i]
			}
		}
	}
	
	// Get top N
	var result []string
	for i := 0; i < n && i < len(rates); i++ {
		result = append(result, rates[i].ID)
	}
	
	return result
}
