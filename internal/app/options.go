package app

import (
	"errors"
	"fmt"
)

// Defaults. Flag registration should reuse these instead of repeating literals.
const (
	DefaultTargetDir = "."
	DefaultOutputDir = "."
	DefaultBaseName  = "PathTick"

	// UnlimitedDepth disables the depth limit.
	UnlimitedDepth = -1
)

// Options contains the settings for one scan and export.
type Options struct {
	TargetDir string // positional argument; defaults to "."

	Scan    ScanOptions
	Display DisplayOptions
	Output  OutputOptions
}

// ScanOptions controls the traversal of a target directory's entries.
type ScanOptions struct {
	// Depth is the number of levels below the target to visit.
	// 1 = direct children only, 0 = nothing (empty document),
	// UnlimitedDepth (-1) = no limit.
	Depth        int
	SkipSymlinks bool
	SkipHidden   bool // names starting with "."
}

// DisplayOptions controls the structure of the output document.
type DisplayOptions struct {
	Title    string // empty = derive from the target directory's name
	Sort     SortOrder
}

// OutputOptions controls the output format and destination.
type OutputOptions struct {
	Format    OutputFormat
	Dir       string
	BaseName  string // file name stem; the format supplies the extension
	Stdout    bool
	Overwrite bool
}

// DefaultOptions returns the default settings for a run.
func DefaultOptions() Options {
	return Options{
		TargetDir: DefaultTargetDir,
		Scan: ScanOptions{
			Depth: UnlimitedDepth,
		},
		Display: DisplayOptions{
			Sort: SortNone,
		},
		Output: OutputOptions{
			Format:   FormatMarkdown,
			Dir:      DefaultOutputDir,
			BaseName: DefaultBaseName,
		},
	}
}

// Validate checks option values and incompatible output modes.
func (options Options) Validate() error {
	if options.TargetDir == "" {
		return errors.New("target directory cannot be empty")
	}
	if options.Scan.Depth < -1 {
		return fmt.Errorf("--depth must be -1 or greater, got %d", options.Scan.Depth)
	}

	switch options.Output.Format {
	case FormatMarkdown, FormatHTML, FormatBoth:
	default:
		return fmt.Errorf(
			"unsupported --format %q: supported formats are %q, %q, and %q",
			options.Output.Format,
			FormatMarkdown,
			FormatHTML,
			FormatBoth,
		)
	}

	switch options.Display.Sort {
	case SortNone, SortAscending, SortDescending:
	default:
		return fmt.Errorf(
			"unsupported --sort %q: supported values are %q, %q, and %q",
			options.Display.Sort,
			SortNone,
			SortAscending,
			SortDescending,
		)
	}

	if options.Output.Stdout && options.Output.Format == FormatBoth {
		return errors.New("--stdout and --format=both cannot be used together; choose either md or html")
	}

	return nil
}
