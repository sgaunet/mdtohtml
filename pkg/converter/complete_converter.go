package converter

import (
	"fmt"
	"os"

	"github.com/sgaunet/mdtohtml/pkg/parser"
	"github.com/sgaunet/mdtohtml/pkg/template"
)

// CompleteConverter combines markdown conversion with title extraction and HTML templating.
type CompleteConverter struct {
	goldmarkConverter *GoldmarkConverter
	titleExtractor    parser.TitleExtractor
	htmlTemplate      template.HTMLTemplate
}

// NewCompleteConverter creates a new complete converter with all components.
func NewCompleteConverter(opts Options) *CompleteConverter {
	return &CompleteConverter{
		goldmarkConverter: NewGoldmarkConverter(opts),
		titleExtractor:    parser.NewMarkdownTitleExtractor(),
		htmlTemplate:      template.NewGitHubTemplate(),
	}
}

// NewCompleteConverterWithComponents creates a complete converter with custom components.
func NewCompleteConverterWithComponents(
	goldmarkConverter *GoldmarkConverter,
	titleExtractor parser.TitleExtractor,
	htmlTemplate template.HTMLTemplate,
) *CompleteConverter {
	return &CompleteConverter{
		goldmarkConverter: goldmarkConverter,
		titleExtractor:    titleExtractor,
		htmlTemplate:      htmlTemplate,
	}
}

// Convert transforms markdown content to complete HTML with title and CSS.
func (c *CompleteConverter) Convert(input []byte) ([]byte, error) {
	// Convert markdown to HTML
	htmlContent, err := c.goldmarkConverter.Convert(input)
	if err != nil {
		return nil, err
	}

	// Extract title
	title := c.titleExtractor.ExtractTitle(input)

	// Wrap in HTML document
	html := c.htmlTemplate.Wrap(string(htmlContent), title)

	// Inject CSS
	html = c.htmlTemplate.InjectCSS(html, "")

	return []byte(html), nil
}

// ConvertFile reads a markdown file and writes the complete HTML output.
func (c *CompleteConverter) ConvertFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", inputPath, err)
	}

	output, err := c.Convert(input)
	if err != nil {
		return err
	}

	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating file %s: %w", outputPath, err)
	}
	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", closeErr)
		}
	}()

	if _, err := out.Write(output); err != nil {
		return fmt.Errorf("error writing file %s: %w", outputPath, err)
	}

	return nil
}