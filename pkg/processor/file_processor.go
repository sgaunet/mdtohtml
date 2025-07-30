package processor

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sgaunet/mdtohtml/pkg/converter"
)

// FileProcessor implements BatchProcessor for file system operations.
type FileProcessor struct {
	converter converter.Converter
}

// NewFileProcessor creates a new file processor with the given converter.
func NewFileProcessor(conv converter.Converter) *FileProcessor {
	return &FileProcessor{
		converter: conv,
	}
}

// ProcessDirectory processes all matching files in a directory.
func (p *FileProcessor) ProcessDirectory(dir string, options ProcessOptions) error {
	// Validate input directory
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return fmt.Errorf("%w: %s", ErrDirectoryNotExist, dir)
	}

	// Create output directory
	const defaultDirMode = 0755
	if err := os.MkdirAll(options.OutputDir, defaultDirMode); err != nil {
		return fmt.Errorf("error creating output directory: %w", err)
	}

	// Find files to process
	files, err := p.findFiles(dir, options)
	if err != nil {
		return fmt.Errorf("error finding files: %w", err)
	}

	if len(files) == 0 {
		fmt.Printf("No files matching pattern '%s' found in '%s'\n", options.Pattern, dir)
		return nil
	}

	// Process each file
	fmt.Printf("Converting %d files...\n", len(files))
	for _, file := range files {
		if err := p.processFile(file, dir, options.OutputDir); err != nil {
			return err
		}
	}

	fmt.Printf("Successfully converted %d files to '%s'\n", len(files), options.OutputDir)
	return nil
}

func (p *FileProcessor) findFiles(dir string, options ProcessOptions) ([]string, error) {
	if options.Recursive {
		return p.findFilesRecursive(dir, options.Pattern)
	}
	files, err := filepath.Glob(filepath.Join(dir, options.Pattern))
	if err != nil {
		return nil, fmt.Errorf("error finding files with pattern %s: %w", options.Pattern, err)
	}
	return files, nil
}

func (p *FileProcessor) findFilesRecursive(dir, pattern string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, walkErr error) error {
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
	if err != nil {
		return nil, fmt.Errorf("error walking directory %s: %w", dir, err)
	}
	return files, nil
}

func (p *FileProcessor) processFile(file, inputDir, outputDir string) error {
	outputPath := GetOutputPath(file, inputDir, outputDir)

	// Create subdirectories if needed
	const defaultDirMode = 0755
	if err := os.MkdirAll(filepath.Dir(outputPath), defaultDirMode); err != nil {
		return fmt.Errorf("error creating directory for '%s': %w", outputPath, err)
	}

	fmt.Printf("Converting %s -> %s\n", file, outputPath)
	if err := p.converter.ConvertFile(file, outputPath); err != nil {
		return fmt.Errorf("error converting '%s': %w", file, err)
	}

	return nil
}