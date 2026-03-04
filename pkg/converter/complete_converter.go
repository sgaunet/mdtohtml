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
	noCSS             bool
}

// NewCompleteConverter creates a new complete converter with all components.
func NewCompleteConverter(opts Options) *CompleteConverter {
	var tmpl template.HTMLTemplate
	switch {
	case opts.CSSSource != "" && opts.AdditionalCSS != "":
		tmpl = template.NewGitHubTemplateWithCSS(opts.CSSSource + "\n" + opts.AdditionalCSS)
	case opts.CSSSource != "":
		tmpl = template.NewGitHubTemplateWithCSS(opts.CSSSource)
	case opts.AdditionalCSS != "":
		tmpl = template.NewGitHubTemplateWithAdditionalCSS(opts.AdditionalCSS)
	default:
		tmpl = template.NewGitHubTemplate()
	}

	return &CompleteConverter{
		goldmarkConverter: NewGoldmarkConverter(opts),
		titleExtractor:    parser.NewMarkdownTitleExtractor(),
		htmlTemplate:      tmpl,
		noCSS:             opts.NoCSS,
	}
}

// NewCompleteConverterWithComponents creates a complete converter with custom components,
// allowing callers to replace the default title extractor or HTML template.
// This is useful for customizing the output format without modifying the core
// markdown conversion logic.
//
// Parameters:
//   - goldmarkConverter: the Goldmark-based markdown converter (use [NewGoldmarkConverter])
//   - titleExtractor: extracts a document title from markdown (implements [parser.TitleExtractor])
//   - htmlTemplate: wraps converted HTML with document structure and CSS (implements [template.HTMLTemplate])
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

	// Inject CSS unless disabled
	if !c.noCSS {
		html = c.htmlTemplate.InjectCSS(html, "")
	}

	return []byte(html), nil
}

// ConvertFile reads a markdown file and writes the complete HTML output.
func (c *CompleteConverter) ConvertFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied reading file '%s': %w", inputPath, err)
		}
		return fmt.Errorf("error reading file '%s': %w", inputPath, err)
	}

	output, err := c.Convert(input)
	if err != nil {
		return err
	}

	const defaultFileMode = 0644
	if err := os.WriteFile(outputPath, output, defaultFileMode); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied writing file '%s': %w", outputPath, err)
		}
		return fmt.Errorf("error writing file '%s': %w", outputPath, err)
	}

	return nil
}