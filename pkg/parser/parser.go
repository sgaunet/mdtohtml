// Package parser provides markdown parsing utilities
package parser

// TitleExtractor defines the interface for extracting titles from markdown content.
type TitleExtractor interface {
	// ExtractTitle extracts the title from markdown content
	ExtractTitle(content []byte) string
}