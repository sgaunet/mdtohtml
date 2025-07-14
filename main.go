// Package main provides a command-line tool to convert markdown files to HTML with GitHub-style CSS.
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

const defaultTitle = ""

//go:embed github-markdown.css
var cssgithub string

var version = "development"

func printVersion() {
	fmt.Println(version)
}

func setupFlags() (bool, bool, bool, bool, bool, bool) {
	var page, toc, smartypants, latexdashes, fractions, vOption bool
	flag.BoolVar(&page, "page", false, "Generate a standalone HTML page")
	flag.BoolVar(&toc, "toc", false, "Generate a table of contents")
	flag.BoolVar(&smartypants, "smartypants", true, "Apply smartypants-style substitutions")
	flag.BoolVar(&latexdashes, "latexdashes", true, "Use LaTeX-style dash rules for smartypants")
	flag.BoolVar(&fractions, "fractions", true, "Use improved fraction rules for smartypants")
	flag.BoolVar(&vOption, "v", false, "Get version")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Markdown Processor "+
			"\nAvailable at http://github.com/gomarkdown/markdown/cmd/mdtohtml\n\n"+
			"Copyright © 2011 Russ Ross <russ@russross.com>\n"+
			"Copyright © 2018 Krzysztof Kowalczyk <https://blog.kowalczyk.info>\n"+
			"Distributed under the Simplified BSD License\n"+
			"Usage:\n"+
			"  %s [options] inputfile outputfile\n\n"+
			"Options:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	return page, toc, smartypants, latexdashes, fractions, vOption
}

func createGoldmarkProcessor(smartypants, fractions, latexdashes bool) goldmark.Markdown {
	if smartypants || fractions || latexdashes {
		return goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.DefinitionList,
				extension.Footnote,
				extension.Typographer,
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithXHTML(),
				html.WithUnsafe(),
				html.WithHardWraps(),
			),
		)
	}

	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(),
		),
	)
}

func processMarkdown(input []byte, md goldmark.Markdown, page bool) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return "", fmt.Errorf("error converting markdown: %w", err)
	}

	output := buf.String()

	if page {
		title := getTitle(input)
		output = fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>%s</title>
</head>
<body>
%s
</body>
</html>`, title, output)
	}

	return output, nil
}

func addCSS(output string, page bool) string {
	if page {
		cssInjection := fmt.Sprintf("<style>\n%s\n</style>\n</head>", cssgithub)
		return strings.Replace(output, "</head>", cssInjection, 1)
	}
	return fmt.Sprintf("<style>\n%s\n</style>\n%s", cssgithub, output)
}

func writeOutput(output, outputFilePath string) error {
	out, err := os.Create(outputFilePath)
	if err != nil {
		return fmt.Errorf("error creating %s: %w", outputFilePath, err)
	}

	defer func() {
		if closeErr := out.Close(); closeErr != nil {
			fmt.Fprintf(os.Stderr, "Error closing file: %v\n", closeErr)
		}
	}()

	if _, err = out.WriteString(output); err != nil {
		return fmt.Errorf("error writing output: %w", err)
	}

	return nil
}

func main() {
	_, toc, smartypants, latexdashes, fractions, vOption := setupFlags()
	_ = toc // Variable assigned but not used

	if vOption {
		printVersion()
		os.Exit(0)
	}

	// enforce implied options
	page := true

	args := flag.Args()
	expectedArgs := 2
	if len(args) != expectedArgs {
		flag.Usage()
		os.Exit(1)
	}

	inputFilePath := args[0]
	outputFilePath := args[1]

	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from", inputFilePath, ":", err)
		os.Exit(-1)
	}

	md := createGoldmarkProcessor(smartypants, fractions, latexdashes)

	output, err := processMarkdown(input, md, page)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}

	output = addCSS(output, page)

	if err := writeOutput(output, outputFilePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}