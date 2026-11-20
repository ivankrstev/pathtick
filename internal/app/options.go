package app

import (
	"errors"
	"fmt"
)

// OutputFormat identifies the format used for exported content.
type OutputFormat string

const (
	FormatMarkdown OutputFormat = "md"
	FormatHTML     OutputFormat = "html"
	FormatBoth     OutputFormat = "both"
)

// Options contains the settings for one scan and export.
type Options struct {
	TargetDir      string
	Depth          int
	Title          string
	OutputDir      string
	OutputBaseName string
	Format         OutputFormat
	Stdout         bool
	Overwrite      bool
}

// DefaultOptions returns the default settings for a run.
func DefaultOptions() Options {
	return Options{
		TargetDir:      ".",
		Title:          "", // empty means derive from the target directory
		OutputDir:      ".",
		OutputBaseName: "PathTick",
		Format:         FormatMarkdown,
		Stdout:         false,
		Overwrite:      false,
		Depth:          -1,
	}
}

// Validate checks option values and incompatible output modes.
func (options Options) Validate() error {
	if options.TargetDir == "" {
		return errors.New("target directory cannot be empty")
	}
	if options.Depth < -1 {
		return fmt.Errorf("--depth must be -1 or greater, got %d", options.Depth)
	}

	switch options.Format {
	case FormatMarkdown, FormatHTML, FormatBoth:
	default:
		return fmt.Errorf(
			"unsupported --format %q: supported formats are %q, %q, and %q",
			options.Format,
			FormatMarkdown,
			FormatHTML,
			FormatBoth,
		)
	}

	if options.Stdout && options.Format == FormatBoth {
		return errors.New("--stdout cannot be used with --format=both; choose either md or html")
	}

	if options.Stdout && options.Overwrite {
		return errors.New("cannot use both --stdout and --overwrite")
	}

	return nil
}
