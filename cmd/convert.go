package cmd

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

//go:embed github-markdown.css
var cssgithub string

func runConversion(inputFilePath, outputFilePath string, smartypants, latexdashes, fractions bool) error {
	input, err := os.ReadFile(inputFilePath)
	if err != nil {
		return fmt.Errorf("error reading from %s: %w", inputFilePath, err)
	}

	md := createGoldmarkProcessor(smartypants, fractions, latexdashes)

	output, err := processMarkdown(input, md)
	if err != nil {
		return err
	}

	output = addCSS(output)

	err = writeOutput(output, outputFilePath)
	if err != nil {
		return err
	}

	return nil
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

func processMarkdown(input []byte, md goldmark.Markdown) (string, error) {
	var buf bytes.Buffer
	if err := md.Convert(input, &buf); err != nil {
		return "", fmt.Errorf("error converting markdown: %w", err)
	}

	title := getTitle(input)
	output := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>%s</title>
</head>
<body>
%s
</body>
</html>`, title, buf.String())

	return output, nil
}

func addCSS(output string) string {
	cssInjection := fmt.Sprintf("<style>\n%s\n</style>\n</head>", cssgithub)
	return strings.Replace(output, "</head>", cssInjection, 1)
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