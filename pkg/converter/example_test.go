package converter_test

import (
	"fmt"

	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/parser"
	"github.com/sgaunet/mdtohtml/pkg/template"
)

func ExampleNewCompleteConverter() {
	opts := converter.DefaultOptions()
	conv := converter.NewCompleteConverter(opts)

	output, err := conv.Convert([]byte("# Hello\n\nWorld"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(len(output) > 0) // true: produces a full HTML document
	// Output: true
}

func ExampleNewCompleteConverterWithComponents() {
	gc := converter.NewGoldmarkConverter(converter.DefaultOptions())
	te := parser.NewMarkdownTitleExtractor()
	ht := template.NewGitHubTemplate()

	conv := converter.NewCompleteConverterWithComponents(gc, te, ht)

	output, err := conv.Convert([]byte("# Custom\n\nUsing custom components"))
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println(len(output) > 0) // true: produces a full HTML document
	// Output: true
}
