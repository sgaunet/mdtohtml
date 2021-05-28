package main

import (
	"bytes"
	"io/ioutil"
	"os"

	wkhtmltopdf "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

func createPDF(htmlFilePath string, pdfFilePath string) error {
	pdfg := wkhtmltopdf.NewPDFPreparer()
	htmlfile, err := ioutil.ReadFile(htmlFilePath)
	if err != nil {
		return err
	}

	pdfg.AddPage(wkhtmltopdf.NewPageReader(bytes.NewReader(htmlfile)))
	pdfg.Dpi.Set(600)

	// The contents of htmlsimple.html are saved as base64 string in the JSON file
	jb, err := pdfg.ToJSON()
	if err != nil {
		return err
	}

	// Server code
	pdfgFromJSON, err := wkhtmltopdf.NewPDFGeneratorFromJSON(bytes.NewReader(jb))
	if err != nil {
		return err
	}

	err = pdfgFromJSON.Create()
	if err != nil {
		return err
	}
	g, err := os.Create(pdfFilePath)

	if err != nil {
		return err
	}

	defer g.Close()
	g.Write(pdfgFromJSON.Bytes())
	return err
}
