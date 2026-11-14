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
	output.WriteString(escapeMarkdownText(title))
	output.WriteByte('\n')
	output.WriteByte('\n')

	writeMarkdownNode(&output, root, 0)

	return []byte(output.String()), nil
}

func writeMarkdownNode(output *strings.Builder, node *tree.Node, depth int) {
	output.WriteString(strings.Repeat("  ", depth))
	output.WriteString("- [ ] ")
	output.WriteString(escapeMarkdownText(node.Name))
	output.WriteByte('\n')

	for _, child := range node.Children {
		writeMarkdownNode(output, child, depth+1)
	}
}

const markdownPunctuation = "!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~"

func escapeMarkdownText(text string) string {
	var escaped strings.Builder

	for _, r := range text {
		switch r {
		case '[':
			escaped.WriteString("&#91;")
		case ']':
			escaped.WriteString("&#93;")
		case '\n':
			escaped.WriteString(`\n`)
		case '\r':
			escaped.WriteString(`\r`)
		case '\t':
			escaped.WriteString(`\t`)
		default:
			if strings.ContainsRune(markdownPunctuation, r) {
				escaped.WriteByte('\\')
			}
			escaped.WriteRune(r)
		}
	}

	return escaped.String()
}
