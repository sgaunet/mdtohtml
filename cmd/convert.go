package cmd

import (
	"fmt"

	"github.com/sgaunet/mdtohtml/pkg/converter"
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
) error {
	// Create converter with options
	options := converter.Options{
		SmartPunctuation: smartypants,
		LaTeXDashes:      latexdashes,
		Fractions:        fractions,
		SafeMode:         safeMode,
		CSSSource:        css.source,
		AdditionalCSS:    css.additional,
		NoCSS:            css.noCSS,
	}

	conv := converter.NewCompleteConverter(options)

	// Convert file
	if err := conv.ConvertFile(inputFilePath, outputFilePath); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	return nil
}
