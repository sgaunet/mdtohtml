package pkg_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/processor"
	"github.com/sgaunet/mdtohtml/pkg/validator"
)

// TestCompleteWorkflow_EndToEnd tests the complete markdown processing workflow
func TestCompleteWorkflow_EndToEnd(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create test markdown files with various content
	testFiles := []struct {
		filename string
		content  string
	}{
		{
			"simple.md",
			"# Simple Document\n\nThis is a simple markdown document with **bold** text.",
		},
		{
			"complex.md",
			`# Complex Document

This document has various markdown features:

## Lists

- Item 1
- Item 2
- Item 3

## Code

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

## Table

| Column 1 | Column 2 |
|----------|----------|
| Cell 1   | Cell 2   |

## Links and Emphasis

This has [links](http://example.com) and *italic* text.
`,
		},
		{
			"with_title.md",
			"# Document Title\n\nContent with a clear title.",
		},
		{
			"underlined_title.md",
			"Underlined Title\n================\n\nContent with underlined title.",
		},
	}

	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create test files
	for _, file := range testFiles {
		path := filepath.Join(inputDir, file.filename)
		if err := os.WriteFile(path, []byte(file.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", file.filename, err)
		}
	}

	// Test complete conversion workflow
	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	proc := processor.NewFileProcessor(conv)

	options := processor.ProcessOptions{
		OutputDir: outputDir,
		Pattern:   "*.md",
		Recursive: false,
	}

	err := proc.ProcessDirectory(inputDir, options)
	if err != nil {
		t.Fatalf("Complete workflow failed: %v", err)
	}

	// Verify all output files were created
	for _, file := range testFiles {
		expectedOutput := strings.Replace(file.filename, ".md", ".html", 1)
		outputPath := filepath.Join(outputDir, expectedOutput)
		
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Errorf("Output file %s was not created", expectedOutput)
			continue
		}

		// Read and verify output content
		output, err := os.ReadFile(outputPath)
		if err != nil {
			t.Errorf("Failed to read output file %s: %v", expectedOutput, err)
			continue
		}

		outputStr := string(output)
		
		// Verify basic HTML structure
		if !strings.Contains(outputStr, "<!DOCTYPE html>") {
			t.Errorf("Output file %s missing DOCTYPE", expectedOutput)
		}
		if !strings.Contains(outputStr, "<html>") {
			t.Errorf("Output file %s missing html tag", expectedOutput)
		}
		if !strings.Contains(outputStr, "<head>") {
			t.Errorf("Output file %s missing head tag", expectedOutput)
		}
		if !strings.Contains(outputStr, "<body>") {
			t.Errorf("Output file %s missing body tag", expectedOutput)
		}
		if !strings.Contains(outputStr, "<style>") {
			t.Errorf("Output file %s missing CSS", expectedOutput)
		}
	}
}

// TestCustomComponentIntegration tests integration with custom components
func TestCustomComponentIntegration(t *testing.T) {
	// Create custom components
	customTemplate := &CustomTemplate{prefix: "CUSTOM"}
	customParser := &CustomParser{suffix: "_PARSED"}
	
	goldmarkConv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	completeConv := converter.NewCompleteConverterWithComponents(goldmarkConv, customParser, customTemplate)

	input := []byte("# Test Title\n\nTest content here.")
	output, err := completeConv.Convert(input)
	
	if err != nil {
		t.Fatalf("Custom component integration failed: %v", err)
	}

	outputStr := string(output)
	
	// Verify custom components were used
	if !strings.Contains(outputStr, "CUSTOM") {
		t.Error("Custom template was not used")
	}
	if !strings.Contains(outputStr, "_PARSED") {
		t.Error("Custom parser was not used")  
	}
}

// TestValidationWorkflow tests the validation workflow
func TestValidationWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	
	testCases := []struct {
		filename string
		content  string
		valid    bool
	}{
		{"valid.md", "# Valid Document\n\nThis is valid markdown.", true},
		{"empty.md", "", true},
		{"complex.md", "# Complex\n\n- List\n- Items\n\n```code```", true},
	}

	for _, tc := range testCases {
		filePath := filepath.Join(tmpDir, tc.filename)
		if err := os.WriteFile(filePath, []byte(tc.content), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}

		conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
		val := validator.NewGoldmarkValidator(conv)

		err := val.ValidateFile(filePath)
		
		if tc.valid && err != nil {
			t.Errorf("File %s should be valid but got error: %v", tc.filename, err)
		}
		if !tc.valid && err == nil {
			t.Errorf("File %s should be invalid but validation passed", tc.filename)
		}
	}
}

// TestStressWorkflow tests the system under stress
func TestStressWorkflow(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping stress test in short mode")
	}

	tmpDir := t.TempDir()
	inputDir := filepath.Join(tmpDir, "input")
	outputDir := filepath.Join(tmpDir, "output")
	
	if err := os.MkdirAll(inputDir, 0755); err != nil {
		t.Fatalf("Failed to create input directory: %v", err)
	}

	// Create many files
	numFiles := 20
	for i := 0; i < numFiles; i++ {
		filename := fmt.Sprintf("test_%d.md", i)
		content := fmt.Sprintf("# Test Document %d\n\nThis is test document number %d with some **bold** text and a [link](http://example.com).\n\n- List item 1\n- List item 2", i, i)
		
		path := filepath.Join(inputDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create test file %d: %v", i, err)
		}
	}

	// Process all files
	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	proc := processor.NewFileProcessor(conv)

	options := processor.ProcessOptions{
		OutputDir: outputDir,
		Pattern:   "*.md",
		Recursive: false,
	}

	err := proc.ProcessDirectory(inputDir, options)
	if err != nil {
		t.Fatalf("Stress test failed: %v", err)
	}

	// Verify all output files were created
	outputFiles, err := filepath.Glob(filepath.Join(outputDir, "*.html"))
	if err != nil {
		t.Fatalf("Failed to list output files: %v", err)
	}

	if len(outputFiles) != numFiles {
		t.Errorf("Expected %d output files, got %d", numFiles, len(outputFiles))
	}
}

// TestConcurrentWorkflow tests concurrent processing
func TestConcurrentWorkflow(t *testing.T) {
	tmpDir := t.TempDir()
	
	// Create separate directories for concurrent processing
	numWorkers := 5
	errChan := make(chan error, numWorkers)
	
	for worker := 0; worker < numWorkers; worker++ {
		go func(workerID int) {
			workerDir := filepath.Join(tmpDir, fmt.Sprintf("worker_%d", workerID))
			inputDir := filepath.Join(workerDir, "input")
			outputDir := filepath.Join(workerDir, "output")
			
			if err := os.MkdirAll(inputDir, 0755); err != nil {
				errChan <- fmt.Errorf("worker %d: failed to create input dir: %v", workerID, err)
				return
			}

			// Create test files for this worker
			for i := 0; i < 10; i++ {
				filename := fmt.Sprintf("worker_%d_file_%d.md", workerID, i)
				content := fmt.Sprintf("# Worker %d File %d\n\nContent for worker %d, file %d.", workerID, i, workerID, i)
				
				path := filepath.Join(inputDir, filename)
				if err := os.WriteFile(path, []byte(content), 0644); err != nil {
					errChan <- fmt.Errorf("worker %d: failed to create file %d: %v", workerID, i, err)
					return
				}
			}

			// Process files
			conv := converter.NewCompleteConverter(converter.DefaultOptions())
			proc := processor.NewFileProcessor(conv)

			options := processor.ProcessOptions{
				OutputDir: outputDir,
				Pattern:   "*.md",
				Recursive: false,
			}

			err := proc.ProcessDirectory(inputDir, options)
			errChan <- err
		}(worker)
	}

	// Collect results
	for i := 0; i < numWorkers; i++ {
		err := <-errChan
		if err != nil {
			t.Errorf("Concurrent worker failed: %v", err)
		}
	}
}

// TestAllComponentIntegration tests all components working together
func TestAllComponentIntegration(t *testing.T) {
	// Test different component combinations
	testCases := []struct {
		name        string
		options     converter.Options
		input       string
		checkOutput func(string) bool
	}{
		{
			name:    "default components",
			options: converter.DefaultOptions(),
			input:   "# Test\n\nSmart \"quotes\" and 1/2 fractions.",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<h1") && strings.Contains(output, "<title>Test</title>")
			},
		},
		{
			name: "no typography",
			options: converter.Options{
				SmartPunctuation: false,
				LaTeXDashes:      false,
				Fractions:        false,
			},
			input: "# Plain Text\n\nPlain \"quotes\" and 1/2.",
			checkOutput: func(output string) bool {
				return strings.Contains(output, "<title>Plain Text</title>") && strings.Contains(output, "&quot;quotes&quot;")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conv := converter.NewCompleteConverter(tc.options)
			output, err := conv.Convert([]byte(tc.input))
			
			if err != nil {
				t.Fatalf("Component integration failed: %v", err)
			}

			if !tc.checkOutput(string(output)) {
				t.Errorf("Output check failed for %s\nOutput: %s", tc.name, string(output))
			}
		})
	}
}

// Custom components for testing
type CustomTemplate struct {
	prefix string
}

func (c *CustomTemplate) Wrap(content, title string) string {
	return fmt.Sprintf("%s_WRAPPED: %s - %s", c.prefix, title, content)
}

func (c *CustomTemplate) InjectCSS(html, css string) string {
	return fmt.Sprintf("%s_CSS: %s", c.prefix, html)
}

type CustomParser struct {
	suffix string
}

func (c *CustomParser) ExtractTitle(content []byte) string {
	// Simple title extraction for testing
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			return strings.TrimPrefix(line, "# ") + c.suffix
		}
	}
	return "No Title" + c.suffix
}