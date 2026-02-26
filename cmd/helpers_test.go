package cmd

import (
	"errors"
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
