package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
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
	// errCSSFileNotFound is returned when the specified CSS file does not exist.
	errCSSFileNotFound = errors.New("CSS file not found")
	// errCSSFetchFailed is returned when fetching CSS from a URL fails.
	errCSSFetchFailed = errors.New("failed to fetch CSS from URL")
	// errConflictingCSS is returned when --no-css is combined with --css-file or --css-url.
	errConflictingCSS = errors.New("--no-css cannot be combined with --css-file or --css-url")
	// errBothCSSSourcesProvided is returned when both --css-file and --css-url are set.
	errBothCSSSourcesProvided = errors.New("--css-file and --css-url are mutually exclusive")
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

// readCSSFile reads a CSS file and returns its content.
func readCSSFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: %s", errCSSFileNotFound, path)
	}
	return string(data), nil
}

// fetchCSSURL fetches CSS content from a URL.
func fetchCSSURL(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec,noctx // user-provided URL is intentional
	if err != nil {
		return "", fmt.Errorf("%w: %w", errCSSFetchFailed, err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"%w: HTTP %d from %s",
			errCSSFetchFailed, resp.StatusCode, url,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("%w: %w", errCSSFetchFailed, err)
	}
	return string(body), nil
}

// validateCSSFlags checks for conflicting CSS flag combinations.
func validateCSSFlags(
	cssFile, cssURL string, noCSS bool,
) error {
	if noCSS && (cssFile != "" || cssURL != "") {
		return errConflictingCSS
	}
	if cssFile != "" && cssURL != "" {
		return errBothCSSSourcesProvided
	}
	return nil
}

// resolveCSSOptions validates flag combinations and resolves CSS content
// from files or URLs. It returns the CSS source text (replacing the
// default) and additional CSS text (appended).
func resolveCSSOptions(
	cssFile, cssURL, additionalCSSFile string, noCSS bool,
) (string, string, error) {
	if err := validateCSSFlags(cssFile, cssURL, noCSS); err != nil {
		return "", "", err
	}
	if noCSS {
		return "", "", nil
	}

	cssSource, err := resolveCSSSource(cssFile, cssURL)
	if err != nil {
		return "", "", err
	}

	var additionalCSS string
	if additionalCSSFile != "" {
		additionalCSS, err = readCSSFile(additionalCSSFile)
		if err != nil {
			return "", "", err
		}
	}

	return cssSource, additionalCSS, nil
}

// resolveCSSSource reads CSS from a file path or URL.
func resolveCSSSource(cssFile, cssURL string) (string, error) {
	if cssFile != "" {
		return readCSSFile(cssFile)
	}
	if cssURL != "" {
		return fetchCSSURL(cssURL)
	}
	return "", nil
}
