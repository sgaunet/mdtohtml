// Package pdf converts Markdown to PDF by piping the existing Markdown→HTML
// pipeline straight into the folio HTML→PDF engine. No temporary HTML files
// are written: the rendered HTML lives only in memory between the two stages.
package pdf

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	folio "github.com/carlos7ags/folio/document"
	folioHTML "github.com/carlos7ags/folio/html"
	"github.com/carlos7ags/folio/layout"

	"github.com/sgaunet/mdtohtml/pkg/converter"
)

// pdfFontOverrideCSS rewrites the GitHub stylesheet's code-related rules to
// values folio's HTML→PDF engine renders correctly:
//
//   - font-family: folio collapses a CSS font-family list to its FIRST entry
//     and maps that single token to a standard PDF font; only "courier",
//     "monospace", "mono", or substrings of "times"/"serif" resolve. The
//     GitHub stylesheet lists "SFMono-Regular,Consolas,...,monospace" — the
//     first token wins, so inline code falls back to Helvetica. Putting the
//     generic "monospace" keyword first restores Courier in the PDF.
//   - background-color: folio's parseColor discards the alpha channel, so the
//     GitHub rule "background-color: rgba(27,31,35,.05)" on inline <code>
//     paints a near-black block over the text instead of a faint shade. We
//     replace it with the same opaque light gray used for fenced blocks so
//     the inline code chip stays readable.
const pdfFontOverrideCSS = `
body code,
body kbd,
body samp,
body tt {
  font-family: monospace;
  background-color: #f6f8fa;
}
body pre,
body pre code,
body pre > code {
  font-family: monospace;
}
`

// Converter renders Markdown directly to PDF. It satisfies converter.Converter
// so it can be plugged into the existing batch processor without any changes
// to the file-walking or path logic.
type Converter struct {
	htmlConv *converter.CompleteConverter
	pageSize folio.PageSize
	margins  layout.Margins
}

// New builds a PDF converter from the same options the HTML pipeline uses,
// plus PDF-specific options (page size, margins).
func New(opts converter.Options, pdfOpts Options) (*Converter, error) {
	ps, err := resolvePageSize(pdfOpts.PageSize)
	if err != nil {
		return nil, err
	}
	if opts.AdditionalCSS == "" {
		opts.AdditionalCSS = pdfFontOverrideCSS
	} else {
		opts.AdditionalCSS = opts.AdditionalCSS + "\n" + pdfFontOverrideCSS
	}
	return &Converter{
		htmlConv: converter.NewCompleteConverter(opts),
		pageSize: ps,
		margins: layout.Margins{
			Top:    pdfOpts.Margins.Top,
			Right:  pdfOpts.Margins.Right,
			Bottom: pdfOpts.Margins.Bottom,
			Left:   pdfOpts.Margins.Left,
		},
	}, nil
}

// Convert renders Markdown to PDF bytes. Relative image references will not
// resolve from this entry point because no input directory is known; use
// ConvertFile when assets need to load from disk.
func (c *Converter) Convert(input []byte) ([]byte, error) {
	htmlBytes, err := c.htmlConv.Convert(input)
	if err != nil {
		return nil, fmt.Errorf("markdown to HTML: %w", err)
	}

	doc, err := c.renderPDF(string(htmlBytes), "")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if _, err := doc.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("writing PDF: %w", err)
	}
	return buf.Bytes(), nil
}

// ConvertFile reads a Markdown file and writes the PDF to outputPath. The
// input file's directory is wired into folio's BasePath so relative image
// references (e.g. ./img/foo.png) resolve correctly.
func (c *Converter) ConvertFile(inputPath, outputPath string) error {
	input, err := os.ReadFile(inputPath)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied reading file '%s': %w", inputPath, err)
		}
		return fmt.Errorf("error reading file '%s': %w", inputPath, err)
	}

	htmlBytes, err := c.htmlConv.Convert(input)
	if err != nil {
		return fmt.Errorf("markdown to HTML: %w", err)
	}

	doc, err := c.renderPDF(string(htmlBytes), filepath.Dir(inputPath))
	if err != nil {
		return err
	}

	if err := doc.Save(outputPath); err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("permission denied writing file '%s': %w", outputPath, err)
		}
		return fmt.Errorf("error writing file '%s': %w", outputPath, err)
	}
	return nil
}

// renderPDF runs folio's HTML→PDF stage, applying any @page configuration
// found in the source HTML and forwarding the document title metadata.
func (c *Converter) renderPDF(htmlStr, basePath string) (*folio.Document, error) {
	result, err := folioHTML.ConvertFull(htmlStr, &folioHTML.Options{BasePath: basePath})
	if err != nil {
		return nil, fmt.Errorf("HTML to PDF: %w", err)
	}

	doc := folio.NewDocument(c.pageSize)
	doc.SetMargins(c.margins)
	if pc := result.PageConfig; pc != nil && pc.HasMargins {
		doc.SetMargins(layout.Margins{
			Top:    pc.MarginTop,
			Right:  pc.MarginRight,
			Bottom: pc.MarginBottom,
			Left:   pc.MarginLeft,
		})
	}
	if result.MarginBoxes != nil {
		doc.SetMarginBoxes(result.MarginBoxes)
	}
	if result.FirstMarginBoxes != nil {
		doc.SetFirstMarginBoxes(result.FirstMarginBoxes)
	}
	for _, e := range result.Elements {
		doc.Add(e)
	}
	if result.Metadata.Title != "" {
		doc.Info.Title = result.Metadata.Title
	}
	doc.SetAutoBookmarks(true)
	return doc, nil
}
