package parser_test

import (
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/parser"
)

// TestMarkdownTitleExtractor_ExtractTitle tests title extraction functionality
func TestMarkdownTitleExtractor_ExtractTitle(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "prefix header",
			input:    "# My Title\n\nContent here.",
			expected: "My Title",
		},
		{
			name:     "prefix header with extra spaces",
			input:    "#  Title with Spaces  \n\nContent.",
			expected: "Title with Spaces",
		},
		{
			name:     "prefix header with tab",
			input:    "#\tTabbed Title\n\nContent.",
			expected: "Tabbed Title",
		},
		{
			name:     "underlined header with equals",
			input:    "Underlined Title\n================\n\nContent here.",
			expected: "Underlined Title",
		},
		{
			name:     "underlined header with trailing spaces",
			input:    "Title With Spaces   \n====     \n\nContent.",
			expected: "Title With Spaces",
		},
		{
			name:     "no title - regular text",
			input:    "This is just regular text.\n\nMore content.",
			expected: "",
		},
		{
			name:     "no title - level 2 header",
			input:    "## Level 2 Header\n\nContent here.",
			expected: "",
		},
		{
			name:     "no title - underlined with dashes",
			input:    "Not a Title\n-----------\n\nContent here.",
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "only whitespace",
			input:    "   \n\n  \n",
			expected: "",
		},
		{
			name:     "title after blank lines",
			input:    "\n\n# Title After Blanks\n\nContent.",
			expected: "Title After Blanks",
		},
		{
			name:     "underlined title after blank lines",
			input:    "\n\nTitle After Blanks\n==================\n\nContent.",
			expected: "Title After Blanks",
		},
		{
			name:     "invalid prefix header - no space",
			input:    "#NoSpaceTitle\n\nContent.",
			expected: "",
		},
		{
			name:     "incomplete underlined title",
			input:    "Title\n==\nMore content on same line",
			expected: "",
		},
		{
			name:     "carriage return line endings",
			input:    "# Title\r\n\r\nContent with CRLF.",
			expected: "Title",
		},
		{
			name:     "mixed line endings",
			input:    "Underlined Title\r\n================\r\n\nContent.",
			expected: "Underlined Title",
		},
	}

	extractor := parser.NewMarkdownTitleExtractor()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractor.ExtractTitle([]byte(tt.input))
			if result != tt.expected {
				t.Errorf("ExtractTitle() = %q, want %q", result, tt.expected)
			}
		})
	}
}

// TestMarkdownTitleExtractor_ExtractTitle_EdgeCases tests edge cases and boundary conditions
func TestMarkdownTitleExtractor_ExtractTitle_EdgeCases(t *testing.T) {
	extractor := parser.NewMarkdownTitleExtractor()

	t.Run("very long title", func(t *testing.T) {
		longTitle := "This is a very long title that exceeds normal length expectations and contains many words to test handling of lengthy content"
		input := "# " + longTitle + "\n\nContent."
		result := extractor.ExtractTitle([]byte(input))
		if result != longTitle {
			t.Errorf("ExtractTitle() = %q, want %q", result, longTitle)
		}
	})

	t.Run("title with special characters", func(t *testing.T) {
		title := "Title with Special Characters: @#$%^&*()_+-={}[]|\\:\";<>?,./~`"
		input := "# " + title + "\n\nContent."
		result := extractor.ExtractTitle([]byte(title))
		if result != "" { // Should not extract from plain text
			t.Errorf("ExtractTitle() should not extract from plain text, got %q", result)
		}
		
		// Test with proper markdown
		result = extractor.ExtractTitle([]byte(input))
		if result != title {
			t.Errorf("ExtractTitle() = %q, want %q", result, title)
		}
	})

	t.Run("unicode characters", func(t *testing.T) {
		title := "Unicode Title: 你好世界 🌍 العالم мир"
		input := "# " + title + "\n\nContent."
		result := extractor.ExtractTitle([]byte(input))
		if result != title {
			t.Errorf("ExtractTitle() = %q, want %q", result, title)
		}
	})

	t.Run("single character title", func(t *testing.T) {
		input := "# A\n\nContent."
		result := extractor.ExtractTitle([]byte(input))
		if result != "A" {
			t.Errorf("ExtractTitle() = %q, want %q", result, "A")
		}
	})
}

// BenchmarkMarkdownTitleExtractor_ExtractTitle benchmarks the title extraction performance
func BenchmarkMarkdownTitleExtractor_ExtractTitle(b *testing.B) {
	extractor := parser.NewMarkdownTitleExtractor()
	input := []byte("# This is a test title\n\nThis is some content that follows the title and contains multiple paragraphs.\n\nMore content here.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		extractor.ExtractTitle(input)
	}
}