# Intelligent Code Review Assistant

A Go-based tool that analyzes code repositories, identifies patterns, and provides automated code reviews.

## Features

- **Code Smell Detection**: Identifies code smells and anti-patterns specific to Go
- **Security Vulnerability Scanning**: Detects security vulnerabilities using static analysis
- **PR Summary Generation**: Generates concise, meaningful summaries of pull requests
- **Optimization Suggestions**: Recommends performance optimizations based on code analysis
- **Machine Learning**: Learns from accepted/rejected code reviews to improve suggestions over time

## Installation

### Pre-built Binaries

Download the pre-built binary for your platform from the releases page:

- [Linux (amd64)](dist/code-review-assistant-linux-amd64-1.0.0.tar.gz)
- [macOS (Intel)](dist/code-review-assistant-darwin-amd64-1.0.0.tar.gz)
- [macOS (Apple Silicon)](dist/code-review-assistant-darwin-arm64-1.0.0.tar.gz)
- [Windows](dist/code-review-assistant-windows-amd64-1.0.0.zip)

Extract the archive and place the executable in your PATH.

### Building from Source

To build from source, you need Go 1.18 or later installed.

```bash
# Clone the repository
git clone https://github.com/user/code-review-assistant.git
cd code-review-assistant

# Build the executable
go build -o code-review-assistant ./cmd/code-review-assistant

# Or use the build script to build for all platforms
chmod +x build.sh
./build.sh
```

## Usage

### Basic Analysis

Analyze a repository for code issues:

```bash
code-review-assistant -analyze -repo /path/to/repo
```

### Generate PR Summary

Generate a summary of changes between two Git references:

```bash
code-review-assistant -summary -base main -head feature-branch
```

### Suggest Optimizations

Analyze a repository for optimization opportunities:

```bash
code-review-assistant -optimize -repo /path/to/repo
```

### Enable Machine Learning

Enable machine learning to improve suggestions over time:

```bash
code-review-assistant -analyze -repo /path/to/repo -learn
```

### Provide Feedback

Provide feedback on an issue to improve future suggestions:

```bash
code-review-assistant -feedback -issue-id "file.go:10:Error not handled" -accepted
```

## Configuration

The tool can be configured using command-line flags or a configuration file. See the [Configuration Guide](docs/configuration.md) for details.

## Documentation

- [Configuration Guide](docs/configuration.md)
- [Architecture Overview](docs/architecture.md)
- [Contributing Guide](docs/contributing.md)

## License

This project is licensed under the MIT License - see the LICENSE file for details.
