// Package main provides a command-line tool to convert markdown files to HTML with GitHub-style CSS.
package main

import (
	"github.com/sgaunet/mdtohtml/cmd"
)

var version = "development"

func main() {
	cmd.Version = version
	cmd.Execute()
}