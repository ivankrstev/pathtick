package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// renderedFile holds rendered content and the extension used for its output filename.
type renderedFile struct {
	extension string
	data      []byte
}

// writeOutputFiles writes the rendered files, refusing to replace existing files unless overwrite is true.
func writeOutputFiles(
	outputDir string,
	baseName string,
	overwrite bool,
	outputs []renderedFile,
) error {
	if len(outputs) == 0 {
		return errors.New("no output files to write")
	}
	if outputDir == "" {
		return errors.New("output directory must not be empty")
	}
	if baseName == "" || baseName == "." || baseName == ".." || filepath.Base(baseName) != baseName {
		return fmt.Errorf("invalid output file name %q", baseName)
	}

	if err := os.MkdirAll(outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory %q: %w", outputDir, err)
	}

	targets := make([]string, 0, len(outputs))
	seenTargets := make(map[string]struct{}, len(outputs))

	for _, output := range outputs {
		extension := strings.TrimPrefix(output.extension, ".")
		if extension == "" || filepath.Base(extension) != extension {
			return fmt.Errorf("invalid output file extension %q", output.extension)
		}

		path := filepath.Join(outputDir, baseName+"."+extension)
		if _, exists := seenTargets[path]; exists {
			return fmt.Errorf("duplicate output file %q", path)
		}
		seenTargets[path] = struct{}{}

		if !overwrite {
			_, err := os.Lstat(path)
			switch {
			case err == nil:
				return fmt.Errorf("output file %q already exists; use --overwrite to replace it", path)
			case !errors.Is(err, os.ErrNotExist):
				return fmt.Errorf("check output file %q: %w", path, err)
			}
		}

		targets = append(targets, path)
	}

	createdPaths := make([]string, 0, len(targets))
	for index, output := range outputs {
		if err := writeOutputFile(targets[index], output.data, overwrite); err != nil {
			if !overwrite {
				for _, path := range createdPaths {
					_ = os.Remove(path)
				}
			}
			return err
		}
		if !overwrite {
			createdPaths = append(createdPaths, targets[index])
		}
	}

	return nil
}

// writeOutputFile writes a single output file, refusing to replace an existing file unless overwrite is true.
func writeOutputFile(path string, data []byte, overwrite bool) error {
	flags := os.O_WRONLY | os.O_CREATE
	if overwrite {
		flags |= os.O_TRUNC
	} else {
		flags |= os.O_EXCL
	}

	file, err := os.OpenFile(path, flags, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("output file %q already exists; use --overwrite to replace it", path)
		}
		return fmt.Errorf("create output file %q: %w", path, err)
	}

	n, writeErr := file.Write(data)
	if writeErr == nil && n != len(data) {
		writeErr = io.ErrShortWrite
	}
	if writeErr == nil {
		writeErr = file.Sync()
	}
	closeErr := file.Close()

	if writeErr != nil {
		if !overwrite {
			_ = os.Remove(path)
		}
		return fmt.Errorf("write output file %q: %w", path, writeErr)
	}
	if closeErr != nil {
		if !overwrite {
			_ = os.Remove(path)
		}
		return fmt.Errorf("close output file %q: %w", path, closeErr)
	}

	return nil
}
