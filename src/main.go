package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"

	_ "embed"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

const defaultTitle = ""

//go:embed github-markdown.css
var cssgithub string

func main() {
	var cssgh, page, toc, xhtml, latex, smartypants, latexdashes, fractions bool
	var css, cpuprofile string
	flag.BoolVar(&page, "page", false,
		"Generate a standalone HTML page (implies -latex=false)")
	flag.BoolVar(&toc, "toc", false,
		"Generate a table of contents (implies -latex=false)")
	flag.BoolVar(&xhtml, "xhtml", true,
		"Use XHTML-style tags in HTML output")
	flag.BoolVar(&cssgh, "cssgh", true,
		"Github style")
	//flag.BoolVar(&latex, "latex", false,
	//	"Generate LaTeX output instead of HTML")
	flag.BoolVar(&smartypants, "smartypants", true,
		"Apply smartypants-style substitutions")
	flag.BoolVar(&latexdashes, "latexdashes", true,
		"Use LaTeX-style dash rules for smartypants")
	flag.BoolVar(&fractions, "fractions", true,
		"Use improved fraction rules for smartypants")
	flag.StringVar(&cpuprofile, "cpuprofile", "",
		"Write cpu profile to a file")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Markdown Processor "+
			"\nAvailable at http://github.com/gomarkdown/markdown/cmd/mdtohtml\n\n"+
			"Copyright © 2011 Russ Ross <russ@russross.com>\n"+
			"Copyright © 2018 Krzysztof Kowalczyk <https://blog.kowalczyk.info>\n"+
			"Distributed under the Simplified BSD License\n"+
			"Usage:\n"+
			"  %s [options] [inputfile [outputfile]]\n\n"+
			"Options:\n",
			os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// enforce implied options
	if css != "" || cssgh {
		page = true
	}
	if page {
		latex = false
	}
	if toc {
		latex = false
	}

	// turn on profiling?
	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

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

	// Create temporary file
	inputDir := filepath.Dir(args[0])
	tmpFile, err := ioutil.TempFile(inputDir, "mdtohtml-")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

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
	if latex {
		// render the data into LaTeX
		//renderer = markdown.LatexRenderer(0)
	} else {
		// render the data into HTML
		var htmlFlags html.Flags
		if xhtml {
			htmlFlags |= html.UseXHTML
		}
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
			CSS:   css,
		}
		renderer = html.NewRenderer(params)
	}

	// parse and render
	var output []byte
	parser := parser.NewWithExtensions(extensions)
	output = markdown.ToHTML(input, parser, renderer)

	// output the result
	var out *os.File
	if len(args) == 2 {
		if out, err = os.Create(tmpFile.Name()); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", tmpFile.Name(), err)
			os.Exit(-1)
		}
		defer out.Close()
	} else {
		out = os.Stdout
	}

	if _, err = out.Write(output); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing output:", err)
		os.Exit(-1)
	}
	out.Close()

	// html with github-markdown.css
	if cssgh && len(args) == 2 {
		tmpFile2, err := ioutil.TempFile(inputDir, "mdtohtml-")
		if err != nil {
			log.Fatal(err)
		}
		defer os.Remove(tmpFile2.Name())

		file, err := os.Open(tmpFile.Name())
		if err != nil {
			panic(err)
		}
		defer file.Close()

		f, err := os.Create(tmpFile2.Name())

		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()

		// Start reading from the file with a reader.
		reader := bufio.NewReader(file)
		var line string
		for {
			line, err = reader.ReadString('\n')
			if err != nil && err != io.EOF {
				break
			}

			if line == "</head>\n" {
				for _, i := range [3]string{"<style>", cssgithub, "</style>"} {
					_, err2 := f.WriteString(i)

					if err2 != nil {
						log.Fatal(err2)
					}
				}
			}
			_, err2 := f.WriteString(line)

			if err2 != nil {
				log.Fatal(err2)
			}

			if err != nil {
				break
			}
		}
		if err != io.EOF {
			fmt.Printf(" > Failed with error: %v\n", err)
			panic(err)
		}

		err = f.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		err = file.Close()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		// err = os.Rename(tmpFile2.Name(), tmpFile.Name())
		// if err != nil {
		// 	fmt.Fprintf(os.Stderr, "Error renaming %s to %s : %v", tmpFile2.Name(), tmpFile.Name(), err)
		// 	os.Exit(1)
		// }
		if input, err = ioutil.ReadFile(tmpFile2.Name()); err != nil {
			fmt.Fprintln(os.Stderr, "Error reading from", tmpFile.Name(), ":", err)
			os.Exit(-1)
		}

		out2, err := os.Create(tmpFile.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating %s: %v", tmpFile.Name(), err)
			os.Exit(-1)
		}
		if _, err = out2.Write(input); err != nil {
			fmt.Fprintln(os.Stderr, "Error writing output:", err)
			os.Exit(-1)
		}
		out2.Close()
		tmpFile2.Close()
		os.Remove(tmpFile2.Name())
	}

	if input, err = ioutil.ReadFile(tmpFile.Name()); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from", tmpFile.Name(), ":", err)
		os.Exit(-1)
	}

	out2, err := os.Create(outputFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating %s: %v", tmpFile.Name(), err)
		os.Exit(-1)
	}
	if _, err = out2.Write(input); err != nil {
		fmt.Fprintln(os.Stderr, "Error writing output:", err)
		os.Exit(-1)
	}
	out2.Close()
	err = tmpFile.Close()
	err = os.Remove(tmpFile.Name())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting %s : %s \n", tmpFile.Name(), err)
		os.Exit(-1)
	}
}
