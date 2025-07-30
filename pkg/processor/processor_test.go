package processor_test

import (
	"fmt"
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

	t.Run("permission denied for output directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputDir := filepath.Join(tmpDir, "input")
		outputDir := "/root/restricted" // System directory that should be restricted

		if err := os.MkdirAll(inputDir, 0755); err != nil {
			t.Fatalf("Failed to create input directory: %v", err)
		}

		// Create a test file
		testFile := filepath.Join(inputDir, "test.md")
		if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		options := processor.ProcessOptions{
			OutputDir: outputDir,
			Pattern:   "*.md",
			Recursive: false,
		}

		err := proc.ProcessDirectory(inputDir, options)
		if err == nil {
			t.Error("Expected error for restricted output directory")
		}
	})

	t.Run("conversion error", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputDir := filepath.Join(tmpDir, "input")
		outputDir := filepath.Join(tmpDir, "output")

		if err := os.MkdirAll(inputDir, 0755); err != nil {
			t.Fatalf("Failed to create input directory: %v", err)
		}

		// Create a test file
		testFile := filepath.Join(inputDir, "test.md")
		if err := os.WriteFile(testFile, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		// Use a mock converter that always fails
		mockConv := &MockFailingConverter{}
		failingProc := processor.NewFileProcessor(mockConv)

		options := processor.ProcessOptions{
			OutputDir: outputDir,
			Pattern:   "*.md",
			Recursive: false,
		}

		err := failingProc.ProcessDirectory(inputDir, options)
		if err == nil {
			t.Error("Expected error from failing converter")
		}
	})
}

// TestGetOutputPath_EdgeCases tests edge cases for output path calculation
func TestGetOutputPath_EdgeCases(t *testing.T) {
	tests := []struct {
		name      string
		inputFile string
		inputDir  string
		outputDir string
		expected  string
	}{
		{
			name:      "file with no extension",
			inputFile: "/input/noext",
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/output/noext.html",
		},
		{
			name:      "file with multiple extensions",
			inputFile: "/input/file.backup.md",
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/output/file.backup.html",
		},
		{
			name:      "deeply nested file",
			inputFile: "/input/a/b/c/d/file.md",
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/output/a/b/c/d/file.html",
		},
		{
			name:      "special characters in filename",
			inputFile: "/input/file with spaces & symbols.md",
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/output/file with spaces & symbols.html",
		},
		{
			name:      "different root path",
			inputFile: "/completely/different/path/file.md", // Different root path
			inputDir:  "/input",
			outputDir: "/output",
			expected:  "/completely/different/path/file.html", // Relative path preserved
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

// TestFileProcessor_FindFiles_EdgeCases tests file finding edge cases
func TestFileProcessor_FindFiles_EdgeCases(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create various file types
	files := []struct {
		path    string
		content string
	}{
		{"normal.md", "# Normal"},
		{"README.MD", "# Uppercase extension"},
		{".hidden.md", "# Hidden file"},
		{"no-extension", "# No extension"},
		{"subdir/nested.md", "# Nested"},
		{"subdir/deep/nested.md", "# Deep nested"},
		{"other.txt", "Not markdown"},
	}

	for _, file := range files {
		fullPath := filepath.Join(tmpDir, file.path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create directory for %s: %v", file.path, err)
		}
		if err := os.WriteFile(fullPath, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file.path, err)
		}
	}

	tests := []struct {
		name      string
		pattern   string
		recursive bool
		expected  int // Number of files expected
	}{
		{
			name:      "case insensitive pattern",
			pattern:   "*.MD",
			recursive: false,
			expected:  1, // Only README.MD
		},
		{
			name:      "hidden files",
			pattern:   ".*md",
			recursive: false,
			expected:  1, // Only .hidden.md
		},
		{
			name:      "no matches",
			pattern:   "*.xyz",
			recursive: false,
			expected:  0,
		},
		{
			name:      "recursive with multiple levels",
			pattern:   "*.md",
			recursive: true,
			expected:  4, // normal.md, .hidden.md, subdir/nested.md, subdir/deep/nested.md
		},
	}

	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	proc := processor.NewFileProcessor(conv)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputDir := filepath.Join(tmpDir, "output_"+tt.name)
			
			options := processor.ProcessOptions{
				OutputDir: outputDir,
				Pattern:   tt.pattern,
				Recursive: tt.recursive,
			}

			err := proc.ProcessDirectory(tmpDir, options)
			if err != nil && tt.expected > 0 {
				t.Errorf("ProcessDirectory() unexpected error = %v", err)
				return
			}
			if err == nil && tt.expected == 0 {
				// Should not fail even if no files found
				return
			}

			// Count output files
			if tt.expected > 0 {
				outputFiles := countHTMLFiles(t, outputDir)
				if outputFiles != tt.expected {
					t.Errorf("Expected %d output files, got %d", tt.expected, outputFiles)
				}
			}
		})
	}
}

// TestFileProcessor_LargeFiles tests processing of large files
func TestFileProcessor_LargeFiles(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large file test in short mode")
	}

	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")

	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create a large markdown file (1MB+)
	largeContent := strings.Repeat("# Section\n\nLorem ipsum dolor sit amet, consectetur adipiscing elit. ", 10000)
	largeFile := filepath.Join(inputDir, "large.md")
	if err := os.WriteFile(largeFile, []byte(largeContent), 0644); err != nil {
		t.Fatalf("Failed to create large file: %v", err)
	}

	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	proc := processor.NewFileProcessor(conv)

	options := processor.ProcessOptions{
		OutputDir: outputDir,
		Pattern:   "*.md",
		Recursive: false,
	}

	err := proc.ProcessDirectory(inputDir, options)
	if err != nil {
		t.Errorf("ProcessDirectory() failed on large file: %v", err)
	}

	// Verify output was created
	outputFile := filepath.Join(outputDir, "large.html")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Large file output was not created")
	}
}

// TestFileProcessor_ConcurrentProcessing tests concurrent safety (if applicable)
func TestFileProcessor_ConcurrentProcessing(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create multiple input directories
	for i := 0; i < 3; i++ {
		inputDir := filepath.Join(tmpDir, fmt.Sprintf("input%d", i))
		if err := os.MkdirAll(inputDir, 0755); err != nil {
			t.Fatalf("Failed to create input directory %d: %v", i, err)
		}
		
		testFile := filepath.Join(inputDir, "test.md")
		if err := os.WriteFile(testFile, []byte(fmt.Sprintf("# Test %d", i)), 0644); err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
	}

	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	proc := processor.NewFileProcessor(conv)

	// Process directories concurrently
	errChan := make(chan error, 3)
	
	for i := 0; i < 3; i++ {
		go func(index int) {
			inputDir := filepath.Join(tmpDir, fmt.Sprintf("input%d", index))
			outputDir := filepath.Join(tmpDir, fmt.Sprintf("output%d", index))
			
			options := processor.ProcessOptions{
				OutputDir: outputDir,
				Pattern:   "*.md",
				Recursive: false,
			}
			
			err := proc.ProcessDirectory(inputDir, options)
			errChan <- err
		}(i)
	}

	// Collect results
	for i := 0; i < 3; i++ {
		err := <-errChan
		if err != nil {
			t.Errorf("Concurrent processing %d failed: %v", i, err)
		}
	}
}

// Helper function to count HTML files in a directory
func countHTMLFiles(t *testing.T, dir string) int {
	t.Helper()
	
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0
	}

	count := 0
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(path, ".html") {
			count++
		}
		return nil
	})
	
	if err != nil {
		t.Errorf("Error counting HTML files: %v", err)
		return 0
	}
	
	return count
}

// Mock converter that always fails
type MockFailingConverter struct{}

func (m *MockFailingConverter) Convert(input []byte) ([]byte, error) {
	return nil, fmt.Errorf("mock conversion error")
}

func (m *MockFailingConverter) ConvertFile(inputPath, outputPath string) error {
	return fmt.Errorf("mock file conversion error")
}