package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/processor"
)

var (
	outputDir  string
	pattern    string
	recursive  bool
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
}

func batchConvert(_ *cobra.Command, args []string) error {
	inputDir := args[0]

	if err := validateInputDir(inputDir); err != nil {
		return err
	}

	// Create converter with options
	options := converter.Options{
		SmartPunctuation: smartypants,
		LaTeXDashes:      latexdashes,
		Fractions:        fractions,
		SafeMode:         safeMode,
	}

	conv := converter.NewCompleteConverter(options)
	proc := processor.NewFileProcessor(conv)

	// Process directory
	processOptions := processor.ProcessOptions{
		OutputDir: outputDir,
		Pattern:   pattern,
		Recursive: recursive,
	}

	if err := proc.ProcessDirectory(inputDir, processOptions); err != nil {
		return fmt.Errorf("batch processing failed for directory '%s': %w", inputDir, err)
	}
	return nil
}