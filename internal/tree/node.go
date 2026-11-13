package tree

// EntryKind identifies whether a node is a directory, file, or symbolic link.
type EntryKind string

const (
	KindDirectory EntryKind = "directory"
	KindFile      EntryKind = "file"
	KindSymlink   EntryKind = "symlink"
)

// Node represents one entry in the scanned directory tree.
type Node struct {
	Name         string
	RelativePath string
	Kind         EntryKind
	Children     []*Node
}
