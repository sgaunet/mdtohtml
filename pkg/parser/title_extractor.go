package parser

import "strings"

// MarkdownTitleExtractor implements TitleExtractor for markdown content.
type MarkdownTitleExtractor struct{}

// NewMarkdownTitleExtractor creates a new markdown title extractor.
func NewMarkdownTitleExtractor() *MarkdownTitleExtractor {
	return &MarkdownTitleExtractor{}
}

// ExtractTitle extracts the title from markdown content.
// It looks for either a # prefix header or an underlined header.
func (e *MarkdownTitleExtractor) ExtractTitle(content []byte) string {
	if len(content) == 0 {
		return ""
	}

	i := e.skipBlankLines(content, 0)
	if i >= len(content) {
		return ""
	}

	line1, nextPos := e.extractFirstLine(content, i)
	if len(line1) == 0 {
		return ""
	}

	// Check for a prefix header (# Title)
	if title := e.checkPrefixHeader(line1); title != "" {
		return title
	}

	// Check for an underlined header
	if e.checkUnderlinedHeader(content, nextPos) {
		return strings.TrimSpace(string(line1))
	}

	return ""
}

func (e *MarkdownTitleExtractor) skipBlankLines(input []byte, start int) int {
	i := start
	for i < len(input) && (input[i] == '\n' || input[i] == '\r') {
		i++
	}
	return i
}

func (e *MarkdownTitleExtractor) extractFirstLine(input []byte, start int) ([]byte, int) {
	i := e.skipCarriageReturn(input, start)
	lineStart := i
	i = e.findLineEnd(input, i)
	line := input[lineStart:i]
	i = e.skipLineEnding(input, i)
	return line, i
}

func (e *MarkdownTitleExtractor) skipCarriageReturn(input []byte, pos int) int {
	if pos < len(input) && input[pos] == '\r' && pos+1 < len(input) && input[pos+1] == '\n' {
		return pos + 1
	}
	return pos
}

func (e *MarkdownTitleExtractor) findLineEnd(input []byte, start int) int {
	i := start
	for i < len(input) && input[i] != '\n' && input[i] != '\r' {
		i++
	}
	return i
}

func (e *MarkdownTitleExtractor) skipLineEnding(input []byte, pos int) int {
	const crlfLen = 2
	if pos < len(input) && input[pos] == '\r' && pos+1 < len(input) && input[pos+1] == '\n' {
		return pos + crlfLen
	}
	if pos < len(input) {
		return pos + 1
	}
	return pos
}

func (e *MarkdownTitleExtractor) checkPrefixHeader(line []byte) string {
	if len(line) >= 3 && line[0] == '#' && (line[1] == ' ' || line[1] == '\t') {
		return strings.TrimSpace(string(line[2:]))
	}
	return ""
}

func (e *MarkdownTitleExtractor) checkUnderlinedHeader(input []byte, pos int) bool {
	if pos >= len(input) || input[pos] != '=' {
		return false
	}

	equalCount := e.countEqualSigns(input, pos)
	endPos := e.skipToLineEnd(input, pos+equalCount)
	
	const minEqualSigns = 3
	return e.isAtLineEnd(input, endPos) && equalCount >= minEqualSigns
}

func (e *MarkdownTitleExtractor) countEqualSigns(input []byte, start int) int {
	count := 0
	for i := start; i < len(input) && input[i] == '='; i++ {
		count++
	}
	return count
}

func (e *MarkdownTitleExtractor) skipToLineEnd(input []byte, start int) int {
	i := start
	for i < len(input) && (input[i] == ' ' || input[i] == '\t') {
		i++
	}
	return i
}

func (e *MarkdownTitleExtractor) isAtLineEnd(input []byte, pos int) bool {
	return pos >= len(input) || input[pos] == '\n' || input[pos] == '\r'
}