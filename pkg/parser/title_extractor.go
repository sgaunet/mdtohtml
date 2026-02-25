package parser

import "strings"

// MarkdownTitleExtractor implements TitleExtractor for markdown content.
type MarkdownTitleExtractor struct{}

// NewMarkdownTitleExtractor creates a new markdown title extractor.
func NewMarkdownTitleExtractor() *MarkdownTitleExtractor {
	return &MarkdownTitleExtractor{}
}

// ExtractTitle extracts the title from markdown content.
// It looks for either a # prefix header or an underlined header (===).
func (e *MarkdownTitleExtractor) ExtractTitle(content []byte) string {
	pos := skipBlankLines(content)
	if pos >= len(content) {
		return ""
	}

	line, nextPos := readLine(content, pos)
	if len(line) == 0 {
		return ""
	}

	if title, ok := parsePrefixHeader(line); ok {
		return title
	}

	if hasEqualUnderline(content, nextPos) {
		return strings.TrimSpace(string(line))
	}

	return ""
}

// skipBlankLines returns the position of the first non-blank-line character.
func skipBlankLines(content []byte) int {
	i := 0
	for i < len(content) && (content[i] == '\n' || content[i] == '\r') {
		i++
	}
	return i
}

// readLine extracts the line starting at pos and returns the position after its line ending.
func readLine(content []byte, pos int) ([]byte, int) {
	start := pos
	for pos < len(content) && content[pos] != '\n' && content[pos] != '\r' {
		pos++
	}
	line := content[start:pos]

	if pos < len(content) && content[pos] == '\r' && pos+1 < len(content) && content[pos+1] == '\n' {
		pos += 2
	} else if pos < len(content) {
		pos++
	}
	return line, pos
}

// parsePrefixHeader checks if a line is a level-1 prefix header (# Title).
func parsePrefixHeader(line []byte) (string, bool) {
	if len(line) >= 3 && line[0] == '#' && (line[1] == ' ' || line[1] == '\t') {
		return strings.TrimSpace(string(line[2:])), true
	}
	return "", false
}

// hasEqualUnderline checks if content at pos is a valid === underline (at least 3 equal signs).
func hasEqualUnderline(content []byte, pos int) bool {
	const minEqualSigns = 3

	if pos >= len(content) || content[pos] != '=' {
		return false
	}

	start := pos
	for pos < len(content) && content[pos] == '=' {
		pos++
	}

	return (pos-start) >= minEqualSigns && isBlankUntilLineEnd(content, pos)
}

// isBlankUntilLineEnd returns true if content from pos to the next line ending contains only whitespace.
func isBlankUntilLineEnd(content []byte, pos int) bool {
	for pos < len(content) && (content[pos] == ' ' || content[pos] == '\t') {
		pos++
	}
	return pos >= len(content) || content[pos] == '\n' || content[pos] == '\r'
}
