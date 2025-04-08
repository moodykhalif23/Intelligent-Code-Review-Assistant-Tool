# Intelligent Code Review Assistant - Configuration Guide

## Overview

This document provides information on how to configure the Intelligent Code Review Assistant to suit your specific needs. The tool can be configured using command-line flags or a configuration file.

## Command-Line Flags

The following command-line flags are available:

### Common Flags

- `-repo`: Path to the repository to analyze (default: current directory)
- `-config`: Path to configuration file
- `-verbose`: Enable verbose output
- `-format`: Output format (text, json, html)
- `-version`: Show version information

### Analysis Flags

- `-include-tests`: Include test files in analysis (default: true)
- `-exclude-dirs`: Comma-separated list of directories to exclude (default: .git,vendor,node_modules)
- `-exclude-files`: Comma-separated list of files to exclude

### PR Summary Flags

- `-base`: Base reference for PR summary (default: main)
- `-head`: Head reference for PR summary (default: HEAD)

### Command Flags

- `-analyze`: Run code analysis
- `-summary`: Generate PR summary
- `-optimize`: Suggest optimizations
- `-learn`: Enable machine learning
- `-feedback`: Provide feedback for an issue
- `-issue-id`: Issue ID for feedback
- `-accepted`: Whether the issue was accepted

## Configuration File

The configuration file is in JSON format and can include the following settings:

```json
{
  "verbose": false,
  "include_tests": true,
  "exclude_dirs": [".git", "vendor", "node_modules"],
  "exclude_files": [],
  "max_file_size": 1048576,
  "enabled_analyzers": ["all"],
  "disabled_analyzers": [],
  "security_severity": "high",
  "pattern_severity": "medium",
  "enable_learning": true,
  "model_path": "",
  "custom_rules_path": ""
}
```

### Configuration Options

- `verbose`: Enable verbose output
- `include_tests`: Include test files in analysis
- `exclude_dirs`: List of directories to exclude
- `exclude_files`: List of files to exclude
- `max_file_size`: Maximum file size to analyze (in bytes)
- `enabled_analyzers`: List of analyzers to enable (use "all" for all analyzers)
- `disabled_analyzers`: List of analyzers to disable
- `security_severity`: Minimum severity for security issues (critical, high, medium, low)
- `pattern_severity`: Minimum severity for pattern issues (critical, high, medium, low)
- `enable_learning`: Enable machine learning
- `model_path`: Path to store machine learning model data
- `custom_rules_path`: Path to custom rules

## Examples

### Basic Analysis

```bash
code-review-assistant -analyze -repo /path/to/repo
```

### Generate PR Summary

```bash
code-review-assistant -summary -base main -head feature-branch
```

### Suggest Optimizations

```bash
code-review-assistant -optimize -repo /path/to/repo
```

### Provide Feedback

```bash
code-review-assistant -feedback -issue-id "file.go:10:Error not handled" -accepted
```

### Using a Configuration File

```bash
code-review-assistant -config config.json -analyze
```

## Machine Learning

The tool includes a machine learning component that improves over time based on feedback. To enable this feature:

1. Use the `-learn` flag when running analysis
2. Provide feedback on issues using the `-feedback` command

The machine learning system will:
- Adjust issue confidence based on historical feedback
- Filter out suggestions with low acceptance rates
- Sort issues by a combination of severity and acceptance probability
- Suggest custom rules based on successful patterns
- Analyze project-specific patterns to provide tailored insights
