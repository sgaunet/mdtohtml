package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
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

var version string = "development"

func printVersion() {
	fmt.Println(version)
}

func main() {
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

	if vOption {
		printVersion()
		os.Exit(0)
	}

	// enforce implied options
	page = true

	// read the input
	var input []byte
	var err error
	args := flag.Args()

	if len(args) != 2 {
		flag.Usage()
		os.Exit(1)
	}

	inputFilePath := args[0]
	outputFilePath := args[1]

	if input, err = ioutil.ReadFile(inputFilePath); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from", inputFilePath, ":", err)
		os.Exit(-1)
	}

	// Set up Goldmark with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM, // GitHub Flavored Markdown
			extension.DefinitionList,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithXHTML(),
			html.WithUnsafe(), // Allow raw HTML if needed
		),
	)

	// If smartypants is enabled, add the typographer extension
	if smartypants || fractions || latexdashes {
		md = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.DefinitionList,
				extension.Footnote,
				extension.Typographer, // This enables typographic substitutions
			),
			goldmark.WithParserOptions(
				parser.WithAutoHeadingID(),
			),
			goldmark.WithRendererOptions(
				html.WithXHTML(),
				html.WithUnsafe(),
				// Enable typographic substitutions in the renderer
				html.WithHardWraps(),
				html.WithXHTML(),
				html.WithUnsafe(),
			),
		)
	}

	// Set page title if needed
	var title string
	if page {
		title = getTitle(input)
	}

	// Create a buffer to hold the HTML output
	var buf bytes.Buffer

	// Convert markdown to HTML
	if err := md.Convert(input, &buf); err != nil {
		fmt.Fprintln(os.Stderr, "Error converting markdown:", err)
		os.Exit(-1)
	}

	// Get the output as a string
	output := buf.String()

	// Wrap in HTML page if needed
	if page {
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

	// Add CSS github style
	if page {
		// Inject CSS before the closing </head> tag
		cssInjection := fmt.Sprintf("<style>\n%s\n</style>\n</head>", cssgithub)
		output = strings.Replace(output, "</head>", cssInjection, 1)
	} else {
		// If not a full page, just prepend the style
		output = fmt.Sprintf("<style>\n%s\n</style>\n%s", cssgithub, output)
	}

	// output the result
	var out *os.File
	if out, err = os.Create(outputFilePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v", outputFilePath, err)
		os.Exit(-1)
	}
	defer out.Close()

	if _, err = out.WriteString(output); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing output:", err)
		os.Exit(-1)
	}
}
