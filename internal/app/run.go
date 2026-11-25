package app

import (
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ivankrstev/pathtick/internal/output"
	"github.com/ivankrstev/pathtick/internal/render"
	"github.com/ivankrstev/pathtick/internal/scanner"
	"github.com/ivankrstev/pathtick/internal/tree"
)

// run scans the target, renders the requested output, and writes it to stdout or files.
func Run(options Options, stdout io.Writer) error {
	// Validate options and resolve output file paths before scanning the target directory.
	extensions, err := outputExtensions(options.Output.Format)
	if err != nil {
		return err
	}

	if !options.Output.Stdout {
		_, err := output.ResolveTargets(
			options.Output.Dir,
			options.Output.BaseName,
			options.Output.Overwrite,
			extensions,
		)
		if err != nil {
			return err
		}
	}

	fmt.Print("Scanning directory: ", options.TargetDir, "\n")

	root, err := scanner.ScanWithOptions(options.TargetDir, scanner.Options{
		MaxDepth:     options.Scan.Depth,
		SkipSymlinks: options.Scan.SkipSymlinks,
		SkipHidden:   options.Scan.SkipHidden,
	})
	if err != nil {
		return err
	}
	if err := sortTree(root, options.Display.Sort); err != nil {
		return err
	}

	if options.Display.Title == "" {
		options.Display.Title = root.Name
	}

	outputs, err := renderOutputs(root, options)
	if err != nil {
		return err
	}

	if options.Output.Stdout {
		if len(outputs) != 1 {
			return errors.New("--stdout requires exactly one rendered output")
		}

		rendered := outputs[0].Data
		_, err := stdout.Write(rendered)
		if err != nil {
			return fmt.Errorf("write to stdout: %w", err)
		}
		return nil
	}

	return output.WriteFiles(
		options.Output.Dir,
		options.Output.BaseName,
		options.Output.Overwrite,
		outputs,
	)
}

// sortTree orders each directory's children by name and applies the same order recursively.
func sortTree(root *tree.Node, order SortOrder) error {
	var less func(i, j int) bool
	switch order {
	case SortNone:
		return nil
	case SortAscending:
		less = func(i, j int) bool {
			return compareNames(root.Children[i].Name, root.Children[j].Name) < 0
		}
	case SortDescending:
		less = func(i, j int) bool {
			return compareNames(root.Children[i].Name, root.Children[j].Name) > 0
		}
	default:
		return fmt.Errorf(
			"unsupported --sort %q: supported values are %q, %q, and %q",
			order,
			SortNone,
			SortAscending,
			SortDescending,
		)
	}

	sort.SliceStable(root.Children, less)
	for _, child := range root.Children {
		if err := sortTree(child, order); err != nil {
			return err
		}
	}

	return nil
}

// compareNames compares names without regard to case, then uses the original
// spelling to keep names that differ only by case in a consistent order.
func compareNames(left, right string) int {
	if comparison := strings.Compare(strings.ToLower(left), strings.ToLower(right)); comparison != 0 {
		return comparison
	}

	return strings.Compare(left, right)
}

// renderOutputs renders the tree and returns each document with its file extension.
func renderOutputs(root *tree.Node, options Options) ([]output.File, error) {
	switch options.Output.Format {
	case FormatMarkdown:
		markdown, err := renderMarkdown(root, options.Display.Title)
		if err != nil {
			return nil, err
		}
		return []output.File{markdown}, nil

	case FormatHTML:
		html, err := renderHTML(root, options.Display.Title)
		if err != nil {
			return nil, err
		}
		return []output.File{html}, nil

	case FormatBoth:
		markdown, err := renderMarkdown(root, options.Display.Title)
		if err != nil {
			return nil, err
		}
		html, err := renderHTML(root, options.Display.Title)
		if err != nil {
			return nil, err
		}
		return []output.File{markdown, html}, nil

	default:
		return nil, fmt.Errorf("unsupported output format: %q", options.Output.Format)
	}
}

// outputExtensions returns the file extensions for a given output format.
func outputExtensions(format OutputFormat) ([]string, error) {
	switch format {
	case FormatMarkdown:
		return []string{"md"}, nil
	case FormatHTML:
		return []string{"html"}, nil
	case FormatBoth:
		return []string{"md", "html"}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %q", format)
	}
}

func renderMarkdown(root *tree.Node, title string) (output.File, error) {
	data, err := render.RenderMarkdown(root, title)
	if err != nil {
		return output.File{}, fmt.Errorf("render markdown: %w", err)
	}
	return output.File{Extension: "md", Data: data}, nil
}

func renderHTML(root *tree.Node, title string) (output.File, error) {
	data, err := render.RenderHTML(root, title)
	if err != nil {
		return output.File{}, fmt.Errorf("render HTML: %w", err)
	}
	return output.File{Extension: "html", Data: data}, nil
}
