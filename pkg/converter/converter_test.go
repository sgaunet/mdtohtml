package converter_test

import (
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