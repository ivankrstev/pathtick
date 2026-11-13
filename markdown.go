package main

import (
	"errors"
	"strings"
)

func renderMarkdown(root *Node, options Options) ([]byte, error) {
	if root == nil {
		return nil, errors.New("cannot render an empty tree")
	}

	var output strings.Builder

	// add title
	output.WriteString("# ")
	output.WriteString(options.Title)
	output.WriteByte('\n')
	output.WriteByte('\n')

	writeMarkdownNode(&output, root, 0)

	return []byte(output.String()), nil
}

func writeMarkdownNode(output *strings.Builder, node *Node, depth int) {
	output.WriteString(strings.Repeat("  ", depth))
	output.WriteString("- [ ] ")
	output.WriteString(node.Name)
	output.WriteByte('\n')

	for _, child := range node.Children {
		writeMarkdownNode(output, child, depth+1)
	}
}
