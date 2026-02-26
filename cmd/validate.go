package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/validator"
)

var validateCmd = &cobra.Command{
	Use:   "validate [input.md]",
	Short: "Validate Markdown syntax without converting",
	Long: `Validate Markdown syntax without generating output.
This checks if the file can be parsed successfully by the Goldmark processor.`,
	Args: cobra.ExactArgs(1),
	RunE: validateMarkdown,
	Example: `  mdtohtml validate README.md
  mdtohtml validate document.md`,
}

func init() {
	rootCmd.AddCommand(validateCmd)
	
	validateCmd.Flags().BoolVar(&smartypants, "smartypants", true,
		`Convert quotes to curly quotes, -- to en/em-dash, ... to ellipsis`)
	validateCmd.Flags().BoolVar(&latexdashes, "latexdashes", true,
		`LaTeX-style dashes: --- for em-dash, -- for en-dash (requires --smartypants)`)
	validateCmd.Flags().BoolVar(&fractions, "fractions", true,
		`Convert fractions: 1/2 to ½, 1/4 to ¼, 3/4 to ¾`)
}

func validateMarkdown(_ *cobra.Command, args []string) error {
	inputFilePath := args[0]

	if err := validateInputFile(inputFilePath); err != nil {
		return err
	}

	// Create converter with options
	options := converter.Options{
		SmartPunctuation: smartypants,
		LaTeXDashes:      latexdashes,
		Fractions:        fractions,
	}

	conv := converter.NewGoldmarkConverter(options)
	val := validator.NewGoldmarkValidator(conv)

	// Validate the file
	if err := val.ValidateFile(inputFilePath); err != nil {
		return fmt.Errorf("validation failed for %s: %w", inputFilePath, err)
	}

	fmt.Printf("✓ %s is valid Markdown\n", inputFilePath)
	return nil
}