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

func TestRootCmdRejectsUnsupportedFormat(t *testing.T) {
	command := cli.NewRootCmd()
	command.SetArgs([]string{"--format", "html"})

	err := command.Execute()
	if err == nil {
		t.Fatal("Execute() error = nil, want an error for an unsupported format")
	}
	if want := `unsupported --format "html"`; !strings.Contains(err.Error(), want) {
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
