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
		"      <li><label><input type=\"checkbox\"> repo</label>\n" +
		"        <ul>\n" +
		"          <li><label><input type=\"checkbox\"> README.md</label></li>\n" +
		"          <li><label><input type=\"checkbox\"> src</label>\n" +
		"            <ul>\n" +
		"              <li><label><input type=\"checkbox\"> main.go</label></li>\n" +
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
