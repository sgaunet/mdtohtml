package cmd

import (
	"fmt"

	"github.com/sgaunet/mdtohtml/pkg/converter"
)

func runConversion(inputFilePath, outputFilePath string, smartypants, latexdashes, fractions bool) error {
	// Create converter with options
	options := converter.Options{
		SmartPunctuation: smartypants,
		LaTeXDashes:      latexdashes,
		Fractions:        fractions,
	}

	conv := converter.NewCompleteConverter(options)

	// Convert file
	if err := conv.ConvertFile(inputFilePath, outputFilePath); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	return nil
}