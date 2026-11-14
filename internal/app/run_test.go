package app_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/ivankrstev/pathtick/internal/app"
)

func TestRunWritesMarkdownFile(t *testing.T) {
	projectDir := createProject(t)
	outputDir := filepath.Join(t.TempDir(), "nested", "output")

	options := app.DefaultOptions()
	options.TargetDir = projectDir
	options.Title = "A custom title"
	options.OutputDir = outputDir
	options.OutputBaseName = "checklist"

	if err := app.Run(options, &bytes.Buffer{}); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	outputPath := filepath.Join(outputDir, "checklist.md")
	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read output file: %v", err)
	}

	want := "# A custom title\n\n" +
		"- [ ] project\n" +
		"  - [ ] README\\.md\n"
	if string(got) != want {
		t.Errorf("output file contents:\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}

func TestRunUsesDirectoryNameAsDefaultTitle(t *testing.T) {
	projectDir := createProject(t)
	outputDir := filepath.Join(t.TempDir(), "should-not-be-created")

	options := app.DefaultOptions()
	options.TargetDir = projectDir
	options.OutputDir = outputDir
	options.OutputBaseName = "ignored-in-stdout-mode"
	options.Stdout = true

	var stdout bytes.Buffer
	if err := app.Run(options, &stdout); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	want := "# project\n\n" +
		"- [ ] project\n" +
		"  - [ ] README\\.md\n"
	if got := stdout.String(); got != want {
		t.Errorf("stdout contents:\ngot:\n%q\nwant:\n%q", got, want)
	}

	if _, err := os.Stat(outputDir); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("output directory stat error = %v, want os.ErrNotExist", err)
	}
}

func TestRunReturnsScanError(t *testing.T) {
	options := app.DefaultOptions()
	options.TargetDir = filepath.Join(t.TempDir(), "missing")

	err := app.Run(options, &bytes.Buffer{})
	if err == nil {
		t.Fatal("Run() error = nil, want an error for a missing target")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Run() error = %v, want an error wrapping os.ErrNotExist", err)
	}
}

func TestRunDoesNotOverwriteExistingFileByDefault(t *testing.T) {
	projectDir := createProject(t)
	outputDir := t.TempDir()
	outputPath := filepath.Join(outputDir, "PathTick.md")
	const existingContents = "keep this file"

	if err := os.WriteFile(outputPath, []byte(existingContents), 0o644); err != nil {
		t.Fatalf("create existing output file: %v", err)
	}

	options := app.DefaultOptions()
	options.TargetDir = projectDir
	options.OutputDir = outputDir

	err := app.Run(options, &bytes.Buffer{})
	if err == nil {
		t.Fatal("Run() error = nil, want an error when the output already exists")
	}

	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read existing output file: %v", err)
	}
	if string(got) != existingContents {
		t.Errorf("existing output file contents = %q, want %q", string(got), existingContents)
	}
}

func TestRunWrapsStdoutWriteError(t *testing.T) {
	projectDir := createProject(t)
	writeErr := errors.New("writer failed")

	options := app.DefaultOptions()
	options.TargetDir = projectDir
	options.Stdout = true

	err := app.Run(options, failingWriter{err: writeErr})
	if !errors.Is(err, writeErr) {
		t.Fatalf("Run() error = %v, want it to wrap %v", err, writeErr)
	}
}

func createProject(t *testing.T) string {
	t.Helper()

	projectDir := filepath.Join(t.TempDir(), "project")
	if err := os.MkdirAll(projectDir, 0o755); err != nil {
		t.Fatalf("create project directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(projectDir, "README.md"), []byte("readme"), 0o644); err != nil {
		t.Fatalf("create project file: %v", err)
	}

	return projectDir
}

type failingWriter struct {
	err error
}

func (writer failingWriter) Write([]byte) (int, error) {
	return 0, writer.err
}
