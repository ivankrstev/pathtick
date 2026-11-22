package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ivankrstev/pathtick/internal/cli"
)

func TestRootCmdStdout(t *testing.T) {
	projectDir := createProject(t)
	command := cli.NewRootCmd()
	var stdout strings.Builder
	command.SetOut(&stdout)
	command.SetArgs([]string{"--stdout", "-t", "CLI title", "-f", "md", projectDir})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	want := "# CLI title\n\n" +
		"- [ ] project\n" +
		"  - [ ] README\\.md\n"
	if got := stdout.String(); got != want {
		t.Errorf("command stdout:\ngot:\n%q\nwant:\n%q", got, want)
	}
}

func TestRootCmdDepthLimitsStdout(t *testing.T) {
	projectDir := createProject(t)
	nestedDir := filepath.Join(projectDir, "nested")
	nestedChildDir := filepath.Join(nestedDir, "child")
	if err := os.MkdirAll(nestedChildDir, 0o755); err != nil {
		t.Fatalf("create nested directories: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "inside.txt"), []byte("inside"), 0o644); err != nil {
		t.Fatalf("create nested file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedChildDir, "deep.txt"), []byte("deep"), 0o644); err != nil {
		t.Fatalf("create deep file: %v", err)
	}

	command := cli.NewRootCmd()
	var stdout strings.Builder
	command.SetOut(&stdout)
	command.SetArgs([]string{"--stdout", "--depth", "1", projectDir})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	for _, want := range []string{
		"- [ ] project",
		"README\\.md",
		"nested",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("command output does not contain %q:\n%s", want, stdout.String())
		}
	}

	for _, excluded := range []string{"inside.txt", "child", "deep.txt"} {
		if strings.Contains(stdout.String(), excluded) {
			t.Errorf("command output unexpectedly contains %q:\n%s", excluded, stdout.String())
		}
	}
}

func TestRootCmdSortFlag(t *testing.T) {
	projectDir := createProject(t)
	for name, contents := range map[string]string{
		"a":           "a",
		"PathTick.md": "PathTick",
		"zeta.txt":    "zeta",
	} {
		if err := os.WriteFile(filepath.Join(projectDir, name), []byte(contents), 0o644); err != nil {
			t.Fatalf("create file %q: %v", name, err)
		}
	}

	command := cli.NewRootCmd()
	var stdout strings.Builder
	command.SetOut(&stdout)
	command.SetArgs([]string{"--stdout", "--sort", "desc", projectDir})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	got := stdout.String()
	orderedLines := []string{
		"  - [ ] zeta\\.txt\n",
		"  - [ ] README\\.md\n",
		"  - [ ] PathTick\\.md\n",
		"  - [ ] a\n",
	}
	previousIndex := -1
	for _, line := range orderedLines {
		index := strings.Index(got, line)
		if index == -1 {
			t.Fatalf("command output does not contain %q:\n%s", line, got)
		}
		if index <= previousIndex {
			t.Errorf("command output does not order entries as %v:\n%s", orderedLines, got)
			break
		}
		previousIndex = index
	}
}

func TestRootCmdHTMLToStdout(t *testing.T) {
	projectDir := createProject(t)
	command := cli.NewRootCmd()
	var stdout strings.Builder
	command.SetOut(&stdout)
	command.SetArgs([]string{"--stdout", "--format", "html", projectDir})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	for _, want := range []string{
		"<!doctype html>",
		"<title>project</title>",
		"<input type=\"checkbox\"> README.md",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("command stdout does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestRootCmdHTMLFormat(t *testing.T) {
	projectDir := createProject(t)
	outputDir := t.TempDir()
	command := cli.NewRootCmd()
	command.SetArgs([]string{
		"--format", "html",
		"--output-dir", outputDir,
		"--output-name", "checklist",
		projectDir,
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	outputPath := filepath.Join(outputDir, "checklist.html")
	got, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("read HTML output: %v", err)
	}

	for _, want := range []string{
		"<!doctype html>",
		"<title>project</title>",
		"<input type=\"checkbox\"> README.md",
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("HTML output does not contain %q:\n%s", want, got)
		}
	}
}

func TestRootCmdBothFormat(t *testing.T) {
	projectDir := createProject(t)
	outputDir := t.TempDir()
	command := cli.NewRootCmd()
	command.SetArgs([]string{
		"--format", "both",
		"--output-dir", outputDir,
		"--output-name", "checklist",
		projectDir,
	})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	markdownPath := filepath.Join(outputDir, "checklist.md")
	markdown, err := os.ReadFile(markdownPath)
	if err != nil {
		t.Fatalf("read Markdown output: %v", err)
	}
	if !strings.Contains(string(markdown), "# project\n") {
		t.Errorf("Markdown output does not contain the project title:\n%s", markdown)
	}

	htmlPath := filepath.Join(outputDir, "checklist.html")
	html, err := os.ReadFile(htmlPath)
	if err != nil {
		t.Fatalf("read HTML output: %v", err)
	}
	if !strings.Contains(string(html), "<title>project</title>") {
		t.Errorf("HTML output does not contain the project title:\n%s", html)
	}
}

func TestRootCmdRejectsOutputFlagsWithStdout(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{
			name:    "output directory",
			args:    []string{"--stdout", "--output-dir", t.TempDir()},
			wantErr: "--stdout and --output-dir cannot be used together",
		},
		{
			name:    "output name",
			args:    []string{"--stdout", "--output-name", "custom"},
			wantErr: "--stdout and --output-name cannot be used together",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			command := cli.NewRootCmd()
			command.SetArgs(test.args)

			err := command.Execute()
			if err == nil {
				t.Fatalf("Execute() error = nil, want an error containing %q", test.wantErr)
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("Execute() error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}

func TestRootCmdRejectsStdoutWithOverwrite(t *testing.T) {
	command := cli.NewRootCmd()
	command.SetArgs([]string{"--stdout", "--overwrite"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want an error for stdout with overwrite")
	}
	if want := "cannot use both --stdout and --overwrite"; !strings.Contains(err.Error(), want) {
		t.Errorf("Execute() error = %q, want it to contain %q", err, want)
	}
}

func TestRootCmdRejectsStdoutWithBothFormats(t *testing.T) {
	command := cli.NewRootCmd()
	command.SetArgs([]string{"--stdout", "--format", "both"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want an error for stdout with both formats")
	}
	if want := "--stdout cannot be used with --format=both"; !strings.Contains(err.Error(), want) {
		t.Errorf("Execute() error = %q, want it to contain %q", err, want)
	}
}

func TestRootCmdRejectsUnsupportedFormat(t *testing.T) {
	command := cli.NewRootCmd()
	command.SetArgs([]string{"--format", "json"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want an error for an unsupported format")
	}
	if want := `unsupported --format "json"`; !strings.Contains(err.Error(), want) {
		t.Errorf("Execute() error = %q, want it to contain %q", err, want)
	}
}

func TestRootCmdHelpListsFlagGroups(t *testing.T) {
	command := cli.NewRootCmd()
	var stdout strings.Builder
	command.SetOut(&stdout)
	command.SetArgs([]string{"--help"})

	if err := command.Execute(); err != nil {
		t.Fatalf("Execute() with --help error = %v", err)
	}

	for _, want := range []string{
		"Arguments:",
		"Formatting:",
		"--format",
		"Output:",
		"--stdout",
		"Help:",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Errorf("help output does not contain %q:\n%s", want, stdout.String())
		}
	}
}

func TestRootCmdRejectsMultipleTargetDirectories(t *testing.T) {
	command := cli.NewRootCmd()
	command.SetArgs([]string{"first", "second"})

	if err := command.Execute(); err == nil {
		t.Fatal("Execute() error = nil, want an error for multiple target directories")
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
