package pdf

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/carlos7ags/folio/document"
)

// PageSizeA4 and PageSizeLetter are the recognised string identifiers for
// the most common page sizes. Additional folio sizes can be added here as
// they are needed.
const (
	PageSizeA4      = "A4"
	PageSizeLetter  = "Letter"
	PageSizeLegal   = "Legal"
	PageSizeA3      = "A3"
	PageSizeA5      = "A5"
	PageSizeTabloid = "Tabloid"
)

// DefaultMargin is the default page margin in PDF points (1.25 inch).
// 72 points = 1 inch, so 90 points = 1.25 inch.
const DefaultMargin = 90.0

// Unit conversion constants (PDF points per unit).
const (
	pointsPerInch = 72.0
	pointsPerCm   = pointsPerInch / 2.54
	pointsPerMm   = pointsPerCm / 10.0
)

// Margins defines page margins in PDF points (72pt = 1 inch).
type Margins struct {
	Top    float64
	Right  float64
	Bottom float64
	Left   float64
}

// Options configures PDF generation.
type Options struct {
	// PageSize is the page size identifier (e.g. "A4", "Letter").
	// Empty defaults to A4.
	PageSize string

	// Margins are the page margins applied unless overridden by an @page CSS rule.
	Margins Margins
}

// DefaultOptions returns sensible PDF defaults (A4, 1.25 inch margins).
func DefaultOptions() Options {
	return Options{
		PageSize: PageSizeA4,
		Margins: Margins{
			Top:    DefaultMargin,
			Right:  DefaultMargin,
			Bottom: DefaultMargin,
			Left:   DefaultMargin,
		},
	}
}

// ParseMargin converts a margin string into PDF points. Empty input returns
// DefaultMargin. Accepted unit suffixes: pt (default), in, cm, mm.
func ParseMargin(s string) (float64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return DefaultMargin, nil
	}

	num, unit := splitMargin(s)
	v, err := strconv.ParseFloat(num, 64)
	if err != nil {
		return 0, fmt.Errorf("%w: %s", ErrInvalidMargin, s)
	}
	if v < 0 {
		return 0, fmt.Errorf("%w: %s (must be non-negative)", ErrInvalidMargin, s)
	}

	switch unit {
	case "", "pt":
		return v, nil
	case "in":
		return v * pointsPerInch, nil
	case "cm":
		return v * pointsPerCm, nil
	case "mm":
		return v * pointsPerMm, nil
	default:
		return 0, fmt.Errorf("%w: unknown unit %q", ErrInvalidMargin, unit)
	}
}

// splitMargin separates the numeric portion of a margin value from its
// trailing unit suffix. The unit is returned lowercased.
func splitMargin(s string) (string, string) {
	for i, r := range s {
		if (r >= '0' && r <= '9') || r == '.' || r == '-' || r == '+' {
			continue
		}
		return s[:i], strings.ToLower(s[i:])
	}
	return s, ""
}

// resolvePageSize maps a case-insensitive identifier to a folio PageSize.
// Returns an error for unknown identifiers.
func resolvePageSize(name string) (document.PageSize, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "", "a4":
		return document.PageSizeA4, nil
	case "letter":
		return document.PageSizeLetter, nil
	case "legal":
		return document.PageSizeLegal, nil
	case "a3":
		return document.PageSizeA3, nil
	case "a5":
		return document.PageSizeA5, nil
	case "tabloid":
		return document.PageSizeTabloid, nil
	default:
		return document.PageSize{}, fmt.Errorf("%w: %s", ErrUnknownPageSize, name)
	}
}
