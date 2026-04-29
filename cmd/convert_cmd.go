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
	convertCmd.Flags().StringVar(&cssFile, "css-file", "", "Path to a CSS file to use instead of the default GitHub CSS")
	convertCmd.Flags().StringVar(&cssURL, "css-url", "", "URL to fetch CSS from instead of the default GitHub CSS")
	convertCmd.Flags().StringVar(&additionalCSSFile, "additional-css", "",
		"Path to a CSS file to append to the default CSS")
	convertCmd.Flags().BoolVar(&noCSS, "no-css", false, "Disable CSS injection entirely")
	convertCmd.Flags().StringVar(&outputFormat, "format", "",
		`Output format: "html" or "pdf" (default: auto-detect from output file extension)`)
	convertCmd.Flags().StringVar(&pageSize, "page-size", "A4",
		`PDF page size when --format=pdf: A4, Letter, Legal, A3, A5, Tabloid`)
	convertCmd.Flags().StringVar(&marginFlag, "margin", defaultMarginFlag,
		`PDF page margin (units: pt, in, cm, mm; bare numbers = pt)`)
}
