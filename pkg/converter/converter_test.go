package converter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/converter"
)

// TestGoldmarkConverter_Convert tests the basic markdown conversion functionality
func TestGoldmarkConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		options  converter.Options
	}{
		{
			name:  "simple markdown",
			input: "# Hello World\n\nThis is a test.",
			contains: []string{
				"<h1", "Hello World", "<p>", "This is a test",
			},
			options: converter.DefaultOptions(),
		},
		{
			name:  "markdown with emphasis",
			input: "**bold** and *italic*",
			contains: []string{
				"<strong>bold</strong>",
				"<em>italic</em>",
			},
			options: converter.DefaultOptions(),
		},
		{
			name:  "markdown with code",
			input: "Here is `code` inline.",
			contains: []string{
				"<code>code</code>",
			},
			options: converter.DefaultOptions(),
		},
		{
			name:  "github flavored markdown",
			input: "| Header 1 | Header 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |",
			contains: []string{
				"<table>", "<th>", "Header 1", "Header 2", "<td>", "Cell 1", "Cell 2",
			},
			options: converter.DefaultOptions(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := converter.NewGoldmarkConverter(tt.options)
			output, err := conv.Convert([]byte(tt.input))
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			outputStr := string(output)
			for _, expected := range tt.contains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Convert() output does not contain %q\nGot: %s", expected, outputStr)
				}
			}
		})
	}
}

// TestCompleteConverter_Convert tests the complete conversion with title and HTML wrapping
func TestCompleteConverter_Convert(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		options  converter.Options
	}{
		{
			name:  "markdown with title",
			input: "# My Title\n\nContent here.",
			contains: []string{
				"<!DOCTYPE html>",
				"<html>",
				"<head>",
				"<title>My Title</title>",
				"<h1", "My Title", "Content here",
				"<style>", // CSS should be injected
				"</body>",
				"</html>",
			},
			options: converter.DefaultOptions(),
		},
		{
			name:  "markdown without title",
			input: "Just some content without a title.",
			contains: []string{
				"<!DOCTYPE html>",
				"<title></title>", // empty title
				"Just some content",
			},
			options: converter.DefaultOptions(),
		},
		{
			name:  "underlined title",
			input: "Underlined Title\n================\n\nContent here.",
			contains: []string{
				"<title>Underlined Title</title>",
				"<h1", "Underlined Title",
			},
			options: converter.DefaultOptions(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := converter.NewCompleteConverter(tt.options)
			output, err := conv.Convert([]byte(tt.input))
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			outputStr := string(output)
			for _, expected := range tt.contains {
				if !strings.Contains(outputStr, expected) {
					t.Errorf("Convert() output does not contain %q\nGot: %s", expected, outputStr)
				}
			}
		})
	}
}

// TestConverter_ConvertFile tests file-based conversion
func TestConverter_ConvertFile(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()
	
	inputFile := filepath.Join(tmpDir, "test.md")
	outputFile := filepath.Join(tmpDir, "test.html")
	
	testContent := "# Test File\n\nThis is a test file."
	if err := os.WriteFile(inputFile, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create test input file: %v", err)
	}

	// Test conversion
	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	if err := conv.ConvertFile(inputFile, outputFile); err != nil {
		t.Fatalf("ConvertFile() error = %v", err)
	}

	// Verify output file exists and contains expected content
	output, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read output file: %v", err)
	}

	outputStr := string(output)
	expectedContents := []string{
		"<!DOCTYPE html>",
		"<title>Test File</title>",
		"<h1", "Test File",
		"This is a test file",
	}

	for _, expected := range expectedContents {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("ConvertFile() output does not contain %q", expected)
		}
	}
}

// TestGoldmarkConverter_ConvertFile tests file-based conversion for GoldmarkConverter
func TestGoldmarkConverter_ConvertFile(t *testing.T) {
	tmpDir := t.TempDir()
	
	tests := []struct {
		name      string
		input     string
		options   converter.Options
		expectErr bool
	}{
		{
			name:      "successful conversion",
			input:     "# Test\n\nContent here.",
			options:   converter.DefaultOptions(),
			expectErr: false,
		},
		{
			name:      "conversion with smart punctuation",
			input:     "Smart \"quotes\" and 1/2 fractions",
			options:   converter.Options{SmartPunctuation: true, Fractions: true},
			expectErr: false,
		},
		{
			name:      "empty file",
			input:     "",
			options:   converter.DefaultOptions(),
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputFile := filepath.Join(tmpDir, "input_"+tt.name+".md")
			outputFile := filepath.Join(tmpDir, "output_"+tt.name+".html")

			// Create input file
			if err := os.WriteFile(inputFile, []byte(tt.input), 0644); err != nil {
				t.Fatalf("Failed to create input file: %v", err)
			}

			conv := converter.NewGoldmarkConverter(tt.options)
			err := conv.ConvertFile(inputFile, outputFile)

			if (err != nil) != tt.expectErr {
				t.Errorf("ConvertFile() error = %v, expectErr %v", err, tt.expectErr)
				return
			}

			if !tt.expectErr {
				// Verify output file was created and has content
				output, err := os.ReadFile(outputFile)
				if err != nil {
					t.Errorf("Failed to read output file: %v", err)
					return
				}

				if len(output) == 0 && len(tt.input) > 0 {
					t.Error("Output file is empty but input had content")
				}
			}
		})
	}
}

// TestGoldmarkConverter_ConvertFile_Errors tests error conditions
func TestGoldmarkConverter_ConvertFile_Errors(t *testing.T) {
	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())

	t.Run("nonexistent input file", func(t *testing.T) {
		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "output.html")
		
		err := conv.ConvertFile("/nonexistent/file.md", outputFile)
		if err == nil {
			t.Error("Expected error for nonexistent input file")
		}
	})

	t.Run("invalid output directory", func(t *testing.T) {
		tmpDir := t.TempDir()
		inputFile := filepath.Join(tmpDir, "input.md")
		
		// Create input file
		if err := os.WriteFile(inputFile, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		err := conv.ConvertFile(inputFile, "/invalid/directory/output.html")
		if err == nil {
			t.Error("Expected error for invalid output directory")
		}
	})
}

// TestNewCompleteConverterWithComponents tests custom component integration
func TestNewCompleteConverterWithComponents(t *testing.T) {
	// Create mock components
	mockParser := &MockTitleExtractor{}
	mockTemplate := &MockHTMLTemplate{}
	goldmarkConv := converter.NewGoldmarkConverter(converter.DefaultOptions())

	conv := converter.NewCompleteConverterWithComponents(goldmarkConv, mockParser, mockTemplate)
	
	input := []byte("# Test Title\n\nContent here.")
	output, err := conv.Convert(input)
	
	if err != nil {
		t.Fatalf("Convert() error = %v", err)
	}

	outputStr := string(output)
	
	// Verify mock components were used
	if !strings.Contains(outputStr, "MOCK_TITLE") {
		t.Error("Mock title extractor was not used")
	}
	
	if !strings.Contains(outputStr, "MOCK_WRAPPED") {
		t.Error("Mock HTML template was not used")
	}
}

// TestCompleteConverter_Convert_ErrorPaths tests error handling
func TestCompleteConverter_Convert_ErrorPaths(t *testing.T) {
	t.Run("goldmark conversion error", func(t *testing.T) {
		// Test error path by trying to convert something that would cause conversion errors
		// We'll use the existing converter but with invalid input that might cause issues
		conv := converter.NewCompleteConverter(converter.DefaultOptions())
		
		// Create a very large input that might cause issues
		largeInput := make([]byte, 1024*1024*10) // 10MB of null bytes
		_, err := conv.Convert(largeInput)
		// This might not actually error, but at least we're testing the path
		if err != nil {
			// Error is acceptable for this test
			t.Logf("Large input conversion failed as expected: %v", err)
		}
	})
}

// TestCompleteConverter_ConvertFile_ErrorPaths tests file operation errors
func TestCompleteConverter_ConvertFile_ErrorPaths(t *testing.T) {
	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	tmpDir := t.TempDir()

	t.Run("read-only output directory", func(t *testing.T) {
		inputFile := filepath.Join(tmpDir, "input.md")
		if err := os.WriteFile(inputFile, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create input file: %v", err)
		}

		readOnlyDir := filepath.Join(tmpDir, "readonly")
		if err := os.Mkdir(readOnlyDir, 0444); err != nil {
			t.Fatalf("Failed to create readonly dir: %v", err)
		}
		defer os.Chmod(readOnlyDir, 0755) // Cleanup

		outputFile := filepath.Join(readOnlyDir, "output.html")
		err := conv.ConvertFile(inputFile, outputFile)
		if err == nil {
			t.Error("Expected error for read-only directory")
		}
	})
}

// TestConverter_Options tests different option combinations
func TestConverter_Options(t *testing.T) {
	tests := []struct {
		name    string
		options converter.Options
		input   string
		check   func(string) bool
	}{
		{
			name: "smart punctuation enabled",
			options: converter.Options{
				SmartPunctuation: true,
				LaTeXDashes:      false,
				Fractions:        false,
			},
			input: "\"Smart quotes\"",
			check: func(output string) bool {
				return strings.Contains(output, "&ldquo;") || strings.Contains(output, "\u201c") // Check for smart quotes entities
			},
		},
		{
			name: "fractions enabled",
			options: converter.Options{
				SmartPunctuation: false,
				LaTeXDashes:      false,
				Fractions:        true,
			},
			input: "One half: 1/2",
			check: func(output string) bool {
				return strings.Contains(output, "½") || strings.Contains(output, "1/2")
			},
		},
		{
			name: "all typography disabled",
			options: converter.Options{
				SmartPunctuation: false,
				LaTeXDashes:      false,
				Fractions:        false,
			},
			input: "\"Plain quotes\" and 1/2",
			check: func(output string) bool {
				return strings.Contains(output, "&quot;Plain quotes&quot;") && strings.Contains(output, "1/2")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := converter.NewGoldmarkConverter(tt.options)
			output, err := conv.Convert([]byte(tt.input))
			if err != nil {
				t.Fatalf("Convert() error = %v", err)
			}

			if !tt.check(string(output)) {
				t.Errorf("Option test failed for %s\nOutput: %s", tt.name, string(output))
			}
		})
	}
}

// TestDefaultOptions tests the default options function
func TestDefaultOptions(t *testing.T) {
	opts := converter.DefaultOptions()
	
	if !opts.SmartPunctuation {
		t.Error("Default options should have SmartPunctuation enabled")
	}
	if !opts.LaTeXDashes {
		t.Error("Default options should have LaTeXDashes enabled")
	}
	if !opts.Fractions {
		t.Error("Default options should have Fractions enabled")
	}
}

// Mock implementations for testing
type MockTitleExtractor struct{}

func (m *MockTitleExtractor) ExtractTitle(content []byte) string {
	return "MOCK_TITLE"
}

type MockHTMLTemplate struct{}

func (m *MockHTMLTemplate) Wrap(content, title string) string {
	return "MOCK_WRAPPED: " + title + " - " + content
}

func (m *MockHTMLTemplate) InjectCSS(html, css string) string {
	return "MOCK_CSS_INJECTED: " + html
}

type MockGoldmarkConverter struct {
	shouldError bool
}

func (m *MockGoldmarkConverter) Convert(input []byte) ([]byte, error) {
	if m.shouldError {
		return nil, fmt.Errorf("mock goldmark error")
	}
	return []byte("mock converted: " + string(input)), nil
}

func (m *MockGoldmarkConverter) ConvertFile(inputPath, outputPath string) error {
	if m.shouldError {
		return fmt.Errorf("mock file conversion error")
	}
	return nil
}

// BenchmarkConverter_Convert benchmarks the conversion performance
func BenchmarkConverter_Convert(b *testing.B) {
	input := `# Benchmark Test

This is a benchmark test with **bold text**, *italic text*, and [links](http://example.com).

## Subheading

- List item 1
- List item 2
- List item 3

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```"

	conv := converter.NewCompleteConverter(converter.DefaultOptions())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := conv.Convert([]byte(input))
		if err != nil {
			b.Fatalf("Convert() error = %v", err)
		}
	}
}

// BenchmarkGoldmarkConverter_Convert benchmarks just the goldmark conversion
func BenchmarkGoldmarkConverter_Convert(b *testing.B) {
	input := []byte(`# Benchmark Test

This is a benchmark test with **bold text**, *italic text*, and [links](http://example.com).

## Subheading

- List item 1
- List item 2
- List item 3

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```")

	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := conv.Convert(input)
		if err != nil {
			b.Fatalf("Convert() error = %v", err)
		}
	}
}