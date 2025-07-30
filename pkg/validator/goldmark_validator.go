package validator

import (
	"fmt"
	"os"

	"github.com/sgaunet/mdtohtml/pkg/converter"
)

// GoldmarkValidator implements Validator using the Goldmark library.
type GoldmarkValidator struct {
	converter converter.Converter
}

// NewGoldmarkValidator creates a new validator using the given converter.
func NewGoldmarkValidator(conv converter.Converter) *GoldmarkValidator {
	return &GoldmarkValidator{
		converter: conv,
	}
}

// Validate validates markdown content by attempting to convert it.
func (v *GoldmarkValidator) Validate(content []byte) error {
	_, err := v.converter.Convert(content)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}

// ValidateFile validates a markdown file by reading and validating its content.
func (v *GoldmarkValidator) ValidateFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file %s: %w", path, err)
	}

	return v.Validate(content)
}