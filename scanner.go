package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/user/code-review-assistant/internal/config"
	"github.com/user/code-review-assistant/internal/models"
)

// Scanner is responsible for scanning repositories and finding files to analyze
type Scanner struct {
	rootPath string
	config   *config.Config
}

// NewScanner creates a new repository scanner
func NewScanner(rootPath string, cfg *config.Config) *Scanner {
	return &Scanner{
		rootPath: rootPath,
		config:   cfg,
	}
}

// Scan scans the repository and returns a list of files to analyze
func (s *Scanner) Scan() ([]*models.File, error) {
	// Check if root path exists
	info, err := os.Stat(s.rootPath)
	if err != nil {
		return nil, fmt.Errorf("failed to access repository path: %w", err)
	}
	
	if !info.IsDir() {
		return nil, fmt.Errorf("repository path is not a directory: %s", s.rootPath)
	}
	
	var files []*models.File
	
	// Walk through the directory tree
	err = filepath.Walk(s.rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			// Check if directory should be excluded
			for _, excludeDir := range s.config.ExcludeDirs {
				if info.Name() == excludeDir {
					if s.config.Verbose {
						fmt.Printf("Skipping excluded directory: %s\n", path)
					}
					return filepath.SkipDir
				}
			}
			return nil
		}
		
		// Skip files that are too large
		if info.Size() > s.config.MaxFileSize {
			if s.config.Verbose {
				fmt.Printf("Skipping file (too large): %s\n", path)
			}
			return nil
		}
		
		// Skip test files if not included
		if !s.config.IncludeTests && strings.HasSuffix(info.Name(), "_test.go") {
			if s.config.Verbose {
				fmt.Printf("Skipping test file: %s\n", path)
			}
			return nil
		}
		
		// Check file extension (only .go files for now)
		if !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}
		
		// Check if file should be excluded
		for _, excludeFile := range s.config.ExcludeFiles {
			if info.Name() == excludeFile {
				if s.config.Verbose {
					fmt.Printf("Skipping excluded file: %s\n", path)
				}
				return nil
			}
		}
		
		// Add file to the list
		relPath, err := filepath.Rel(s.rootPath, path)
		if err != nil {
			relPath = path
		}
		
		file := &models.File{
			Path:     path,
			RelPath:  relPath,
			Size:     info.Size(),
			ModTime:  info.ModTime(),
			IsVendor: strings.Contains(path, "vendor/"),
		}
		
		files = append(files, file)
		
		return nil
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to scan repository: %w", err)
	}
	
	return files, nil
}

// GetRepositoryInfo returns information about the repository
func (s *Scanner) GetRepositoryInfo() (*models.Repository, error) {
	// Check if it's a Git repository
	gitDir := filepath.Join(s.rootPath, ".git")
	isGit, _ := exists(gitDir)
	
	repo := &models.Repository{
		Path:  s.rootPath,
		IsGit: isGit,
	}
	
	// If it's a Git repository, try to get additional information
	if isGit {
		// This is a placeholder for Git repository information retrieval
		// In a real implementation, we would use a Git library or execute Git commands
		// to get information like remote URL, current branch, etc.
		repo.VCS = "git"
	}
	
	return repo, nil
}

// exists checks if a file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
