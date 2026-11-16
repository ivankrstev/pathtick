package app

import (
	"errors"
	"fmt"
	"io"

	"github.com/ivankrstev/pathtick/internal/output"
	"github.com/ivankrstev/pathtick/internal/render"
	"github.com/ivankrstev/pathtick/internal/scanner"
	"github.com/ivankrstev/pathtick/internal/tree"
)

// run scans the target, renders the requested output, and writes it to stdout or files.
func Run(options Options, stdout io.Writer) error {
	root, err := scanner.Scan(options.TargetDir)
	if err != nil {
		return err
	}

	if options.Title == "" {
		options.Title = root.Name
	}

	outputs, err := renderOutputs(root, options)
	if err != nil {
		return err
	}

	if options.Stdout {
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
		options.OutputDir,
		options.OutputBaseName,
		options.Overwrite,
		outputs,
	)
}

// renderOutputs renders the tree and returns each document with its file extension.
func renderOutputs(root *tree.Node, options Options) ([]output.File, error) {
	switch options.Format {
	case FormatMarkdown:
		markdown, err := renderMarkdown(root, options.Title)
		if err != nil {
			return nil, err
		}
		return []output.File{markdown}, nil

	case FormatHTML:
		html, err := renderHTML(root, options.Title)
		if err != nil {
			return nil, err
		}
		return []output.File{html}, nil

	case FormatBoth:
		markdown, err := renderMarkdown(root, options.Title)
		if err != nil {
			return nil, err
		}
		html, err := renderHTML(root, options.Title)
		if err != nil {
			return nil, err
		}
		return []output.File{markdown, html}, nil

	default:
		return nil, fmt.Errorf("unsupported output format: %q", options.Format)
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
