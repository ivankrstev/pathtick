package render

import (
	"errors"
	"html"
	"strings"

	"github.com/ivankrstev/pathtick/internal/tree"
)

// RenderHTML renders the directory tree as a self-contained HTML checklist.
func RenderHTML(root *tree.Node, title string) ([]byte, error) {
	if root == nil {
		return nil, errors.New("cannot render an empty tree")
	}

	var output strings.Builder
	output.WriteString("<!doctype html>\n")
	output.WriteString("<html lang=\"en\">\n")
	output.WriteString("<head>\n")
	output.WriteString("  <meta charset=\"utf-8\">\n")
	output.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n")
	output.WriteString("  <title>")
	output.WriteString(html.EscapeString(title))
	output.WriteString("</title>\n")
	output.WriteString("  <style>\n")
	output.WriteString("    body { font-family: system-ui, sans-serif; line-height: 1.5; margin: 2rem auto; max-width: 52rem; padding: 0 1rem; }\n")
	output.WriteString("    ul { list-style: none; padding-left: 1.5rem; }\n")
	output.WriteString("    li { margin: 0.25rem 0; }\n")
	output.WriteString("    label { white-space: pre-wrap; }\n")
	output.WriteString("    input { margin-right: 0.5rem; }\n")
	output.WriteString("  </style>\n")
	output.WriteString("</head>\n")
	output.WriteString("<body>\n")
	output.WriteString("  <main>\n")
	output.WriteString("    <h1>")
	output.WriteString(html.EscapeString(title))
	output.WriteString("</h1>\n")
	output.WriteString("    <ul>\n")
	writeHTMLNode(&output, root, 3)
	output.WriteString("    </ul>\n")
	output.WriteString("  </main>\n")
	output.WriteString("</body>\n")
	output.WriteString("</html>\n")

	return []byte(output.String()), nil
}

func writeHTMLNode(output *strings.Builder, node *tree.Node, depth int) {
	indent := strings.Repeat("  ", depth)
	output.WriteString(indent)
	output.WriteString("<li><label><input type=\"checkbox\"> ")
	output.WriteString(html.EscapeString(node.Name))
	output.WriteString("</label>")

	if len(node.Children) == 0 {
		output.WriteString("</li>\n")
		return
	}

	output.WriteString("\n")
	output.WriteString(indent)
	output.WriteString("  <ul>\n")
	for _, child := range node.Children {
		writeHTMLNode(output, child, depth+2)
	}
	output.WriteString(indent)
	output.WriteString("  </ul>\n")
	output.WriteString(indent)
	output.WriteString("</li>\n")
}
