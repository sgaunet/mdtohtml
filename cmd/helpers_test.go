package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestValidateInputFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a valid file
	validFile := filepath.Join(tmpDir, "test.md")
	if err := os.WriteFile(validFile, []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name: "valid file",
			path: validFile,
		},
		{
			name:    "nonexistent file",
			path:    filepath.Join(tmpDir, "nonexistent.md"),
			wantErr: errFileNotFound,
		},
		{
			name:    "directory instead of file",
			path:    tmpDir,
			wantErr: errNotAFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputFile(tt.path)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("validateInputFile() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Errorf("validateInputFile() expected error wrapping %v, got nil", tt.wantErr)
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validateInputFile() error = %v, want wrapping %v", err, tt.wantErr)
			}
		})
	}
}

func TestValidateInputDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a regular file
	regularFile := filepath.Join(tmpDir, "file.txt")
	if err := os.WriteFile(regularFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	tests := []struct {
		name    string
		path    string
		wantErr error
	}{
		{
			name: "valid directory",
			path: tmpDir,
		},
		{
			name:    "nonexistent directory",
			path:    filepath.Join(tmpDir, "nonexistent"),
			wantErr: errDirNotFound,
		},
		{
			name:    "file instead of directory",
			path:    regularFile,
			wantErr: errNotADir,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateInputDir(tt.path)
			if tt.wantErr == nil {
				if err != nil {
					t.Errorf("validateInputDir() unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Errorf("validateInputDir() expected error wrapping %v, got nil", tt.wantErr)
				return
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("validateInputDir() error = %v, want wrapping %v", err, tt.wantErr)
			}
		})
	}
}

func TestResolveCSSOptions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a CSS file for testing
	cssFile := filepath.Join(tmpDir, "custom.css")
	if err := os.WriteFile(cssFile, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatalf("Failed to create CSS file: %v", err)
	}

	additionalFile := filepath.Join(tmpDir, "extra.css")
	if err := os.WriteFile(additionalFile, []byte("p { margin: 0; }"), 0644); err != nil {
		t.Fatalf("Failed to create additional CSS file: %v", err)
	}

	// HTTP test server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/style.css" {
			fmt.Fprint(w, "h1 { font-size: 2em; }")
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	tests := []struct {
		name              string
		cssFile           string
		cssURL            string
		additionalCSSFile string
		noCSS             bool
		wantSource        string
		wantAdditional    string
		wantErr           error
	}{
		{
			name: "no flags",
		},
		{
			name:       "cssFile set",
			cssFile:    cssFile,
			wantSource: "body { color: red; }",
		},
		{
			name:       "cssURL set",
			cssURL:     srv.URL + "/style.css",
			wantSource: "h1 { font-size: 2em; }",
		},
		{
			name:              "additionalCSSFile set",
			additionalCSSFile: additionalFile,
			wantAdditional:    "p { margin: 0; }",
		},
		{
			name:    "noCSS + cssFile conflict",
			noCSS:   true,
			cssFile: cssFile,
			wantErr: errConflictingCSS,
		},
		{
			name:    "noCSS + cssURL conflict",
			noCSS:   true,
			cssURL:  srv.URL + "/style.css",
			wantErr: errConflictingCSS,
		},
		{
			name:    "cssFile + cssURL conflict",
			cssFile: cssFile,
			cssURL:  srv.URL + "/style.css",
			wantErr: errBothCSSSourcesProvided,
		},
		{
			name:              "noCSS + additionalCSS silently returns empty",
			noCSS:             true,
			additionalCSSFile: additionalFile,
		},
		{
			name:    "nonexistent file",
			cssFile: filepath.Join(tmpDir, "nonexistent.css"),
			wantErr: errCSSFileNotFound,
		},
		{
			name:    "URL returns 404",
			cssURL:  srv.URL + "/notfound.css",
			wantErr: errCSSFetchFailed,
		},
		{
			name:    "URL unreachable",
			cssURL:  "http://127.0.0.1:1/unreachable.css",
			wantErr: errCSSFetchFailed,
		},
		{
			name:              "cssFile + additionalCSS both populated",
			cssFile:           cssFile,
			additionalCSSFile: additionalFile,
			wantSource:        "body { color: red; }",
			wantAdditional:    "p { margin: 0; }",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			source, additional, err := resolveCSSOptions(tt.cssFile, tt.cssURL, tt.additionalCSSFile, tt.noCSS)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error wrapping %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("error = %v, want wrapping %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if source != tt.wantSource {
				t.Errorf("cssSource = %q, want %q", source, tt.wantSource)
			}
			if additional != tt.wantAdditional {
				t.Errorf("additionalCSS = %q, want %q", additional, tt.wantAdditional)
			}
		})
	}
}
