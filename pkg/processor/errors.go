package processor

import "errors"

var (
	// ErrDirectoryNotExist is returned when a directory does not exist.
	ErrDirectoryNotExist = errors.New("directory does not exist")
)