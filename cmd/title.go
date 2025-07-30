package cmd

import "strings"

const defaultTitle = ""

// getTitle tries to guess the title from the input buffer.
// It checks if it starts with an <h1> element and uses that.
func getTitle(input []byte) string {
	if len(input) == 0 {
		return defaultTitle
	}

	i := skipBlankLines(input, 0)
	if i >= len(input) {
		return defaultTitle
	}

	line1, nextPos := extractFirstLine(input, i)
	if len(line1) == 0 {
		return defaultTitle
	}

	// check for a prefix header
	if title := checkPrefixHeader(line1); title != "" {
		return title
	}

	// check for an underlined header
	if checkUnderlinedHeader(input, nextPos) {
		return strings.TrimSpace(string(line1))
	}

	return defaultTitle
}

func skipBlankLines(input []byte, start int) int {
	i := start
	for i < len(input) && (input[i] == '\n' || input[i] == '\r') {
		i++
	}

	return i
}

func extractFirstLine(input []byte, start int) ([]byte, int) {
	i := skipCarriageReturn(input, start)
	lineStart := i
	i = findLineEnd(input, i)
	line := input[lineStart:i]
	i = skipLineEnding(input, i)

	return line, i
}

func skipCarriageReturn(input []byte, pos int) int {
	if pos < len(input) && input[pos] == '\r' && pos+1 < len(input) && input[pos+1] == '\n' {
		return pos + 1
	}

	return pos
}

func findLineEnd(input []byte, start int) int {
	i := start
	for i < len(input) && input[i] != '\n' && input[i] != '\r' {
		i++
	}
	return i
}

func skipLineEnding(input []byte, pos int) int {
	const crlfLen = 2
	if pos < len(input) && input[pos] == '\r' && pos+1 < len(input) && input[pos+1] == '\n' {
		return pos + crlfLen
	}
	return pos + 1
}

func checkPrefixHeader(line []byte) string {
	const minHeaderLen = 3
	if len(line) >= minHeaderLen && line[0] == '#' && (line[1] == ' ' || line[1] == '\t') {
		return strings.TrimSpace(string(line[2:]))
	}
	return ""
}

func checkUnderlinedHeader(input []byte, pos int) bool {
	i := pos
	if i >= len(input) || input[i] != '=' {
		return false
	}
	
	for i < len(input) && input[i] == '=' {
		i++
	}
	
	for i < len(input) && (input[i] == ' ' || input[i] == '\t') {
		i++
	}

	return i < len(input) && (input[i] == '\n' || input[i] == '\r')
}