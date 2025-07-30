// Package cmd contains the CLI commands for mdtohtml
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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
	batchCmd.Flags().BoolVar(&smartypants, "smartypants", true, "Apply smartypants-style substitutions")
	batchCmd.Flags().BoolVar(&latexdashes, "latexdashes", true, "Use LaTeX-style dash rules for smartypants")
	batchCmd.Flags().BoolVar(&fractions, "fractions", true, "Use improved fraction rules for smartypants")
}

func batchConvert(_ *cobra.Command, args []string) error {
	inputDir := args[0]
	
	// Validate input directory exists
	if err := validateInputDir(inputDir); err != nil {
		return err
	}

	// Create output directory if it doesn't exist
	if err := createOutputDir(); err != nil {
		return err
	}

	// Find files to convert
	files, err := findFiles(inputDir)
	if err != nil {
		return fmt.Errorf("error finding files: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No files matching pattern '%s' found in '%s'\n", pattern, inputDir)
		return nil
	}

	// Convert all files
	return convertFiles(files, inputDir)
}

func validateInputDir(inputDir string) error {
	_, err := os.Stat(inputDir)
	if os.IsNotExist(err) {
		return fmt.Errorf("%w: directory '%s'", ErrInputNotExist, inputDir)
	}
	return err //nolint:wrapcheck // os.Stat error is self-explanatory
}

func createOutputDir() error {
	const dirPerm = 0o755
	if err := os.MkdirAll(outputDir, dirPerm); err != nil {
		return fmt.Errorf("error creating output directory '%s': %w", outputDir, err)
	}
	return nil
}

func findFiles(inputDir string) ([]string, error) {
	if recursive {
		return findFilesRecursive(inputDir)
	}
	return filepath.Glob(filepath.Join(inputDir, pattern)) //nolint:wrapcheck // filepath.Glob error is self-explanatory
}

func findFilesRecursive(inputDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(inputDir, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			files = append(files, path)
		}
		return nil
	})
	return files, err //nolint:wrapcheck // filepath.WalkDir error is self-explanatory
}

func convertFiles(files []string, inputDir string) error {
	fmt.Printf("Converting %d files...\n", len(files))

	for _, file := range files {
		if err := convertFile(file, inputDir); err != nil {
			return err
		}
	}

	fmt.Printf("Successfully converted %d files to '%s'\n", len(files), outputDir)
	return nil
}

func convertFile(file, inputDir string) error {
	outputPath := getOutputPath(file, inputDir)
	
	// Create subdirectories if needed
	const dirPerm = 0o755
	if err := os.MkdirAll(filepath.Dir(outputPath), dirPerm); err != nil {
		return fmt.Errorf("error creating directory for '%s': %w", outputPath, err)
	}

	fmt.Printf("Converting %s -> %s\n", file, outputPath)
	if err := runConversion(file, outputPath, smartypants, latexdashes, fractions); err != nil {
		return fmt.Errorf("error converting '%s': %w", file, err)
	}
	
	return nil
}

func getOutputPath(file, inputDir string) string {
	relPath, err := filepath.Rel(inputDir, file)
	if err != nil {
		relPath = filepath.Base(file)
	}
	
	outputFile := strings.TrimSuffix(relPath, filepath.Ext(relPath)) + ".html"
	return filepath.Join(outputDir, outputFile)
}