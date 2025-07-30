package validator_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/validator"
)

// TestGoldmarkValidator_Validate tests markdown validation functionality
func TestGoldmarkValidator_Validate(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		shouldErr bool
	}{
		{
			name:      "valid simple markdown",
			input:     "# Hello World\n\nThis is a test.",
			shouldErr: false,
		},
		{
			name:      "valid complex markdown",
			input:     "# Title\n\n**Bold** and *italic* text.\n\n- List item\n- Another item\n\n```go\nfunc main() {}\n```",
			shouldErr: false,
		},
		{
			name:      "valid github flavored markdown",
			input:     "| Header 1 | Header 2 |\n|----------|----------|\n| Cell 1   | Cell 2   |",
			shouldErr: false,
		},
		{
			name:      "empty content",
			input:     "",
			shouldErr: false, // Empty content is valid
		},
		{
			name:      "only whitespace",
			input:     "   \n\n  \t  \n",
			shouldErr: false, // Whitespace-only content is valid
		},
		{
			name:      "unicode content",
			input:     "# Unicode Test\n\n你好世界 🌍 العالم мир",
			shouldErr: false,
		},
		{
			name:      "malformed table",
			input:     "| Header 1 | Header 2\n|----------|\n| Cell 1   | Cell 2   |",
			shouldErr: false, // Goldmark is generally tolerant of malformed markdown
		},
	}

	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	val := validator.NewGoldmarkValidator(conv)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := val.Validate([]byte(tt.input))
			if (err != nil) != tt.shouldErr {
				t.Errorf("Validate() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}
}

// TestGoldmarkValidator_ValidateFile tests file-based validation functionality
func TestGoldmarkValidator_ValidateFile(t *testing.T) {
	// Create temporary test files
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		filename  string
		content   string
		shouldErr bool
	}{
		{
			name:      "valid markdown file",
			filename:  "valid.md",
			content:   "# Valid File\n\nThis is valid markdown content.",
			shouldErr: false,
		},
		{
			name:      "empty file",
			filename:  "empty.md",
			content:   "",
			shouldErr: false,
		},
		{
			name:      "complex valid file",
			filename:  "complex.md",
			content:   "---\ntitle: YAML frontmatter\n---\n\n# Document\n\nContent here.",
			shouldErr: false, // Goldmark handles YAML frontmatter gracefully
		},
	}

	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	val := validator.NewGoldmarkValidator(conv)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test file
			filePath := filepath.Join(tmpDir, tt.filename)
			if err := os.WriteFile(filePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Validate file
			err := val.ValidateFile(filePath)
			if (err != nil) != tt.shouldErr {
				t.Errorf("ValidateFile() error = %v, shouldErr %v", err, tt.shouldErr)
			}
		})
	}

	// Test nonexistent file
	t.Run("nonexistent file", func(t *testing.T) {
		err := val.ValidateFile(filepath.Join(tmpDir, "nonexistent.md"))
		if err == nil {
			t.Error("ValidateFile() should return error for nonexistent file")
		}
		if !strings.Contains(err.Error(), "no such file") && !strings.Contains(err.Error(), "cannot find") {
			t.Errorf("ValidateFile() error should indicate file not found, got: %v", err)
		}
	})
}

// TestGoldmarkValidator_WithDifferentOptions tests validation with different converter options
func TestGoldmarkValidator_WithDifferentOptions(t *testing.T) {
	input := []byte("# Test\n\nSmart quotes: \"hello\" and fractions: 1/2")

	tests := []struct {
		name    string
		options converter.Options
	}{
		{
			name:    "default options",
			options: converter.DefaultOptions(),
		},
		{
			name: "no smart punctuation",
			options: converter.Options{
				SmartPunctuation: false,
				LaTeXDashes:      false,
				Fractions:        false,
			},
		},
		{
			name: "only fractions",
			options: converter.Options{
				SmartPunctuation: false,
				LaTeXDashes:      false,
				Fractions:        true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := converter.NewGoldmarkConverter(tt.options)
			val := validator.NewGoldmarkValidator(conv)

			err := val.Validate(input)
			if err != nil {
				t.Errorf("Validate() with %s options failed: %v", tt.name, err)
			}
		})
	}
}

// BenchmarkGoldmarkValidator_Validate benchmarks validation performance
func BenchmarkGoldmarkValidator_Validate(b *testing.B) {
	input := []byte(`# Benchmark Test Document

This is a test document for benchmarking the validation performance.

## Features

- **Bold text** and *italic text*
- [Links](http://example.com)  
- ` + "`code blocks`" + `

### Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `

### Table

| Column 1 | Column 2 | Column 3 |
|----------|----------|----------|
| Value 1  | Value 2  | Value 3  |
| Data A   | Data B   | Data C   |

More content here to make the document longer and more realistic for benchmarking purposes.
`)

	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	val := validator.NewGoldmarkValidator(conv)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := val.Validate(input)
		if err != nil {
			b.Fatalf("Validate() error = %v", err)
		}
	}
}