package main

import (
	"errors"
	"fmt"
	"io"
)

// run scans the target, renders the requested output, and writes it to stdout or files.
func run(options Options, stdout io.Writer) error {
	tree, err := scanDirectory(options.TargetDir, options)
	if err != nil {
		return err
	}

	if options.Title == "" {
		options.Title = tree.Name
	}

	outputs, err := renderOutputs(tree, options)
	if err != nil {
		return err
	}

	if options.Stdout {
		if len(outputs) != 1 {
			return errors.New("--stdout requires exactly one rendered output")
		}

		rendered := outputs[0].data
		_, err := stdout.Write(rendered)
		if err != nil {
			return fmt.Errorf("write to stdout: %w", err)
		}
		return nil
	}

	return writeOutputFiles(
		options.OutputDir,
		options.OutputBaseName,
		options.Overwrite,
		outputs,
	)
}

// renderOutputs renders the tree and returns each document with its file extension.
func renderOutputs(root *Node, options Options) ([]renderedFile, error) {
	switch options.Format {
	case FormatMarkdown:
		data, err := renderMarkdown(root, options)
		if err != nil {
			return nil, fmt.Errorf("render markdown: %w", err)
		}
		return []renderedFile{{extension: "md", data: data}}, nil

	default:
		return nil, fmt.Errorf("unsupported output format: %q", options.Format)
	}
}
