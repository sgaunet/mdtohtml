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
	// SmartPunctuation enables typographic substitutions:
	//   - Straight quotes ("x", 'x') become curly quotes (\u201cx\u201d, \u2018x\u2019)
	//   - Double hyphens (--) become en-dash (\u2013) or em-dash (\u2014) depending on LaTeXDashes
	//   - Triple hyphens (---) become em-dash (\u2014)
	//   - Three dots (...) become ellipsis (\u2026)
	SmartPunctuation bool

	// LaTeXDashes controls dash substitution style when SmartPunctuation is enabled.
	// When true (default): --- produces em-dash (\u2014), -- produces en-dash (\u2013).
	// When false: -- produces em-dash (\u2014), - between words produces en-dash (\u2013).
	LaTeXDashes bool

	// Fractions converts common fraction sequences to Unicode characters:
	//   - 1/2 becomes \u00bd, 1/4 becomes \u00bc, 3/4 becomes \u00be
	Fractions bool

	// SafeMode disables raw HTML pass-through to prevent XSS
	SafeMode bool

	// CSSSource is the full CSS text replacing the default embedded CSS.
	// Resolved by the CLI layer from --css-file or --css-url.
	CSSSource string

	// AdditionalCSS is extra CSS text appended to the default (or replaced) CSS.
	// Resolved by the CLI layer from --additional-css.
	AdditionalCSS string

	// NoCSS skips CSS injection entirely.
	NoCSS bool
}

// DefaultOptions returns the default converter options.
func DefaultOptions() Options {
	return Options{
		SmartPunctuation: true,
		LaTeXDashes:      true,
		Fractions:        true,
	}
}