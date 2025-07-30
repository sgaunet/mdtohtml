// Package validator provides markdown validation functionality
package validator

// Validator defines the interface for markdown validation.
type Validator interface {
	// Validate validates markdown content and returns an error if invalid
	Validate(content []byte) error

	// ValidateFile validates a markdown file and returns an error if invalid
	ValidateFile(path string) error
}