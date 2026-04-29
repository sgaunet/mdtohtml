package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version holds the application version string, injected at build time via ldflags.
	Version           = "development"
	smartypants       bool
	latexdashes       bool
	fractions         bool
	safeMode          bool
	cssFile           string
	cssURL            string
	additionalCSSFile string
	noCSS             bool
	outputFormat      string // "", "html", "pdf"; empty = auto-detect from extension
	pageSize          string // PDF page size, e.g. "A4", "Letter"
	marginFlag        string // PDF margin, e.g. "1.25in", "90", "2.5cm"
)

const defaultMarginFlag = "1.25in"

// Recognised output formats.
const (
	formatHTML = "html"
	formatPDF  = "pdf"
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
	rootCmd.Flags().StringVar(&cssFile, "css-file", "", "Path to a CSS file to use instead of the default GitHub CSS")
	rootCmd.Flags().StringVar(&cssURL, "css-url", "", "URL to fetch CSS from instead of the default GitHub CSS")
	rootCmd.Flags().StringVar(&additionalCSSFile, "additional-css", "", "Path to a CSS file to append to the default CSS")
	rootCmd.Flags().BoolVar(&noCSS, "no-css", false, "Disable CSS injection entirely")
	rootCmd.Flags().StringVar(&outputFormat, "format", "",
		`Output format: "html" or "pdf" (default: auto-detect from output file extension)`)
	rootCmd.Flags().StringVar(&pageSize, "page-size", "A4",
		`PDF page size when --format=pdf: A4, Letter, Legal, A3, A5, Tabloid`)
	rootCmd.Flags().StringVar(&marginFlag, "margin", defaultMarginFlag,
		`PDF page margin (units: pt, in, cm, mm; bare numbers = pt)`)
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

	source, additional, err := resolveCSSOptions(
		cssFile, cssURL, additionalCSSFile, noCSS,
	)
	if err != nil {
		return err
	}

	css := cssOptions{
		source: source, additional: additional, noCSS: noCSS,
	}
	format, err := resolveFormat(outputFormat, outputFilePath)
	if err != nil {
		return err
	}
	return runConversion(
		inputFilePath, outputFilePath,
		smartypants, latexdashes, fractions, safeMode, css,
		format, pageSize, marginFlag,
	)
}