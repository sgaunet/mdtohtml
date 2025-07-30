package processor_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/processor"
)

// TestGetOutputPath tests the output path calculation
func TestGetOutputPath(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		inputDir  string
		outputDir string
		expected  string
	}{
		{
			name:      "simple file",
			inputFile: "/input/test.md",
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/output/test.html",
		},
		{
			name:      "nested file",
			inputFile: "/input/subdir/test.md",
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/output/subdir/test.html",
		},
		{
			name:      "relative paths",
			inputFile: "docs/readme.md",
			inputDir:  "docs",
			outputDir: "html",
			expected:  "html/readme.html",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := processor.GetOutputPath(tt.inputFile, tt.inputDir, tt.outputDir)
			// Normalize paths for cross-platform testing
			result = filepath.ToSlash(result)
			expected := filepath.ToSlash(tt.expected)
			
			if result != expected {
				t.Errorf("GetOutputPath() = %q, want %q", result, expected)
			}
		})
	}
}

// TestFileProcessor_ProcessDirectory tests batch processing functionality
func TestFileProcessor_ProcessDirectory(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	
	// Create test directory structure
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}
	
	subDir := filepath.Join(inputDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create test files
	testFiles := []struct {
		path    string
		content string
	}{
		{
			path:    filepath.Join(inputDir, "test1.md"),
			content: "# Test 1\n\nThis is test file 1.",
		},
		{
			path:    filepath.Join(inputDir, "test2.md"),
			content: "# Test 2\n\nThis is test file 2.",
		},
		{
			path:    filepath.Join(subDir, "nested.md"),
			content: "# Nested\n\nThis is a nested file.",
		},
		{
			path:    filepath.Join(inputDir, "not-markdown.txt"),
			content: "This should be ignored.",
		},
	}

	for _, file := range testFiles {
		if err := os.WriteFile(file.path, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.path, err)
		}
	}

	// Test non-recursive processing
	t.Run("non-recursive", func(t *testing.T) {
		conv := converter.NewCompleteConverter(converter.DefaultOptions())
		proc := processor.NewFileProcessor(conv)

		options := processor.ProcessOptions{
			OutputDir: outputDir,
			Pattern:   "*.md",
			Recursive: false,
		}

		if err := proc.ProcessDirectory(inputDir, options); err != nil {
			t.Fatalf("ProcessDirectory() error = %v", err)
		}

		// Verify output files exist
		expectedFiles := []string{
			filepath.Join(outputDir, "test1.html"),
			filepath.Join(outputDir, "test2.html"),
		}

		for _, expectedFile := range expectedFiles {
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				t.Errorf("Expected output file %s does not exist", expectedFile)
			} else {
				// Verify content
				content, err := os.ReadFile(expectedFile)
				if err != nil {
					t.Errorf("Failed to read output file %s: %v", expectedFile, err)
					continue
				}
				
				contentStr := string(content)
				if !strings.Contains(contentStr, "<!DOCTYPE html>") {
					t.Errorf("Output file %s does not contain HTML doctype", expectedFile)
				}
			}
		}

		// Verify nested file was not processed
		nestedOutput := filepath.Join(outputDir, "subdir", "nested.html")
		if _, err := os.Stat(nestedOutput); !os.IsNotExist(err) {
			t.Errorf("Unexpected nested output file exists: %s", nestedOutput)
		}
	})

	// Clean up for recursive test
	if err := os.RemoveAll(outputDir); err != nil {
		t.Fatalf("Failed to clean up output directory: %v", err)
	}

	// Test recursive processing
	t.Run("recursive", func(t *testing.T) {
		conv := converter.NewCompleteConverter(converter.DefaultOptions())
		proc := processor.NewFileProcessor(conv)

		options := processor.ProcessOptions{
			OutputDir: outputDir,
			Pattern:   "*.md",
			Recursive: true,
		}

		if err := proc.ProcessDirectory(inputDir, options); err != nil {
			t.Fatalf("ProcessDirectory() error = %v", err)
		}

		// Verify all output files exist, including nested
		expectedFiles := []string{
			filepath.Join(outputDir, "test1.html"),
			filepath.Join(outputDir, "test2.html"),
			filepath.Join(outputDir, "subdir", "nested.html"),
		}

		for _, expectedFile := range expectedFiles {
			if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
				t.Errorf("Expected output file %s does not exist", expectedFile)
			}
		}
	})
}

// TestFileProcessor_ProcessDirectory_Errors tests error conditions
func TestFileProcessor_ProcessDirectory_Errors(t *testing.T) {
	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	proc := processor.NewFileProcessor(conv)

	t.Run("nonexistent directory", func(t *testing.T) {
		options := processor.ProcessOptions{
			OutputDir: "/tmp",
			Pattern:   "*.md",
			Recursive: false,
		}

		err := proc.ProcessDirectory("/nonexistent/directory", options)
		if err == nil {
			t.Error("Expected error for nonexistent directory, got nil")
		}
	})
}