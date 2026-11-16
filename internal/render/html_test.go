package render_test

import (
	"strings"
	"testing"

	"github.com/ivankrstev/pathtick/internal/render"
	"github.com/ivankrstev/pathtick/internal/tree"
)

func TestRenderHTMLNestedTree(t *testing.T) {
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

	got, err := render.RenderHTML(root, "Sample Project")
	if err != nil {
		t.Fatalf("RenderHTML() error = %v", err)
	}

	wantChecklist := "    <h1>Sample Project</h1>\n" +
		"    <ul>\n" +
		"      <li><label data-relative-path=\".\"><input type=\"checkbox\"> repo</label>\n" +
		"        <ul>\n" +
		"          <li><label data-relative-path=\"README.md\"><input type=\"checkbox\"> README.md</label></li>\n" +
		"          <li><label data-relative-path=\"src\"><input type=\"checkbox\"> src</label>\n" +
		"            <ul>\n" +
		"              <li><label data-relative-path=\"src/main.go\"><input type=\"checkbox\"> main.go</label></li>\n" +
		"            </ul>\n" +
		"          </li>\n" +
		"        </ul>\n" +
		"      </li>\n" +
		"    </ul>\n"
	if !strings.Contains(string(got), wantChecklist) {
		t.Errorf("RenderHTML() output does not contain the expected nested checklist:\n%s", got)
	}
}

func TestRenderHTMLNilRoot(t *testing.T) {
	_, err := render.RenderHTML(nil, "Sample Project")
	if err == nil {
		t.Fatal("RenderHTML(nil, ...) error = nil, want an error")
	}
}

func TestRenderHTMLEscapesTitleAndFileNames(t *testing.T) {
	root := &tree.Node{
		Name: `report <draft> & "notes" 'final'.md`,
	}
	title := `A <project> & "title" 'quoted'`

	got, err := render.RenderHTML(root, title)
	if err != nil {
		t.Fatalf("RenderHTML() error = %v", err)
	}

	escapedTitle := `A &lt;project&gt; &amp; &#34;title&#34; &#39;quoted&#39;`
	if count := strings.Count(string(got), escapedTitle); count != 2 {
		t.Errorf("escaped title appears %d times, want 2; output:\n%s", count, got)
	}

	escapedName := `report &lt;draft&gt; &amp; &#34;notes&#34; &#39;final&#39;.md`
	if !strings.Contains(string(got), escapedName) {
		t.Errorf("output does not contain escaped file name %q:\n%s", escapedName, got)
	}
}

func TestRenderHTMLProgressHelp(t *testing.T) {
	got, err := render.RenderHTML(&tree.Node{Name: "project"}, "Project")
	if err != nil {
		t.Fatalf("RenderHTML() error = %v", err)
	}

	for _, want := range []string{
		`<button id="progress-help-toggle" type="button" aria-label="Progress info" aria-expanded="false" aria-controls="progress-help-panel">?</button>`,
		`<aside id="progress-help-panel" aria-labelledby="progress-help-title" hidden>`,
		"Progress autosaves in this browser for this HTML file.",
		`.progress-help { position: fixed; right: 1rem; bottom: 1rem;`,
		`#progress-help-panel { position: absolute;`,
		`#progress-help-toggle:hover::after, #progress-help-toggle:focus-visible::after`,
		`font-size: 1.5rem;`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("RenderHTML() output does not contain %q:\n%s", want, got)
		}
	}
}

func TestRenderHTMLIncludesAutosaveScript(t *testing.T) {
	root := &tree.Node{
		Name:         "project",
		RelativePath: ".",
		Children: []*tree.Node{
			{Name: "README.md", RelativePath: "README.md"},
		},
	}

	got, err := render.RenderHTML(root, "Project")
	if err != nil {
		t.Fatalf("RenderHTML() error = %v", err)
	}

	for _, want := range []string{
		`<label data-relative-path="README.md"><input type="checkbox"> README.md</label>`,
		`window.localStorage.getItem(storageKey)`,
		`checkbox.checked = progress[path] === true;`,
		`checkbox.addEventListener("change", () => {`,
		`window.localStorage.setItem(storageKey, JSON.stringify(progress))`,
	} {
		if !strings.Contains(string(got), want) {
			t.Errorf("RenderHTML() output does not contain %q:\n%s", want, got)
		}
	}
}
