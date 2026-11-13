package render

import (
	"errors"
	"strings"

	"github.com/ivankrstev/pathtick/internal/tree"
)

func RenderMarkdown(root *tree.Node, title string) ([]byte, error) {
	if root == nil {
		return nil, errors.New("cannot render an empty tree")
	}

	var output strings.Builder

	// add title
	output.WriteString("# ")
	output.WriteString(title)
	output.WriteByte('\n')
	output.WriteByte('\n')

	writeMarkdownNode(&output, root, 0)

	return []byte(output.String()), nil
}

func writeMarkdownNode(output *strings.Builder, node *tree.Node, depth int) {
	output.WriteString(strings.Repeat("  ", depth))
	output.WriteString("- [ ] ")
	output.WriteString(node.Name)
	output.WriteByte('\n')

	for _, child := range node.Children {
		writeMarkdownNode(output, child, depth+1)
	}
}
