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
	convertCmd.Flags().BoolVar(&smartypants, "smartypants", true,
		`Convert quotes to curly quotes, -- to en/em-dash, ... to ellipsis`)
	convertCmd.Flags().BoolVar(&latexdashes, "latexdashes", true,
		`LaTeX-style dashes: --- for em-dash, -- for en-dash (requires --smartypants)`)
	convertCmd.Flags().BoolVar(&fractions, "fractions", true,
		`Convert fractions: 1/2 to ½, 1/4 to ¼, 3/4 to ¾`)
}