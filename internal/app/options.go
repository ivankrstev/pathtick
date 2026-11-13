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
	}
}

// Validate checks option values and incompatible output modes.
func (options Options) Validate() error {
	if options.TargetDir == "" {
		return errors.New("target directory cannot be empty")
	}

	switch options.Format {
	case FormatMarkdown:
	default:
		return fmt.Errorf(
			"unsupported --format %q: currently only %q is supported",
			options.Format,
			FormatMarkdown,
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
