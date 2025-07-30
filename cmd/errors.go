package cmd

import "errors"

var (
	// ErrInputNotExist is returned when the input file or directory does not exist.
	ErrInputNotExist = errors.New("input does not exist")
)