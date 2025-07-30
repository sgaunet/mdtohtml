// Package converter provides interfaces and implementations for converting markdown to HTML
package converter

// Converter defines the interface for markdown to HTML conversion.
type Converter interface {
	// Convert transforms markdown content to HTML
	Convert(input []byte) ([]byte, error)

	// ConvertFile reads a markdown file and writes the HTML output
	ConvertFile(inputPath, outputPath string) error
}

// Options configures the converter behavior.
type Options struct {
	// SmartPunctuation enables smart quotes, dashes, and ellipses
	SmartPunctuation bool

	// LaTeXDashes uses LaTeX-style dash rules
	LaTeXDashes bool

	// Fractions converts fractions like 1/2 to ½
	Fractions bool
}

// DefaultOptions returns the default converter options.
func DefaultOptions() Options {
	return Options{
		SmartPunctuation: true,
		LaTeXDashes:      true,
		Fractions:        true,
	}
}