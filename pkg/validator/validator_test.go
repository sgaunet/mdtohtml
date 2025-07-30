package validator_test

import (
	"fmt"
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

// TestGoldmarkValidator_Validate_ErrorScenarios tests error conditions
func TestGoldmarkValidator_Validate_ErrorScenarios(t *testing.T) {
	tests := []struct {
		name        string
		setupConv   func() converter.Converter
		input       []byte
		expectError bool
	}{
		{
			name: "converter error",
			setupConv: func() converter.Converter {
				return &MockErrorConverter{}
			},
			input:       []byte("# Test"),
			expectError: true,
		},
		{
			name: "extremely large input",
			setupConv: func() converter.Converter {
				return converter.NewGoldmarkConverter(converter.DefaultOptions())
			},
			input:       make([]byte, 50*1024*1024), // 50MB
			expectError: false, // Should handle large inputs gracefully
		},
		{
			name: "binary input",
			setupConv: func() converter.Converter {
				return converter.NewGoldmarkConverter(converter.DefaultOptions())
			},
			input:       []byte{0x00, 0x01, 0x02, 0xFF, 0xFE, 0xFD},
			expectError: false, // Should handle binary data gracefully
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conv := tt.setupConv()
			val := validator.NewGoldmarkValidator(conv)
			
			err := val.Validate(tt.input)
			
			if (err != nil) != tt.expectError {
				t.Errorf("Validate() error = %v, expectError %v", err, tt.expectError)
			}
		})
	}
}

// TestGoldmarkValidator_ValidateFile_AdvancedErrors tests advanced file error scenarios
func TestGoldmarkValidator_ValidateFile_AdvancedErrors(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("permission denied", func(t *testing.T) {
		restrictedFile := filepath.Join(tmpDir, "restricted.md")
		if err := os.WriteFile(restrictedFile, []byte("# Test"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		
		// Make file unreadable
		if err := os.Chmod(restrictedFile, 0000); err != nil {
			t.Fatalf("Failed to change permissions: %v", err)
		}
		defer os.Chmod(restrictedFile, 0644) // Cleanup
		
		conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
		val := validator.NewGoldmarkValidator(conv)
		
		err := val.ValidateFile(restrictedFile)
		if err == nil {
			t.Error("Expected error for unreadable file")
		}
	})

	t.Run("empty file", func(t *testing.T) {
		emptyFile := filepath.Join(tmpDir, "empty.md")
		if err := os.WriteFile(emptyFile, []byte{}, 0644); err != nil {
			t.Fatalf("Failed to create empty file: %v", err)
		}
		
		conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
		val := validator.NewGoldmarkValidator(conv)
		
		err := val.ValidateFile(emptyFile)
		if err != nil {
			t.Errorf("Empty file should be valid: %v", err)
		}
	})

	t.Run("very large file", func(t *testing.T) {
		if testing.Short() {
			t.Skip("Skipping large file test in short mode")
		}
		
		largeFile := filepath.Join(tmpDir, "large.md")
		largeContent := strings.Repeat("# Section\n\nContent here.\n\n", 100000) // ~1.5MB
		if err := os.WriteFile(largeFile, []byte(largeContent), 0644); err != nil {
			t.Fatalf("Failed to create large file: %v", err)
		}
		
		conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
		val := validator.NewGoldmarkValidator(conv)
		
		err := val.ValidateFile(largeFile)
		if err != nil {
			t.Errorf("Large file validation failed: %v", err)
		}
	})
}

// TestGoldmarkValidator_ConcurrentValidation tests concurrent validation
func TestGoldmarkValidator_ConcurrentValidation(t *testing.T) {
	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	val := validator.NewGoldmarkValidator(conv)
	
	inputs := [][]byte{
		[]byte("# Test 1\n\nContent 1"),
		[]byte("# Test 2\n\nContent 2"),
		[]byte("# Test 3\n\nContent 3"),
		[]byte("# Test 4\n\nContent 4"),
		[]byte("# Test 5\n\nContent 5"),
	}
	
	errChan := make(chan error, len(inputs))
	
	for i, input := range inputs {
		go func(index int, content []byte) {
			err := val.Validate(content)
			errChan <- err
		}(i, input)
	}
	
	for i := 0; i < len(inputs); i++ {
		err := <-errChan
		if err != nil {
			t.Errorf("Concurrent validation %d failed: %v", i, err)
		}
	}
}

// TestGoldmarkValidator_MemoryUsage tests memory efficiency
func TestGoldmarkValidator_MemoryUsage(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping memory test in short mode")
	}
	
	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	val := validator.NewGoldmarkValidator(conv)
	
	// Validate many different inputs to test for memory leaks
	for i := 0; i < 1000; i++ {
		input := []byte(fmt.Sprintf("# Test %d\n\nThis is test content number %d with some **bold** and *italic* text.", i, i))
		err := val.Validate(input)
		if err != nil {
			t.Errorf("Validation %d failed: %v", i, err)
			break
		}
	}
}

// TestValidator_Interface tests that our implementation satisfies the interface
func TestValidator_Interface(t *testing.T) {
	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	
	// Ensure our implementation satisfies the interface
	var _ validator.Validator = validator.NewGoldmarkValidator(conv)
	
	t.Log("GoldmarkValidator correctly implements Validator interface")
}

// Mock converter that always returns an error
type MockErrorConverter struct{}

func (m *MockErrorConverter) Convert(input []byte) ([]byte, error) {
	return nil, fmt.Errorf("mock converter error")
}

func (m *MockErrorConverter) ConvertFile(inputPath, outputPath string) error {
	return fmt.Errorf("mock converter file error")
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

// BenchmarkGoldmarkValidator_ValidateFile benchmarks file validation performance
func BenchmarkGoldmarkValidator_ValidateFile(b *testing.B) {
	tmpDir := b.TempDir()
	testFile := filepath.Join(tmpDir, "benchmark.md")
	
	content := `# Benchmark Test Document

This is a test document for benchmarking file validation performance.

## Features

- **Bold text** and *italic text*
- [Links](http://example.com)  
- ` + "`code blocks`" + `

### Code Example

` + "```go\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}\n```" + `
`
	
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		b.Fatalf("Failed to create test file: %v", err)
	}

	conv := converter.NewGoldmarkConverter(converter.DefaultOptions())
	val := validator.NewGoldmarkValidator(conv)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := val.ValidateFile(testFile)
		if err != nil {
			b.Fatalf("ValidateFile() error = %v", err)
		}
	}
}