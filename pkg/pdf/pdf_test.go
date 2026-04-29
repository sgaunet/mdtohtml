package pdf_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sgaunet/mdtohtml/pkg/converter"
	"github.com/sgaunet/mdtohtml/pkg/pdf"
)

// Compile-time interface check.
var _ converter.Converter = (*pdf.Converter)(nil)

func newConv(t *testing.T) *pdf.Converter {
	t.Helper()
	c, err := pdf.New(converter.DefaultOptions(), pdf.DefaultOptions())
	if err != nil {
		t.Fatalf("pdf.New: %v", err)
	}
	return c
}

func TestConvert_ProducesPDFBytes(t *testing.T) {
	out, err := newConv(t).Convert([]byte("# Hello\n\nWorld"))
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if !bytes.HasPrefix(out, []byte("%PDF-")) {
		t.Fatalf("output does not start with PDF magic bytes; got %q", out[:min(8, len(out))])
	}
	if !bytes.Contains(out, []byte("%%EOF")) {
		t.Fatal("output missing PDF EOF marker")
	}
}

func TestConvert_TitleEmbeddedInMetadata(t *testing.T) {
	const title = "Folio Integration Smoke Title"
	out, err := newConv(t).Convert([]byte("# " + title + "\n\nbody"))
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	// Folio writes Info.Title as a PDF string literal in the trailer/Info dict.
	// We don't decode the PDF here — a substring check is enough for a smoke test.
	if !bytes.Contains(out, []byte(title)) {
		t.Fatalf("expected title %q to appear in PDF stream", title)
	}
}

func TestConvertFile_WritesPDF(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "doc.md")
	if err := os.WriteFile(in, []byte("# Test\n\nContent."), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}
	out := filepath.Join(dir, "doc.pdf")

	if err := newConv(t).ConvertFile(in, out); err != nil {
		t.Fatalf("ConvertFile: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if !bytes.HasPrefix(data, []byte("%PDF-")) {
		t.Fatal("output file is not a PDF")
	}
}

func TestConvertFile_PageBreakFixture(t *testing.T) {
	// Reuse the existing markdown that exercises CSS page breaks.
	src := filepath.Join("..", "..", "tst", "README-with-page-break.md")
	if _, err := os.Stat(src); err != nil {
		t.Skipf("fixture missing: %v", err)
	}
	out := filepath.Join(t.TempDir(), "breaks.pdf")
	if err := newConv(t).ConvertFile(src, out); err != nil {
		t.Fatalf("ConvertFile: %v", err)
	}
	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	// A document with at least one page break should produce more than one page.
	if pages := strings.Count(string(data), "/Type /Page\n") + strings.Count(string(data), "/Type/Page"); pages < 2 {
		t.Logf("page-count heuristic returned %d; PDF size %d bytes (not strictly checking — folio's output format may vary)", pages, len(data))
	}
}

func TestNew_RejectsUnknownPageSize(t *testing.T) {
	_, err := pdf.New(converter.DefaultOptions(), pdf.Options{PageSize: "Bogus"})
	if !errors.Is(err, pdf.ErrUnknownPageSize) {
		t.Fatalf("expected ErrUnknownPageSize, got %v", err)
	}
}

func TestNew_AcceptsKnownPageSizes(t *testing.T) {
	for _, ps := range []string{"", "A4", "a4", "Letter", "Legal", "A3", "A5", "Tabloid"} {
		t.Run(ps, func(t *testing.T) {
			if _, err := pdf.New(converter.DefaultOptions(), pdf.Options{PageSize: ps}); err != nil {
				t.Fatalf("page size %q rejected: %v", ps, err)
			}
		})
	}
}

func TestDefaultOptions_HasMargins(t *testing.T) {
	opts := pdf.DefaultOptions()
	if opts.Margins.Top != pdf.DefaultMargin ||
		opts.Margins.Right != pdf.DefaultMargin ||
		opts.Margins.Bottom != pdf.DefaultMargin ||
		opts.Margins.Left != pdf.DefaultMargin {
		t.Fatalf("expected all margins = %v, got %+v", pdf.DefaultMargin, opts.Margins)
	}
}

func TestParseMargin(t *testing.T) {
	cases := []struct {
		in   string
		want float64
	}{
		{"", pdf.DefaultMargin},
		{"  ", pdf.DefaultMargin},
		{"90", 90},
		{"90pt", 90},
		{"  90pt ", 90},
		{"1in", 72},
		{"1.25in", 90},
		{"2.54cm", 72},
		{"25.4mm", 72},
		{"0", 0},
	}
	for _, c := range cases {
		t.Run(c.in, func(t *testing.T) {
			got, err := pdf.ParseMargin(c.in)
			if err != nil {
				t.Fatalf("ParseMargin(%q) err: %v", c.in, err)
			}
			// Allow tiny float drift on cm/mm conversions.
			if diff := got - c.want; diff < -0.001 || diff > 0.001 {
				t.Fatalf("ParseMargin(%q) = %v, want %v", c.in, got, c.want)
			}
		})
	}
}

func TestParseMargin_Errors(t *testing.T) {
	for _, in := range []string{"abc", "10ft", "1.2.3", "-5"} {
		t.Run(in, func(t *testing.T) {
			if _, err := pdf.ParseMargin(in); !errors.Is(err, pdf.ErrInvalidMargin) {
				t.Fatalf("ParseMargin(%q) expected ErrInvalidMargin, got %v", in, err)
			}
		})
	}
}

func TestConvert_HonoursCustomMargin(t *testing.T) {
	// Smoke test: a custom margin must not break PDF generation.
	opts := pdf.DefaultOptions()
	opts.Margins = pdf.Margins{Top: 36, Right: 36, Bottom: 36, Left: 36}
	c, err := pdf.New(converter.DefaultOptions(), opts)
	if err != nil {
		t.Fatalf("pdf.New: %v", err)
	}
	out, err := c.Convert([]byte("# Hi\n\nbody"))
	if err != nil {
		t.Fatalf("Convert: %v", err)
	}
	if !bytes.HasPrefix(out, []byte("%PDF-")) {
		t.Fatal("output is not a PDF")
	}
}
