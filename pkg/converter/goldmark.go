package converter

import (
	"bytes"
	"fmt"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// GoldmarkConverter implements the Converter interface using the Goldmark library.
type GoldmarkConverter struct {
	md      goldmark.Markdown
	options Options
}

// NewGoldmarkConverter creates a new converter with the given options.
func NewGoldmarkConverter(opts Options) *GoldmarkConverter {
	extensions := []goldmark.Extender{
		extension.GFM,
		extension.DefinitionList,
		extension.Footnote,
	}

	if opts.SmartPunctuation || opts.LaTeXDashes || opts.Fractions {
		extensions = append(extensions, extension.Typographer)
	}

	md := goldmark.New(
		goldmark.WithExtensions(extensions...),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)

	return &GoldmarkConverter{
		md:      md,
		options: opts,
	}
}

// Convert transforms markdown content to HTML.
func (c *GoldmarkConverter) Convert(input []byte) ([]byte, error) {
	var buf bytes.Buffer
	if err := c.md.Convert(input, &buf); err != nil {
		return nil, fmt.Errorf("error converting markdown: %w", err)
	}
	return buf.Bytes(), nil
}

// ConvertFile reads a markdown file and writes the HTML output.
func (c *GoldmarkConverter) ConvertFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", inputPath, err)
	}

	output, err := c.Convert(input)
	if err != nil {
		return err
	}

	const defaultFileMode = 0644
	if err := os.WriteFile(outputPath, output, defaultFileMode); err != nil {
		return fmt.Errorf("error writing file %s: %w", outputPath, err)
	}

	return nil
}