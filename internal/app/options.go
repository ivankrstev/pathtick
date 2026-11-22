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

// SortOrder controls the order in which directory entries are exported.
type SortOrder string

const (
	SortNone       SortOrder = "none"
	SortAscending  SortOrder = "asc"
	SortDescending SortOrder = "desc"
)

// Options contains the settings for one scan and export.
type Options struct {
	TargetDir      string
	Depth          int
	Title          string
	OutputDir      string
	OutputBaseName string
	Format         OutputFormat
	Sort           SortOrder
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
		Sort:           SortNone,
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

	switch options.Sort {
	case SortNone, SortAscending, SortDescending:
	default:
		return fmt.Errorf(
			"unsupported --sort %q: supported values are %q, %q, and %q",
			options.Sort,
			SortNone,
			SortAscending,
			SortDescending,
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
