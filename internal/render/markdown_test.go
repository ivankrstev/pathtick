package render_test

import (
	"testing"

	"github.com/ivankrstev/pathtick/internal/render"
	"github.com/ivankrstev/pathtick/internal/tree"
)

func TestRenderMarkdownNestedTree(t *testing.T) {
	root := &tree.Node{
		Name: "repo",
		Children: []*tree.Node{
			{Name: "README.md"},
			{
				Name: "src",
				Children: []*tree.Node{
					{Name: "main.go"},
				},
			},
		},
	}

	got, err := render.RenderMarkdown(root, "Sample Project")
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	want := "# Sample Project\n\n" +
		"- [ ] repo\n" +
		"  - [ ] README\\.md\n" +
		"  - [ ] src\n" +
		"    - [ ] main\\.go\n"

	if string(got) != want {
		t.Errorf("RenderMarkdown() output:\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}

func TestRenderMarkdownNilRoot(t *testing.T) {
	_, err := render.RenderMarkdown(nil, "Sample Project")
	if err == nil {
		t.Fatal("RenderMarkdown(nil, ...) error = nil, want an error")
	}
}

func TestRenderMarkdownEscapesTitleAndFilenames(t *testing.T) {
	root := &tree.Node{
		Name: "repo",
		Children: []*tree.Node{
			{Name: "draft [todo] *v1*.md"},
			{Name: "résumé-日本語.md"},
			{Name: "line\nbreak.md"},
		},
	}

	got, err := render.RenderMarkdown(root, "a [blabla] **a** *1*")
	if err != nil {
		t.Fatalf("RenderMarkdown() error = %v", err)
	}

	want := "# a &#91;blabla&#93; \\*\\*a\\*\\* \\*1\\*\n\n" +
		"- [ ] repo\n" +
		"  - [ ] draft &#91;todo&#93; \\*v1\\*\\.md\n" +
		"  - [ ] résumé\\-日本語\\.md\n" +
		"  - [ ] line\\nbreak\\.md\n"
	if string(got) != want {
		t.Errorf("RenderMarkdown() output:\ngot:\n%q\nwant:\n%q", string(got), want)
	}
}
