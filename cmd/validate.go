package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
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
		"Apply smartypants-style substitutions during validation")
	validateCmd.Flags().BoolVar(&latexdashes, "latexdashes", true,
		"Use LaTeX-style dash rules for smartypants during validation")
	validateCmd.Flags().BoolVar(&fractions, "fractions", true,
		"Use improved fraction rules for smartypants during validation")
}

func validateMarkdown(_ *cobra.Command, args []string) error {
	inputFilePath := args[0]

	// Validate input file exists
	_, err := os.Stat(inputFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: file '%s'", ErrInputNotExist, inputFilePath)
	}

	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("error reading from %s: %w", inputFilePath, err)
	}

	md := createGoldmarkProcessor(smartypants, fractions, latexdashes)

	// Try to process the markdown
	_, err = processMarkdown(input, md)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	fmt.Printf("✓ %s is valid Markdown\n", inputFilePath)
	return nil
}