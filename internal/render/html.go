package render

import (
	"bytes"
	_ "embed"
	"errors"
	"fmt"
	"html/template"
	"strings"

	"github.com/ivankrstev/pathtick/internal/tree"
)

//go:embed templates/checklist.html
var checklistHTML string

var checklistTemplate = template.Must(template.New("checklist.html").Funcs(template.FuncMap{
	"indent": func(depth int) string {
		return strings.Repeat("  ", depth)
	},
}).Parse(checklistHTML))

type htmlTemplateData struct {
	Title string
	Root  htmlTemplateNode
}

type htmlTemplateNode struct {
	Name         string
	RelativePath string
	Depth        int
	Children     []htmlTemplateNode
}

// RenderHTML renders the directory tree as a self-contained HTML checklist.
func RenderHTML(root *tree.Node, title string) ([]byte, error) {
	if root == nil {
		return nil, errors.New("cannot render an empty tree")
	}

	data := htmlTemplateData{
		Title: title,
		Root:  prepareHTMLNode(root, "", 3),
	}

	var output bytes.Buffer
	if err := checklistTemplate.Execute(&output, data); err != nil {
		return nil, fmt.Errorf("execute HTML template: %w", err)
	}

	return output.Bytes(), nil
}

func prepareHTMLNode(node *tree.Node, parentPath string, depth int) htmlTemplateNode {
	relativePath := node.RelativePath
	if relativePath == "" {
		switch parentPath {
		case "":
			relativePath = "."
		case ".":
			relativePath = node.Name
		default:
			relativePath = parentPath + "/" + node.Name
		}
	}

	prepared := htmlTemplateNode{
		Name:         node.Name,
		RelativePath: relativePath,
		Depth:        depth,
		Children:     make([]htmlTemplateNode, 0, len(node.Children)),
	}
	for _, child := range node.Children {
		prepared.Children = append(prepared.Children, prepareHTMLNode(child, relativePath, depth+2))
	}
	return prepared
}
