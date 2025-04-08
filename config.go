package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// General settings
	Verbose bool `json:"verbose"`
	
	// Analysis settings
	IncludeTests      bool     `json:"include_tests"`
	ExcludeDirs       []string `json:"exclude_dirs"`
	ExcludeFiles      []string `json:"exclude_files"`
	MaxFileSize       int64    `json:"max_file_size"`
	
	// Analyzer settings
	EnabledAnalyzers  []string `json:"enabled_analyzers"`
	DisabledAnalyzers []string `json:"disabled_analyzers"`
	
	// Security settings
	SecuritySeverity  string   `json:"security_severity"`
	
	// Pattern detection settings
	PatternSeverity   string   `json:"pattern_severity"`
	
	// Machine learning settings
	EnableLearning    bool     `json:"enable_learning"`
	ModelPath         string   `json:"model_path"`
	
	// Custom rules
	CustomRulesPath   string   `json:"custom_rules_path"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Verbose:           false,
		IncludeTests:      true,
		ExcludeDirs:       []string{".git", "vendor", "node_modules"},
		ExcludeFiles:      []string{},
		MaxFileSize:       1024 * 1024, // 1MB
		EnabledAnalyzers:  []string{"all"},
		DisabledAnalyzers: []string{},
		SecuritySeverity:  "high",
		PatternSeverity:   "medium",
		EnableLearning:    true,
		ModelPath:         "",
		CustomRulesPath:   "",
	}
}

// LoadConfig loads configuration from a file
func LoadConfig(configPath string) (*Config, error) {
	config := DefaultConfig()
	
	// If no config file specified, return default config
	if configPath == "" {
		return config, nil
	}
	
	// Resolve absolute path
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path: %w", err)
	}
	
	// Read config file
	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	
	// Parse JSON
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}
	
	return config, nil
}

// SaveConfig saves configuration to a file
func SaveConfig(config *Config, configPath string) error {
	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	return nil
}
