package scanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/ivankrstev/pathtick/internal/tree"
)

// scanError presents a concise message while preserving the underlying cause for error inspection.
type scanError struct {
	message string
	cause   error
}

func (err *scanError) Error() string { return err.message }

func (err *scanError) Unwrap() error { return err.cause }

func newScanError(action, path string, cause error) error {
	return &scanError{
		message: fmt.Sprintf("%s %q: %s", action, path, scanErrorReason(cause)),
		cause:   cause,
	}
}

// scanErrorReason returns a concise reason for the error, if known.
func scanErrorReason(err error) string {
	switch {
	case errors.Is(err, os.ErrNotExist):
		return "no such file or directory"
	case errors.Is(err, os.ErrPermission):
		return "permission denied"
	case errors.Is(err, syscall.ENOTDIR):
		return "not a directory"
	}

	if pathErr, ok := errors.AsType[*os.PathError](err); ok {
		return pathErr.Err.Error()
	}

	return err.Error()
}

// Scan builds a tree rooted at the provided directory.
func Scan(targetDir string) (*tree.Node, error) {
	return ScanWithDepth(targetDir, -1)
}

// ScanWithDepth builds a tree rooted at the provided directory, including entries
// through maxDepth. The root is at depth 0; -1 scans without a depth limit.
func ScanWithDepth(targetDir string, maxDepth int) (*tree.Node, error) {
	return ScanWithOptions(targetDir, Options{MaxDepth: maxDepth})
}

// Options controls traversal and filtering of a target directory's entries.
type Options struct {
	MaxDepth     int // The root is at depth 0. Use -1 for unlimited traversal.
	SkipSymlinks bool
	SkipHidden   bool
}

// ScanWithOptions builds a tree rooted at targetDir using the given traversal options.
func ScanWithOptions(targetDir string, options Options) (*tree.Node, error) {
	if options.MaxDepth < -1 {
		return nil, fmt.Errorf("maximum depth must be -1 or greater, got %d", options.MaxDepth)
	}

	if targetDir == "" {
		return nil, fmt.Errorf("target directory must not be empty")
	}

	absolutePath, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, newScanError("cannot resolve target directory", targetDir, err)
	}

	info, err := os.Stat(absolutePath)
	if err != nil {
		return nil, newScanError("cannot access target directory", absolutePath, err)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("target %q is not a directory", absolutePath)
	}

	children, err := scanChildren(absolutePath, ".", 0, options)
	if err != nil {
		return nil, err
	}

	return &tree.Node{
		Name:         filepath.Base(absolutePath),
		RelativePath: ".",
		Kind:         tree.KindDirectory,
		Children:     children,
	}, nil
}

// scanChildren reads a directory and recursively builds its child nodes.
// relativePath uses native filesystem separators during traversal.
func scanChildren(
	directoryPath,
	relativePath string,
	directoryDepth int,
	options Options,
) ([]*tree.Node, error) {
	if options.MaxDepth >= 0 && directoryDepth >= options.MaxDepth {
		return nil, nil
	}

	entries, err := readDirectoryEntries(directoryPath)
	if err != nil {
		return nil, err
	}

	children := make([]*tree.Node, 0, len(entries))

	for _, entry := range entries {
		name := entry.Name()
		if options.SkipHidden && strings.HasPrefix(name, ".") {
			continue
		}
		isSymlink := entry.Type()&os.ModeSymlink != 0
		if options.SkipSymlinks && isSymlink {
			continue
		}

		childRelativePath := filepath.Join(relativePath, name)

		node := &tree.Node{
			Name:         name,
			RelativePath: filepath.ToSlash(childRelativePath),
			Kind:         tree.KindFile,
		}

		switch {
		case isSymlink:
			node.Kind = tree.KindSymlink

		case entry.IsDir():
			node.Kind = tree.KindDirectory

			childPath := filepath.Join(directoryPath, name)

			node.Children, err = scanChildren(
				childPath,
				childRelativePath,
				directoryDepth+1,
				options,
			)

			if err != nil {
				return nil, err
			}
		}

		children = append(children, node)
	}

	return children, nil
}

// readDirectoryEntries preserves filesystem enumeration order.
// The directory handle closes before recursive scanning begins.
func readDirectoryEntries(directoryPath string) ([]os.DirEntry, error) {
	directory, err := os.Open(directoryPath)
	if err != nil {
		return nil, newScanError("cannot access directory", directoryPath, err)
	}
	defer directory.Close()

	entries, err := directory.ReadDir(-1)
	if err != nil {
		return nil, newScanError("cannot read directory", directoryPath, err)
	}

	return entries, nil
}
