package scanner_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ivankrstev/pathtick/internal/scanner"
	"github.com/ivankrstev/pathtick/internal/tree"
)

func TestScanDirectoryNestedTree(t *testing.T) {
	targetDir := t.TempDir()
	nestedDir := filepath.Join(targetDir, "nested")
	emptyDir := filepath.Join(nestedDir, "empty")

	if err := os.MkdirAll(emptyDir, 0o755); err != nil {
		t.Fatalf("create nested directories: %v", err)
	}
	if err := os.WriteFile(filepath.Join(targetDir, "README.md"), []byte("readme"), 0o644); err != nil {
		t.Fatalf("create root file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nestedDir, "inside.txt"), []byte("inside"), 0o644); err != nil {
		t.Fatalf("create nested file: %v", err)
	}

	got, err := scanner.Scan(targetDir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}
	if got == nil {
		t.Fatal("Scan() returned a nil tree without an error")
	}

	want := map[string]struct {
		name string
		kind tree.EntryKind
	}{
		".": {
			name: filepath.Base(targetDir),
			kind: tree.KindDirectory,
		},
		"README.md": {
			name: "README.md",
			kind: tree.KindFile,
		},
		"nested": {
			name: "nested",
			kind: tree.KindDirectory,
		},
		"nested/inside.txt": {
			name: "inside.txt",
			kind: tree.KindFile,
		},
		"nested/empty": {
			name: "empty",
			kind: tree.KindDirectory,
		},
	}

	gotNodes := collectNodes(got)
	if len(gotNodes) != len(want) {
		t.Fatalf("Scan() returned %d nodes, want %d; paths = %v", len(gotNodes), len(want), nodePaths(gotNodes))
	}

	for relativePath, expected := range want {
		node, ok := gotNodes[relativePath]
		if !ok {
			t.Errorf("Scan() missing node at relative path %q", relativePath)
			continue
		}
		if node.Name != expected.name {
			t.Errorf("node %q name = %q, want %q", relativePath, node.Name, expected.name)
		}
		if node.Kind != expected.kind {
			t.Errorf("node %q kind = %q, want %q", relativePath, node.Kind, expected.kind)
		}
	}
}

func TestScanDirectoryRejectsMissingPath(t *testing.T) {
	missingPath := filepath.Join(t.TempDir(), "missing")

	_, err := scanner.Scan(missingPath)
	if err == nil {
		t.Fatal("Scan() error = nil, want an error for a missing path")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("Scan() error = %v, want an error wrapping os.ErrNotExist", err)
	}
}

func TestScanDirectoryRejectsFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0o644); err != nil {
		t.Fatalf("create file target: %v", err)
	}

	_, err := scanner.Scan(filePath)
	if err == nil {
		t.Fatal("Scan() error = nil, want an error for a file target")
	}
	if !strings.Contains(err.Error(), "is not a directory") {
		t.Errorf("Scan() error = %q, want it to explain that the target is not a directory", err)
	}
}

func TestScanDirectoryRejectsEmptyPath(t *testing.T) {
	_, err := scanner.Scan("")
	if err == nil {
		t.Fatal("Scan() error = nil, want an error for an empty path")
	}
	if !strings.Contains(err.Error(), "target directory must not be empty") {
		t.Errorf("Scan() error = %q, want an empty-target error", err)
	}
}

func TestScanDirectoryDoesNotFollowDirectorySymlink(t *testing.T) {
	targetDir := t.TempDir()
	actualDir := filepath.Join(targetDir, "actual")
	if err := os.MkdirAll(actualDir, 0o755); err != nil {
		t.Fatalf("create actual directory: %v", err)
	}
	if err := os.WriteFile(filepath.Join(actualDir, "inside.txt"), []byte("inside"), 0o644); err != nil {
		t.Fatalf("create file inside actual directory: %v", err)
	}

	linkPath := filepath.Join(targetDir, "linked")
	if err := os.Symlink(actualDir, linkPath); err != nil {
		t.Skipf("cannot create symlink in this environment: %v", err)
	}

	got, err := scanner.Scan(targetDir)
	if err != nil {
		t.Fatalf("Scan() error = %v", err)
	}

	nodes := collectNodes(got)
	link, ok := nodes["linked"]
	if !ok {
		t.Fatal("Scan() did not include the symlink")
	}
	if link.Kind != tree.KindSymlink {
		t.Errorf("symlink kind = %q, want %q", link.Kind, tree.KindSymlink)
	}
	if len(link.Children) != 0 {
		t.Errorf("symlink has %d children, want none", len(link.Children))
	}
	if _, ok := nodes["linked/inside.txt"]; ok {
		t.Error("Scan() traversed the directory symlink")
	}
}

func collectNodes(root *tree.Node) map[string]*tree.Node {
	nodes := make(map[string]*tree.Node)
	var visit func(*tree.Node)
	visit = func(node *tree.Node) {
		nodes[node.RelativePath] = node
		for _, child := range node.Children {
			visit(child)
		}
	}
	visit(root)
	return nodes
}

func nodePaths(nodes map[string]*tree.Node) []string {
	paths := make([]string, 0, len(nodes))
	for path := range nodes {
		paths = append(paths, path)
	}
	return paths
}
