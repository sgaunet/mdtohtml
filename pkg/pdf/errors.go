package pdf

import "errors"

// ErrUnknownPageSize is returned when an unrecognised page size identifier is supplied.
var ErrUnknownPageSize = errors.New("unknown page size")

// ErrInvalidMargin is returned when a margin string cannot be parsed.
var ErrInvalidMargin = errors.New("invalid margin")
