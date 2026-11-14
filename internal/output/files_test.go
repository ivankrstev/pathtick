package output_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ivankrstev/pathtick/internal/output"
)

func TestWriteOutputFilesCreatesFiles(t *testing.T) {
	outputDir := filepath.Join(t.TempDir(), "nested", "exports")
	files := []output.File{
		{Extension: ".md", Data: []byte("markdown content")},
		{Extension: "html", Data: []byte("html content")},
	}

	if err := output.WriteFiles(outputDir, "PathTick", false, files); err != nil {
		t.Fatalf("WriteFiles() error = %v", err)
	}

	for _, test := range []struct {
		name string
		want string
	}{
		{name: "PathTick.md", want: "markdown content"},
		{name: "PathTick.html", want: "html content"},
	} {
		got, err := os.ReadFile(filepath.Join(outputDir, test.name))
		if err != nil {
			t.Fatalf("read %s: %v", test.name, err)
		}
		if string(got) != test.want {
			t.Errorf("%s contents = %q, want %q", test.name, string(got), test.want)
		}
	}
}

func TestWriteOutputFilesRefusesOverwrite(t *testing.T) {
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "PathTick.md")
	const existingContents = "keep this file"
	if err := os.WriteFile(outputPath, []byte(existingContents), 0o644); err != nil {
		t.Fatalf("create existing file: %v", err)
	}

	err := output.WriteFiles(outputDir, "PathTick", false, []output.File{
		{Extension: "md", Data: []byte("new content")},
	})
	if err == nil {
		t.Fatal("WriteFiles() error = nil, want an error for an existing output file")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("WriteFiles() error = %q, want an existing-file error", err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read existing file: %v", err)
	}
	if string(got) != existingContents {
		t.Errorf("existing file contents = %q, want %q", string(got), existingContents)
	}
}

func TestWriteOutputFilesOverwrites(t *testing.T) {
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "PathTick.md")
	if err := os.WriteFile(outputPath, []byte("old content"), 0o644); err != nil {
		t.Fatalf("create existing file: %v", err)
	}

	if err := output.WriteFiles(outputDir, "PathTick", true, []output.File{
		{Extension: "md", Data: []byte("new content")},
	}); err != nil {
		t.Fatalf("WriteFiles() with overwrite error = %v", err)
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read overwritten file: %v", err)
	}
	if string(got) != "new content" {
		t.Errorf("overwritten file contents = %q, want %q", string(got), "new content")
	}
}

func TestWriteOutputFilesRejectsInvalidName(t *testing.T) {
	tests := []struct {
		name     string
		baseName string
	}{
		{name: "empty", baseName: ""},
		{name: "dot", baseName: "."},
		{name: "dot dot", baseName: ".."},
		{name: "contains a directory", baseName: filepath.Join("nested", "PathTick")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := output.WriteFiles(t.TempDir(), test.baseName, false, []output.File{
				{Extension: "md", Data: []byte("content")},
			})
			if err == nil {
				t.Fatal("WriteFiles() error = nil, want an invalid-name error")
			}
			if !strings.Contains(err.Error(), "invalid output file name") {
				t.Errorf("WriteFiles() error = %q, want an invalid-name error", err)
			}
		})
	}
}

func TestWriteOutputFilesRejectsInvalidExtension(t *testing.T) {
	tests := []struct {
		name      string
		extension string
	}{
		{name: "empty", extension: ""},
		{name: "dot only", extension: "."},
		{name: "parent path", extension: "../md"},
		{name: "nested path", extension: filepath.Join("nested", "md")},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := output.WriteFiles(t.TempDir(), "PathTick", false, []output.File{
				{Extension: test.extension, Data: []byte("content")},
			})
			if err == nil {
				t.Fatal("WriteFiles() error = nil, want an invalid-extension error")
			}
			if !strings.Contains(err.Error(), "invalid output file extension") {
				t.Errorf("WriteFiles() error = %q, want an invalid-extension error", err)
			}
		})
	}
}

func TestWriteOutputFilesRejectsDuplicateExtensions(t *testing.T) {
	err := output.WriteFiles(t.TempDir(), "PathTick", false, []output.File{
		{Extension: ".md", Data: []byte("first")},
		{Extension: "md", Data: []byte("second")},
	})
	if err == nil {
		t.Fatal("WriteFiles() error = nil, want a duplicate-output error")
	}
	if !strings.Contains(err.Error(), "duplicate output file") {
		t.Errorf("WriteFiles() error = %q, want a duplicate-output error", err)
	}
}

func TestWriteOutputFilesRejectsEmptyOutputList(t *testing.T) {
	err := output.WriteFiles(t.TempDir(), "PathTick", false, nil)
	if err == nil {
		t.Fatal("WriteFiles() error = nil, want an empty-output error")
	}
	if !strings.Contains(err.Error(), "no output files to write") {
		t.Errorf("WriteFiles() error = %q, want an empty-output error", err)
	}
}

func TestWriteOutputFilesRejectsEmptyOutputDirectory(t *testing.T) {
	err := output.WriteFiles("", "PathTick", false, []output.File{
		{Extension: "md", Data: []byte("content")},
	})
	if err == nil {
		t.Fatal("WriteFiles() error = nil, want an empty-output-directory error")
	}
	if !strings.Contains(err.Error(), "output directory must not be empty") {
		t.Errorf("WriteFiles() error = %q, want an empty-output-directory error", err)
	}
}

func TestWriteOutputFilesChecksAllTargetsBeforeWriting(t *testing.T) {
	outputDir := t.TempDir()
	existingPath := filepath.Join(outputDir, "PathTick.html")
	const existingContents = "keep this file"
	if err := os.WriteFile(existingPath, []byte(existingContents), 0o644); err != nil {
		t.Fatalf("create existing file: %v", err)
	}

	err := output.WriteFiles(outputDir, "PathTick", false, []output.File{
		{Extension: "md", Data: []byte("markdown content")},
		{Extension: "html", Data: []byte("html content")},
	})
	if err == nil {
		t.Fatal("WriteFiles() error = nil, want an error for the existing HTML output")
	}

	markdownPath := filepath.Join(outputDir, "PathTick.md")
	if _, err := os.Stat(markdownPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Markdown output stat error = %v, want os.ErrNotExist", err)
	}

	got, err := os.ReadFile(existingPath)
	if err != nil {
		t.Fatalf("read existing HTML file: %v", err)
	}
	if string(got) != existingContents {
		t.Errorf("existing HTML file contents = %q, want %q", string(got), existingContents)
	}
}
