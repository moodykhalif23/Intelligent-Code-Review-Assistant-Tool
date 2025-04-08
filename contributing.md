# Contributing to Intelligent Code Review Assistant

We welcome contributions to the Intelligent Code Review Assistant! This document provides guidelines and instructions for contributing to the project.

## Development Setup

1. **Prerequisites**
   - Go 1.18 or later
   - Git

2. **Clone the Repository**
   ```bash
   git clone https://github.com/user/code-review-assistant.git
   cd code-review-assistant
   ```

3. **Install Dependencies**
   ```bash
   go mod download
   ```

4. **Build the Project**
   ```bash
   go build -o code-review-assistant ./cmd/code-review-assistant
   ```

## Project Structure

- `cmd/code-review-assistant/`: Main application entry point
- `internal/`: Internal packages
  - `analyzer/`: Code analysis engine
  - `scanner/`: Repository scanner
  - `security/`: Security vulnerability scanning
  - `prsummary/`: PR summary generation
  - `optimization/`: Optimization suggestion system
  - `ml/`: Machine learning components
  - `config/`: Configuration management
  - `models/`: Data models
  - `cmd/`: Command implementations
- `test/`: Test files and test data
- `docs/`: Documentation

## Testing

Run the tests with:

```bash
go test ./...
```

Run the integration tests with:

```bash
go test ./test/integration_test.go
```

## Adding New Features

1. **Create a New Branch**
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Implement Your Feature**
   - Follow the existing code style and patterns
   - Add appropriate tests
   - Update documentation as needed

3. **Run Tests**
   ```bash
   go test ./...
   ```

4. **Submit a Pull Request**
   - Provide a clear description of the changes
   - Reference any related issues

## Adding New Analyzers

To add a new analyzer:

1. Create a new file in the appropriate package (e.g., `internal/analyzer/patterns/newpattern.go`)
2. Implement the analyzer following the existing pattern structure
3. Add the analyzer to the appropriate list in the main analyzer

Example:

```go
// newpattern.go
package patterns

import (
    "go/ast"
    "go/token"
    
    "github.com/user/code-review-assistant/internal/models"
)

// GetNewPatterns returns a list of new patterns
func GetNewPatterns() []*Pattern {
    return []*Pattern{
        {
            ID:          "NEW001",
            Name:        "new-pattern",
            Description: "Description of the new pattern",
            Detector:    detectNewPattern,
        },
    }
}

// detectNewPattern detects the new pattern
func detectNewPattern(fset *token.FileSet, node ast.Node) *models.Issue {
    // Implementation
}
```

## Code Style

- Follow standard Go code style and conventions
- Use `gofmt` to format your code
- Add comments for exported functions and types

## License

By contributing to this project, you agree that your contributions will be licensed under the project's MIT License.
