// Package heading provides markdown title/heading extraction utilities
package heading

// TitleExtractor defines the interface for extracting titles from markdown content.
type TitleExtractor interface {
	// ExtractTitle extracts the title from markdown content
	ExtractTitle(content []byte) string
}