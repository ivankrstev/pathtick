package app_test

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/ivankrstev/pathtick/internal/app"
	"github.com/ivankrstev/pathtick/internal/cli"
)

func TestDefaultOptions(t *testing.T) {
	want := app.Options{
		TargetDir: ".",
		Scan: app.ScanOptions{
			Depth: -1,
		},
		Output: app.OutputOptions{
			Dir:       ".",
			BaseName:  "PathTick",
			Format:    app.FormatMarkdown,
			Stdout:    false,
			Overwrite: false,
		},
		Display: app.DisplayOptions{
			Sort:  app.SortNone,
			Title: "",
		},
		// Filter: app.FilterOptions{
		// 	Only: app.EntryAll,
		// },
	}

	got := app.DefaultOptions()
	if !reflect.DeepEqual(got, want) {
		t.Errorf("DefaultOptions() = %+v, want %+v", got, want)
	}
}

func TestOptionsValidate(t *testing.T) {
	validDefaults := app.DefaultOptions()
	validDepthZero := validDefaults
	validDepthZero.Scan.Depth = 0
	validStdout := validDefaults
	validStdout.Output.Stdout = true
	validOverwrite := validDefaults
	validOverwrite.Output.Overwrite = true
	validHTML := validDefaults
	validHTML.Output.Format = app.FormatHTML
	validBoth := validDefaults
	validBoth.Output.Format = app.FormatBoth
	validSortAscending := validDefaults
	validSortAscending.Display.Sort = app.SortAscending
	validSortDescending := validDefaults
	validSortDescending.Display.Sort = app.SortDescending

	emptyTarget := validDefaults
	emptyTarget.TargetDir = ""
	unsupportedUnknown := validDefaults
	unsupportedUnknown.Output.Format = "json"
	unsupportedSort := validDefaults
	unsupportedSort.Display.Sort = "random"
	invalidDepth := validDefaults
	invalidDepth.Scan.Depth = -2
	stdoutAndOverwrite := validDefaults
	stdoutAndOverwrite.Output.Stdout = true
	stdoutAndOverwrite.Output.Overwrite = true

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
			name:    "depth zero is valid",
			options: validDepthZero,
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
			name:    "ascending sort",
			options: validSortAscending,
		},
		{
			name:    "descending sort",
			options: validSortDescending,
		},
		{
			name:    "unknown format",
			options: unsupportedUnknown,
			wantErr: "unsupported --format",
		},
		{
			name:    "unknown sort order",
			options: unsupportedSort,
			wantErr: "unsupported --sort",
		},
		{
			name:    "depth below unlimited",
			options: invalidDepth,
			wantErr: "--depth must be -1 or greater",
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

func TestRootRejectsStdoutWithFileOutputFlags(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{"output-dir", []string{"--stdout", "--output-dir", "out"}, "--stdout and --output-dir cannot be used together"},
		{"output-name", []string{"--stdout", "--output-name", "list"}, "--stdout and --output-name cannot be used together"},
		{"overwrite", []string{"--stdout", "--overwrite"}, "--stdout and --overwrite cannot be used together"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := cli.NewRootCmd()
			cmd.SetArgs(tt.args)
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			err := cmd.Execute()
			if err == nil {
				t.Fatalf("Execute() error = nil, want an error containing %q", tt.want)
			}
			if !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Execute() error = %q, want it to contain %q", err, tt.want)
			}
		})
	}
}
