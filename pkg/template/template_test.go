package template_test

import (
	"strings"
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/template"
)

// TestGitHubTemplate_Wrap tests HTML document wrapping functionality
func TestGitHubTemplate_Wrap(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		title    string
		contains []string
		notContains []string
	}{
		{
			name:    "basic wrapping",
			content: "<p>Hello World</p>",
			title:   "Test Title",
			contains: []string{
				"<!DOCTYPE html>",
				"<html>",
				"<head>",
				"<meta charset=\"UTF-8\">",
				"<title>Test Title</title>",
				"</head>",
				"<body>",
				"<p>Hello World</p>",
				"</body>",
				"</html>",
			},
		},
		{
			name:    "empty title",
			content: "<h1>Content</h1>",
			title:   "",
			contains: []string{
				"<title></title>",
				"<h1>Content</h1>",
			},
		},
		{
			name:    "title with special characters",
			content: "<p>Content</p>",
			title:   "Title with <special> & \"characters\"",
			contains: []string{
				"<title>Title with <special> & \"characters\"</title>",
			},
		},
		{
			name:    "complex content",
			content: "<h1>Header</h1><p>Paragraph with <strong>bold</strong> text.</p><ul><li>Item 1</li><li>Item 2</li></ul>",
			title:   "Complex Document",
			contains: []string{
				"<title>Complex Document</title>",
				"<h1>Header</h1>",
				"<strong>bold</strong>",
				"<ul><li>Item 1</li>",
			},
		},
	}

	tmpl := template.NewGitHubTemplate()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tmpl.Wrap(tt.content, tt.title)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Wrap() result does not contain %q\nGot: %s", expected, result)
				}
			}

			for _, notExpected := range tt.notContains {
				if strings.Contains(result, notExpected) {
					t.Errorf("Wrap() result should not contain %q\nGot: %s", notExpected, result)
				}
			}
		})
	}
}

// TestGitHubTemplate_InjectCSS tests CSS injection functionality
func TestGitHubTemplate_InjectCSS(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		css      string
		contains []string
	}{
		{
			name: "inject custom CSS",
			html: "<!DOCTYPE html>\n<html>\n<head>\n<title>Test</title>\n</head>\n<body></body>\n</html>",
			css:  "body { margin: 0; }",
			contains: []string{
				"<style>\nbody { margin: 0; }\n</style>\n</head>",
				"<title>Test</title>",
			},
		},
		{
			name: "inject default CSS",
			html: "<!DOCTYPE html>\n<html>\n<head>\n<title>Test</title>\n</head>\n<body></body>\n</html>",
			css:  "", // empty CSS should use default
			contains: []string{
				"<style>",
				"</style>\n</head>",
				".octicon", // part of GitHub CSS
			},
		},
		{
			name: "multiple head sections",
			html: "<!DOCTYPE html>\n<html>\n<head>\n<title>Test</title>\n</head>\n<body>\n<head>fake</head>\n</body>\n</html>",
			css:  "p { color: red; }",
			contains: []string{
				"<style>\np { color: red; }\n</style>\n</head>",
				"<head>fake</head>", // should not affect fake head in body
			},
		},
	}

	tmpl := template.NewGitHubTemplate()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tmpl.InjectCSS(tt.html, tt.css)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("InjectCSS() result does not contain %q\nGot: %s", expected, result)
				}
			}
		})
	}
}

// TestGitHubTemplate_WithCustomCSS tests custom CSS functionality
func TestGitHubTemplate_WithCustomCSS(t *testing.T) {
	customCSS := "body { background: red; }"
	tmpl := template.NewGitHubTemplateWithCSS(customCSS)

	html := "<!DOCTYPE html>\n<html>\n<head>\n<title>Test</title>\n</head>\n<body></body>\n</html>"
	result := tmpl.InjectCSS(html, "")

	if !strings.Contains(result, customCSS) {
		t.Errorf("Custom CSS template should inject custom CSS by default")
	}

	// Test overriding custom CSS
	overrideCSS := "body { background: blue; }"
	result = tmpl.InjectCSS(html, overrideCSS)

	if !strings.Contains(result, overrideCSS) {
		t.Errorf("Should be able to override custom CSS")
	}

	if strings.Contains(result, customCSS) {
		t.Errorf("Should not contain original custom CSS when overridden")
	}
}

// TestGitHubTemplate_Integration tests the complete workflow
func TestGitHubTemplate_Integration(t *testing.T) {
	tmpl := template.NewGitHubTemplate()

	content := "<h1>My Document</h1><p>This is a test document.</p>"
	title := "Integration Test"

	// First wrap the content
	wrapped := tmpl.Wrap(content, title)

	// Then inject CSS
	final := tmpl.InjectCSS(wrapped, "")

	expectedElements := []string{
		"<!DOCTYPE html>",
		"<title>Integration Test</title>",
		"<h1>My Document</h1>",
		"<p>This is a test document.</p>",
		"<style>",
		".octicon",
		"</style>",
	}

	for _, expected := range expectedElements {
		if !strings.Contains(final, expected) {
			t.Errorf("Integration test result does not contain %q", expected)
		}
	}

	// Verify proper structure
	if !strings.Contains(final, "<style>") || !strings.Contains(final, "</style>") {
		t.Error("CSS should be properly wrapped in <style> tags")
	}

	// Verify CSS comes before </head>
	styleIndex := strings.Index(final, "<style>")
	headEndIndex := strings.Index(final, "</head>")
	if styleIndex == -1 || headEndIndex == -1 || styleIndex >= headEndIndex {
		t.Error("CSS should be injected before </head> tag")
	}
}