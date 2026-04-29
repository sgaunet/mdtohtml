package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/processor"
)

var (
	outputDir string
	pattern   string
	recursive bool
)

var batchCmd = &cobra.Command{
	Use:   "batch [input-dir]",
	Short: "Convert multiple Markdown files to HTML",
	Long: `Convert multiple Markdown files to HTML with GitHub-style CSS.
By default, processes all *.md files in the specified directory.`,
	Args: cobra.ExactArgs(1),
	RunE: batchConvert,
	Example: `  mdtohtml batch ./docs --out-dir ./html
  mdtohtml batch ./docs --pattern "*.markdown" --out-dir ./public
  mdtohtml batch ./docs --recursive --out-dir ./output`,
}

func init() {
	rootCmd.AddCommand(batchCmd)

	batchCmd.Flags().StringVarP(&outputDir, "out-dir", "o", ".", "Output directory for HTML files")
	batchCmd.Flags().StringVarP(&pattern, "pattern", "p", "*.md", "File pattern to match")
	batchCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Process directories recursively")
	batchCmd.Flags().BoolVar(&smartypants, "smartypants", true,
		`Convert quotes to curly quotes, -- to en/em-dash, ... to ellipsis`)
	batchCmd.Flags().BoolVar(&latexdashes, "latexdashes", true,
		`LaTeX-style dashes: --- for em-dash, -- for en-dash (requires --smartypants)`)
	batchCmd.Flags().BoolVar(&fractions, "fractions", true,
		`Convert fractions: 1/2 to ½, 1/4 to ¼, 3/4 to ¾`)
	batchCmd.Flags().BoolVar(&safeMode, "safe-mode", false, "Disable raw HTML pass-through to prevent XSS")
	batchCmd.Flags().StringVar(&cssFile, "css-file", "", "Path to a CSS file to use instead of the default GitHub CSS")
	batchCmd.Flags().StringVar(&cssURL, "css-url", "", "URL to fetch CSS from instead of the default GitHub CSS")
	batchCmd.Flags().StringVar(&additionalCSSFile, "additional-css", "", "Path to a CSS file to append to the default CSS")
	batchCmd.Flags().BoolVar(&noCSS, "no-css", false, "Disable CSS injection entirely")
	batchCmd.Flags().StringVar(&outputFormat, "format", formatHTML,
		`Output format: "html" or "pdf"`)
	batchCmd.Flags().StringVar(&pageSize, "page-size", "A4",
		`PDF page size when --format=pdf: A4, Letter, Legal, A3, A5, Tabloid`)
	batchCmd.Flags().StringVar(&marginFlag, "margin", defaultMarginFlag,
		`PDF page margin (units: pt, in, cm, mm; bare numbers = pt)`)
}

func batchConvert(_ *cobra.Command, args []string) error {
	inputDir := args[0]

	if err := validateInputDir(inputDir); err != nil {
		return err
	}

	source, additional, err := resolveCSSOptions(
		cssFile, cssURL, additionalCSSFile, noCSS,
	)
	if err != nil {
		return err
	}

	// Create converter with options
	options := converter.Options{
		SmartPunctuation: smartypants,
		LaTeXDashes:      latexdashes,
		Fractions:        fractions,
		SafeMode:         safeMode,
		CSSSource:        source,
		AdditionalCSS:    additional,
		NoCSS:            noCSS,
	}

	format, err := resolveFormat(outputFormat, "")
	if err != nil {
		return err
	}

	conv, err := buildConverter(options, format, pageSize, marginFlag)
	if err != nil {
		return err
	}
	proc := processor.NewFileProcessor(conv)

	processOptions := processor.ProcessOptions{
		OutputDir: outputDir,
		Pattern:   pattern,
		Recursive: recursive,
		OutputExt: extForFormat(format),
	}

	if err := proc.ProcessDirectory(inputDir, processOptions); err != nil {
		return fmt.Errorf("batch processing failed for directory '%s': %w", inputDir, err)
	}
	return nil
}
