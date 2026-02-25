package processor

import "errors"

var (
	// ErrDirectoryNotExist is returned when a directory does not exist.
	ErrDirectoryNotExist = errors.New("directory does not exist")
	// ErrInvalidPattern is returned when a glob pattern is malformed.
	ErrInvalidPattern = errors.New("invalid file pattern")
)