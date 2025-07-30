package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version is set at build time.
	Version     = "development"
	smartypants bool
	latexdashes bool
	fractions   bool
)

var rootCmd = &cobra.Command{
	Use:   "mdtohtml [input.md] [output.html]",
	Short: "Convert Markdown files to HTML with GitHub-style CSS",
	Long: `mdtohtml is a command-line tool that converts Markdown files to HTML with GitHub-style CSS.
It supports GitHub Flavored Markdown, definition lists, footnotes, and typographic enhancements.`,
	Args: cobra.ExactArgs(2), //nolint:mnd // requires exactly 2 args: input and output
	RunE: convert,
	Example: `  mdtohtml README.md README.html
  mdtohtml -smartypants=false input.md output.html`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&smartypants, "smartypants", true, "Apply smartypants-style substitutions")
	rootCmd.Flags().BoolVar(&latexdashes, "latexdashes", true, "Use LaTeX-style dash rules for smartypants")
	rootCmd.Flags().BoolVar(&fractions, "fractions", true, "Use improved fraction rules for smartypants")
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(`{{.Version}}
`)
}

func convert(_ *cobra.Command, args []string) error {
	inputFilePath := args[0]
	outputFilePath := args[1]

	// Validate input file exists
	_, err := os.Stat(inputFilePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: file '%s'", ErrInputNotExist, inputFilePath)
	}

	return runConversion(inputFilePath, outputFilePath, smartypants, latexdashes, fractions)
}