package app

import (
	"strings"
)

// This file holds the allowed values for the enum-like options.
//
// Each enum implements Set, String and Type, which is the interface cobra and
// pflag expect for Var/VarP flags (no import needed here), so bad values are
// rejected at parse time. Input is case-insensitive.

// ---------------------------------------------------------------------------
// OutputFormat
// ---------------------------------------------------------------------------

// OutputFormat identifies the format used for exported content.
type OutputFormat string

const (
	FormatMarkdown OutputFormat = "md"
	FormatHTML     OutputFormat = "html"
	FormatBoth     OutputFormat = "both"
)

// String returns the flag value, e.g. "md".
func (f OutputFormat) String() string { return string(f) }

// Type is the value name shown in --help.
func (f *OutputFormat) Type() string { return "md|html|both" }

// Set parses a --format value.
func (f *OutputFormat) Set(s string) error {
	*f = OutputFormat(strings.ToLower(strings.TrimSpace(s)))
	return nil
}

// Extension returns the file extension of a concrete format, or "" for
// FormatBoth (use Concrete first).
func (f OutputFormat) Extension() string {
	switch f {
	case FormatMarkdown:
		return ".md"
	case FormatHTML:
		return ".html"
	default:
		return ""
	}
}

// Concrete expands FormatBoth into the formats that actually get rendered.
func (f OutputFormat) Concrete() []OutputFormat {
	if f == FormatBoth {
		return []OutputFormat{FormatMarkdown, FormatHTML}
	}
	return []OutputFormat{f}
}

// ---------------------------------------------------------------------------
// SortOrder
// ---------------------------------------------------------------------------

// SortOrder controls the order in which directory entries are exported.
type SortOrder string

const (
	SortNone       SortOrder = "none" // walker order: unspecified, filesystem dependent
	SortAscending  SortOrder = "asc"
	SortDescending SortOrder = "desc"
)

// String returns the flag value, e.g. "asc".
func (s SortOrder) String() string { return string(s) }

// Type is the value name shown in --help.
func (s *SortOrder) Type() string { return "none|asc|desc" }

// Set parses a --sort value.
func (s *SortOrder) Set(in string) error {
	*s = SortOrder(strings.ToLower(strings.TrimSpace(in)))
	return nil
}
