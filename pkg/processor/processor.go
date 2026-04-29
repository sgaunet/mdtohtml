// Package processor handles batch processing of markdown files
package processor

import (
	"fmt"
	"path/filepath"
	"strings"
)

// BatchProcessor defines the interface for processing multiple files.
type BatchProcessor interface {
	// ProcessDirectory processes all matching files in a directory
	ProcessDirectory(dir string, options ProcessOptions) error
}

// DefaultOutputExt is the default extension applied when ProcessOptions.OutputExt is empty.
const DefaultOutputExt = ".html"

// ProcessOptions configures batch processing behavior.
type ProcessOptions struct {
	// OutputDir is the directory where converted files will be written
	OutputDir string

	// Pattern is the file pattern to match (e.g., "*.md")
	Pattern string

	// Recursive determines if subdirectories should be processed
	Recursive bool

	// OutputExt is the extension applied to converted files (e.g., ".html", ".pdf").
	// Empty defaults to DefaultOutputExt.
	OutputExt string
}

// FileInfo represents information about a file to be processed.
type FileInfo struct {
	InputPath  string
	OutputPath string
}

// GetOutputPath calculates the output path for a given input file using the
// default ".html" extension. Use GetOutputPathExt to choose an extension.
func GetOutputPath(inputFile, inputDir, outputDir string) string {
	return GetOutputPathExt(inputFile, inputDir, outputDir, DefaultOutputExt)
}

// GetOutputPathExt is like GetOutputPath but applies the supplied extension
// (which must include the leading dot, e.g. ".pdf").
func GetOutputPathExt(inputFile, inputDir, outputDir, ext string) string {
	relPath, err := filepath.Rel(inputDir, inputFile)
	if err != nil {
		relPath = filepath.Base(inputFile)
	}

	outputFile := relPath[:len(relPath)-len(filepath.Ext(relPath))] + ext
	return filepath.Join(outputDir, outputFile)
}

// ValidateOutputPath checks that outputPath is within the outputDir.
func ValidateOutputPath(outputPath, outputDir string) error {
	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("%w: cannot resolve output path '%s': %w", ErrPathTraversal, outputPath, err)
	}

	absDir, err := filepath.Abs(outputDir)
	if err != nil {
		return fmt.Errorf("%w: cannot resolve output directory '%s': %w", ErrPathTraversal, outputDir, err)
	}

	if !strings.HasPrefix(absOutput, absDir+string(filepath.Separator)) {
		return fmt.Errorf("%w: output path '%s' escapes output directory '%s'", ErrPathTraversal, outputPath, outputDir)
	}

	return nil
}