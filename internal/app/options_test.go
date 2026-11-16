package app_test

import (
	"strings"
	"testing"

	"github.com/ivankrstev/pathtick/internal/app"
)

func TestDefaultOptions(t *testing.T) {
	want := app.Options{
		TargetDir:      ".",
		Title:          "",
		OutputDir:      ".",
		OutputBaseName: "PathTick",
		Format:         app.FormatMarkdown,
		Stdout:         false,
		Overwrite:      false,
	}

	got := app.DefaultOptions()
	if got != want {
		t.Errorf("DefaultOptions() = %#v, want %#v", got, want)
	}
}

func TestOptionsValidate(t *testing.T) {
	validDefaults := app.DefaultOptions()
	validStdout := validDefaults
	validStdout.Stdout = true
	validOverwrite := validDefaults
	validOverwrite.Overwrite = true
	validHTML := validDefaults
	validHTML.Format = app.FormatHTML
	validBoth := validDefaults
	validBoth.Format = app.FormatBoth

	emptyTarget := validDefaults
	emptyTarget.TargetDir = ""
	unsupportedUnknown := validDefaults
	unsupportedUnknown.Format = "json"
	stdoutAndOverwrite := validDefaults
	stdoutAndOverwrite.Stdout = true
	stdoutAndOverwrite.Overwrite = true

	tests := []struct {
		name    string
		options app.Options
		wantErr string
	}{
		{
			name:    "default markdown options are valid",
			options: validDefaults,
		},
		{
			name:    "markdown to stdout is valid",
			options: validStdout,
		},
		{
			name:    "markdown overwrite is valid",
			options: validOverwrite,
		},
		{
			name:    "empty target directory",
			options: emptyTarget,
			wantErr: "target directory cannot be empty",
		},
		{
			name:    "HTML format",
			options: validHTML,
		},
		{
			name:    "both formats",
			options: validBoth,
		},
		{
			name:    "unknown format",
			options: unsupportedUnknown,
			wantErr: "unsupported --format",
		},
		{
			name:    "stdout with overwrite",
			options: stdoutAndOverwrite,
			wantErr: "cannot use both --stdout and --overwrite",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.options.Validate()
			if test.wantErr == "" {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			if err == nil {
				t.Fatalf("Validate() error = nil, want an error containing %q", test.wantErr)
			}
			if !strings.Contains(err.Error(), test.wantErr) {
				t.Errorf("Validate() error = %q, want it to contain %q", err, test.wantErr)
			}
		})
	}
}
