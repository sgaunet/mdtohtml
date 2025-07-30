// Package processor handles batch processing of markdown files
package processor

import (
	"path/filepath"
)

// BatchProcessor defines the interface for processing multiple files.
type BatchProcessor interface {
	// ProcessDirectory processes all matching files in a directory
	ProcessDirectory(dir string, options ProcessOptions) error
}

// ProcessOptions configures batch processing behavior.
type ProcessOptions struct {
	// OutputDir is the directory where HTML files will be written
	OutputDir string

	// Pattern is the file pattern to match (e.g., "*.md")
	Pattern string

	// Recursive determines if subdirectories should be processed
	Recursive bool
}

// FileInfo represents information about a file to be processed.
type FileInfo struct {
	InputPath  string
	OutputPath string
}

// GetOutputPath calculates the output path for a given input file.
func GetOutputPath(inputFile, inputDir, outputDir string) string {
	relPath, err := filepath.Rel(inputDir, inputFile)
	if err != nil {
		relPath = filepath.Base(inputFile)
	}

	outputFile := relPath[:len(relPath)-len(filepath.Ext(relPath))] + ".html"
	return filepath.Join(outputDir, outputFile)
}