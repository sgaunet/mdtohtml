package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	_ "embed"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
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

	// set up options
	var extensions = parser.NoIntraEmphasis |
		parser.Tables |
		parser.FencedCode |
		parser.Autolink |
		parser.Strikethrough |
		parser.SpaceHeadings

	var renderer markdown.Renderer

	// render the data into HTML
	var htmlFlags html.Flags
	htmlFlags |= html.UseXHTML
	if smartypants {
		htmlFlags |= html.Smartypants
	}
	if fractions {
		htmlFlags |= html.SmartypantsFractions
	}
	if latexdashes {
		htmlFlags |= html.SmartypantsLatexDashes
	}
	title := ""
	if page {
		htmlFlags |= html.CompletePage
		title = getTitle(input)
	}
	if toc {
		htmlFlags |= html.TOC
	}
	params := html.RendererOptions{
		Flags: htmlFlags,
		Title: title,
		//CSS:   css,
	}
	renderer = html.NewRenderer(params)

	// parse and render
	var output string
	parser := parser.NewWithExtensions(extensions)
	output = string(markdown.ToHTML(input, parser, renderer))

	// Add css github style
	var outputCssGitHub string
	reader := bufio.NewScanner(strings.NewReader(output))
	for reader.Scan() {
		if reader.Text() == "</head>" {
			for _, i := range [3]string{"<style>", cssgithub, "</style>"} {
				outputCssGitHub = outputCssGitHub + i + "\n"
			}
		}
		outputCssGitHub = outputCssGitHub + reader.Text() + "\n"
	}
	err = reader.Err()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error : %v", err)
		os.Exit(-1)
	}
	output = outputCssGitHub

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
