package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version holds the application version string, injected at build time via ldflags.
	Version = "development"
	smartypants bool
	latexdashes bool
	fractions   bool
	safeMode    bool
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

// Execute adds all child commands to the root command and sets flags appropriately.
// It is called by main.main() and exits the process with status 1 on error.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolVar(&smartypants, "smartypants", true,
		`Convert quotes to curly quotes, -- to en/em-dash, ... to ellipsis`)
	rootCmd.Flags().BoolVar(&latexdashes, "latexdashes", true,
		`LaTeX-style dashes: --- for em-dash, -- for en-dash (requires --smartypants)`)
	rootCmd.Flags().BoolVar(&fractions, "fractions", true,
		`Convert fractions: 1/2 to ½, 1/4 to ¼, 3/4 to ¾`)
	rootCmd.Flags().BoolVar(&safeMode, "safe-mode", false, "Disable raw HTML pass-through to prevent XSS")
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(`{{.Version}}
`)
}

func convert(_ *cobra.Command, args []string) error {
	inputFilePath := args[0]
	outputFilePath := args[1]

	if err := validateInputFile(inputFilePath); err != nil {
		return err
	}

	return runConversion(inputFilePath, outputFilePath, smartypants, latexdashes, fractions, safeMode)
}