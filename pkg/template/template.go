// Package template provides HTML templating functionality
package template

// HTMLTemplate defines the interface for HTML template operations.
type HTMLTemplate interface {
	// Wrap wraps HTML content with a complete HTML document structure
	Wrap(content, title string) string

	// InjectCSS injects CSS into an HTML document
	InjectCSS(html, css string) string
}