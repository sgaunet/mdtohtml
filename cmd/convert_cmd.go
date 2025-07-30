package cmd

import (
	"github.com/spf13/cobra"
)

var convertCmd = &cobra.Command{
	Use:   "convert [input.md] [output.html]",
	Short: "Convert a single Markdown file to HTML",
	Long: `Convert a single Markdown file to HTML with GitHub-style CSS.
This is the default behavior when no subcommand is specified.`,
	Args: cobra.ExactArgs(2), //nolint:mnd // requires exactly 2 args: input and output
	RunE: convert,
	Example: `  mdtohtml convert README.md README.html
  mdtohtml convert -smartypants=false input.md output.html`,
}

func init() {
	rootCmd.AddCommand(convertCmd)
	
	// Add the same flags to convert command
	convertCmd.Flags().BoolVar(&smartypants, "smartypants", true, "Apply smartypants-style substitutions")
	convertCmd.Flags().BoolVar(&latexdashes, "latexdashes", true, "Use LaTeX-style dash rules for smartypants")
	convertCmd.Flags().BoolVar(&fractions, "fractions", true, "Use improved fraction rules for smartypants")
}