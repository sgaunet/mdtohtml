package cmd

import (
	"errors"
	"fmt"
	"os"
)

var (
	// errFileNotFound is returned when an input file does not exist.
	errFileNotFound = errors.New("input file not found")
	// errNotAFile is returned when the input path is a directory instead of a file.
	errNotAFile = errors.New("input path is a directory, not a file")
	// errDirNotFound is returned when an input directory does not exist.
	errDirNotFound = errors.New("input directory not found")
	// errNotADir is returned when the input path is a file instead of a directory.
	errNotADir = errors.New("input path is a file, not a directory")
	// errPermissionDenied is returned when the input path cannot be accessed.
	errPermissionDenied = errors.New("permission denied")
)

// validateInputFile checks that path exists and is a regular file.
func validateInputFile(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", errFileNotFound, path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("%w: %s", errPermissionDenied, path)
		}
		return fmt.Errorf("cannot access input file '%s': %w", path, err)
	}
	if info.IsDir() {
		return fmt.Errorf("%w: %s", errNotAFile, path)
	}
	return nil
}

// validateInputDir checks that path exists and is a directory.
func validateInputDir(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: %s", errDirNotFound, path)
		}
		if os.IsPermission(err) {
			return fmt.Errorf("%w: %s", errPermissionDenied, path)
		}
		return fmt.Errorf("cannot access input directory '%s': %w", path, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%w: %s", errNotADir, path)
	}
	return nil
}
