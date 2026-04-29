package cmd

import (
	"fmt"

	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/pdf"
)

// cssOptions groups the CSS-related flags for runConversion.
type cssOptions struct {
	source     string
	additional string
	noCSS      bool
}

func runConversion(
	inputFilePath, outputFilePath string,
	smartypants, latexdashes, fractions, safeMode bool,
	css cssOptions,
	format, pageSize, margin string,
) error {
	options := converter.Options{
		SmartPunctuation: smartypants,
		LaTeXDashes:      latexdashes,
		Fractions:        fractions,
		SafeMode:         safeMode,
		CSSSource:        css.source,
		AdditionalCSS:    css.additional,
		NoCSS:            css.noCSS,
	}

	conv, err := buildConverter(options, format, pageSize, margin)
	if err != nil {
		return err
	}

	if err := conv.ConvertFile(inputFilePath, outputFilePath); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}
	return nil
}

// buildConverter returns a converter.Converter implementation for the given
// output format. PDF wraps the HTML pipeline; HTML uses it directly.
func buildConverter(options converter.Options, format, pageSize, margin string) (converter.Converter, error) {
	if format == formatPDF {
		m, err := pdf.ParseMargin(margin)
		if err != nil {
			return nil, fmt.Errorf("invalid PDF options: %w", err)
		}
		pdfConv, err := pdf.New(options, pdf.Options{
			PageSize: pageSize,
			Margins:  pdf.Margins{Top: m, Right: m, Bottom: m, Left: m},
		})
		if err != nil {
			return nil, fmt.Errorf("invalid PDF options: %w", err)
		}
		return pdfConv, nil
	}
	return converter.NewCompleteConverter(options), nil
}
