package security

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
)

// GosecResult represents a security issue found by gosec
type GosecResult struct {
	Severity   string `json:"severity"`
	Confidence string `json:"confidence"`
	CWE        struct {
		ID          string `json:"id"`
		URL         string `json:"url"`
		Description string `json:"description"`
	} `json:"cwe"`
	Rule    string `json:"rule_id"`
	Details string `json:"details"`
	File    string `json:"file"`
	Code    string `json:"code"`
	Line    string `json:"line"`
	Column  string `json:"column"`
}

// GosecResults represents the results of a gosec scan
type GosecResults struct {
	Issues []GosecResult `json:"Issues"`
	Stats  struct {
		Files int `json:"files"`
		Lines int `json:"lines"`
		Nosec int `json:"nosec"`
		Found int `json:"found"`
	} `json:"Stats"`
}

// GosecScanner is responsible for scanning code for security vulnerabilities using gosec
type GosecScanner struct {
	config *config.Config
}

// NewGosecScanner creates a new gosec scanner
func NewGosecScanner(cfg *config.Config) *GosecScanner {
	return &GosecScanner{
		config: cfg,
	}
}

// Scan scans a repository for security vulnerabilities
func (s *GosecScanner) Scan(repoPath string) ([]*models.Issue, error) {
	// Create a temporary file to store gosec results
	tmpFile, err := os.CreateTemp("", "gosec-results-*.json")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Build gosec command
	cmd := exec.Command("gosec", "-fmt=json", "-out="+tmpFile.Name(), "-exclude-dir=vendor", "./...")
	cmd.Dir = repoPath

	// Run gosec
	if s.config.Verbose {
		fmt.Println("Running gosec...")
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		// gosec returns non-zero exit code when issues are found, so we need to check if the output file exists
		if _, statErr := os.Stat(tmpFile.Name()); statErr != nil {
			return nil, fmt.Errorf("gosec failed: %w\nOutput: %s", err, string(output))
		}
	}

	// Read gosec results
	resultsData, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read gosec results: %w", err)
	}

	// Parse gosec results
	var gosecResults GosecResults
	if err := json.Unmarshal(resultsData, &gosecResults); err != nil {
		return nil, fmt.Errorf("failed to parse gosec results: %w", err)
	}

	// Convert gosec results to our model
	issues := make([]*models.Issue, 0, len(gosecResults.Issues))
	for _, result := range gosecResults.Issues {
		// Convert line and column to integers
		line, _ := strconv.Atoi(result.Line)
		column, _ := strconv.Atoi(result.Column)

		// Make file path relative to repository root
		relPath, err := filepath.Rel(repoPath, result.File)
		if err != nil {
			relPath = result.File
		}

		// Map gosec severity to our severity levels
		severity := mapGosecSeverity(result.Severity)

		// Create issue
		issue := &models.Issue{
			File:       relPath,
			Line:       line,
			Column:     column,
			Message:    result.Details,
			Category:   "security",
			Severity:   severity,
			Confidence: strings.ToLower(result.Confidence),
			Rule:       result.Rule,
			Code:       result.Code,
			Suggestion: getSuggestionForRule(result.Rule, result.CWE.Description),
		}

		issues = append(issues, issue)
	}

	if s.config.Verbose {
		fmt.Printf("Found %d security issues\n", len(issues))
	}

	return issues, nil
}

// mapGosecSeverity maps gosec severity to our severity levels
func mapGosecSeverity(severity string) string {
	switch severity {
	case "HIGH":
		return "critical"
	case "MEDIUM":
		return "high"
	case "LOW":
		return "medium"
	default:
		return "low"
	}
}

// getSuggestionForRule returns a suggestion for a specific gosec rule
func getSuggestionForRule(rule string, cweDescription string) string {
	suggestions := map[string]string{
		"G101": "Remove hardcoded credentials and use environment variables or a secure vault.",
		"G102": "Avoid binding to all interfaces unless absolutely necessary. Specify a specific IP address.",
		"G103": "Avoid using the unsafe package. If necessary, add proper bounds checking.",
		"G104": "Always check returned errors. Consider using the 'errcheck' tool.",
		"G106": "Avoid using ssh.InsecureIgnoreHostKey. Implement proper host key verification.",
		"G107": "Validate and sanitize URL input to prevent SSRF attacks.",
		"G108": "Disable profiling endpoints in production or protect them with authentication.",
		"G109": "Use appropriate integer types to prevent overflow.",
		"G110": "Implement size limits when handling compressed data to prevent DoS attacks.",
		"G111": "Validate file paths to prevent directory traversal attacks.",
		"G112": "Set appropriate timeouts to prevent slowloris attacks.",
		"G114": "Use http.Server with configured timeouts instead of http.Serve.",
		"G201": "Use parameterized queries or an ORM to prevent SQL injection.",
		"G202": "Use parameterized queries instead of string concatenation to prevent SQL injection.",
		"G203": "Use proper HTML templating to prevent XSS attacks.",
		"G204": "Validate and sanitize input before passing to exec functions.",
		"G301": "Use appropriate file permissions (e.g., 0600 for sensitive files).",
		"G302": "Use appropriate file permissions with chmod.",
		"G303": "Use crypto/rand for generating temporary file paths.",
		"G304": "Validate and sanitize file paths from user input.",
		"G305": "Validate archive contents before extraction to prevent zip slip attacks.",
		"G306": "Use appropriate file permissions when writing to a new file.",
		"G307": "Use appropriate file permissions when creating a file with os.Create.",
		"G401": "Use crypto/sha256 or crypto/sha512 instead of crypto/md5 or crypto/sha1.",
		"G402": "Use secure TLS configuration with modern cipher suites and protocols.",
		"G403": "Use RSA keys of at least 2048 bits.",
		"G404": "Use crypto/rand instead of math/rand for security-sensitive operations.",
		"G405": "Use AES instead of DES or RC4.",
		"G406": "Use SHA-256 or SHA-3 instead of MD4 or RIPEMD160.",
		"G407": "Use a secure random generator for IV/nonce values.",
		"G501": "Replace crypto/md5 imports with crypto/sha256 or crypto/sha512.",
		"G502": "Replace crypto/des imports with crypto/aes.",
		"G503": "Replace crypto/rc4 imports with crypto/aes.",
		"G504": "Replace net/http/cgi imports with net/http.",
		"G505": "Replace crypto/sha1 imports with crypto/sha256 or crypto/sha512.",
		"G506": "Replace golang.org/x/crypto/md4 imports with crypto/sha256 or crypto/sha512.",
		"G507": "Replace golang.org/x/crypto/ripemd160 imports with crypto/sha256 or crypto/sha512.",
		"G601": "Fix implicit memory aliasing in range statements.",
		"G602": "Add bounds checking for slice access.",
	}

	if suggestion, ok := suggestions[rule]; ok {
		return suggestion
	}

	// If no specific suggestion is available, return a generic one based on CWE description
	if cweDescription != "" {
		return fmt.Sprintf("Address the security issue: %s", cweDescription)
	}

	return "Review the code for security vulnerabilities and apply appropriate fixes."
}
