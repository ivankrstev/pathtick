package render

import (
	"errors"
	"html"
	"strings"

	"github.com/ivankrstev/pathtick/internal/tree"
)

const htmlPersistenceScript = `  <script>
    (() => {
      const storageKey = "pathtick:checklist:" + window.location.pathname;
      const status = document.getElementById("save-status");
      const labels = document.querySelectorAll("label[data-relative-path]");
      const progress = Object.create(null);
      const helpToggle = document.getElementById("progress-help-toggle");
      const helpPanel = document.getElementById("progress-help-panel");
      const helpTitle = document.getElementById("progress-help-title");
      const helpClose = document.getElementById("progress-help-close");
      let storageAvailable = true;

      const closeHelp = () => {
        helpPanel.hidden = true;
        helpToggle.setAttribute("aria-expanded", "false");
      };

      helpToggle.addEventListener("click", () => {
        const open = helpToggle.getAttribute("aria-expanded") !== "true";
        helpPanel.hidden = !open;
        helpToggle.setAttribute("aria-expanded", String(open));
        if (open) {
          helpTitle.focus();
        }
      });
      helpClose.addEventListener("click", () => {
        closeHelp();
        helpToggle.focus();
      });
      document.addEventListener("keydown", (event) => {
        if (event.key === "Escape" && !helpPanel.hidden) {
          closeHelp();
          helpToggle.focus();
        }
      });

      try {
        const saved = window.localStorage.getItem(storageKey);
        if (saved) {
          try {
            const parsed = JSON.parse(saved);
            if (parsed && typeof parsed === "object" && !Array.isArray(parsed)) {
              for (const [path, checked] of Object.entries(parsed)) {
                progress[path] = checked === true;
              }
            }
          } catch {}
        }
      } catch {
        storageAvailable = false;
      }

      if (!storageAvailable) {
        status.textContent = "Browser storage is unavailable; progress will not persist.";
      }

      for (const label of labels) {
        const checkbox = label.querySelector("input[type=checkbox]");
        const path = label.dataset.relativePath;
        checkbox.checked = progress[path] === true;
        checkbox.addEventListener("change", () => {
          progress[path] = checkbox.checked;
          try {
            window.localStorage.setItem(storageKey, JSON.stringify(progress));
            status.textContent = "Progress saved in this browser.";
          } catch {
            status.textContent = "Browser storage is unavailable; progress will not persist.";
          }
        });
      }
    })();
  </script>
`

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
	output.WriteString("    body { font-family: system-ui, sans-serif; line-height: 1.5; margin: 2rem auto; max-width: 52rem; padding: 0 1rem 5rem; }\n")
	output.WriteString("    ul { list-style: none; padding-left: 1.5rem; }\n")
	output.WriteString("    li { margin: 0.25rem 0; }\n")
	output.WriteString("    label { white-space: pre-wrap; }\n")
	output.WriteString("    input { margin-right: 0.5rem; }\n")
	output.WriteString("    #save-status { color: #555; font-size: 0.9rem; }\n")
	output.WriteString("    .progress-help { position: fixed; right: 1rem; bottom: 1rem; z-index: 1000; }\n")
	output.WriteString("    #progress-help-toggle { position: relative; display: grid; place-items: center; width: 2.75rem; height: 2.75rem; padding: 0; cursor: pointer; border: 1px solid #555; border-radius: 50%; background: #fff; color: #222; font: inherit; font-size: 1.5rem; font-weight: 700; box-shadow: 0 0.2rem 0.8rem rgba(0, 0, 0, 0.2); }\n")
	output.WriteString("    #progress-help-toggle::after { content: attr(aria-label); position: absolute; right: calc(100% + 0.65rem); bottom: 50%; transform: translateY(50%); visibility: hidden; opacity: 0; pointer-events: none; white-space: nowrap; border-radius: 0.3rem; background: #222; color: #fff; padding: 0.4rem 0.6rem; font-size: 0.85rem; font-weight: 400; transition: opacity 0.1s ease; }\n")
	output.WriteString("    #progress-help-toggle:hover::after, #progress-help-toggle:focus-visible::after { visibility: visible; opacity: 1; }\n")
	output.WriteString("    #progress-help-panel { position: absolute; right: 0; bottom: calc(100% + 0.5rem); box-sizing: border-box; width: min(22rem, calc(100vw - 2rem)); max-height: min(60vh, 24rem); overflow: auto; border: 1px solid #aaa; border-radius: 0.5rem; background: #fff; color: #222; padding: 1rem; box-shadow: 0 0.5rem 1.5rem rgba(0, 0, 0, 0.25); }\n")
	output.WriteString("    .progress-help-header { display: flex; align-items: start; justify-content: space-between; gap: 1rem; }\n")
	output.WriteString("    .progress-help-header h2 { margin: 0; font-size: 1rem; }\n")
	output.WriteString("    #progress-help-panel p { margin-bottom: 0; }\n")
	output.WriteString("  </style>\n")
	output.WriteString("</head>\n")
	output.WriteString("<body>\n")
	output.WriteString("  <main>\n")
	output.WriteString("    <h1>")
	output.WriteString(html.EscapeString(title))
	output.WriteString("</h1>\n")
	output.WriteString("    <ul>\n")
	writeHTMLNode(&output, root, 3, "")
	output.WriteString("    </ul>\n")
	output.WriteString("    <p id=\"save-status\" role=\"status\" aria-live=\"polite\">Progress saves in this browser when storage is available.</p>\n")
	output.WriteString("  </main>\n")
	output.WriteString("  <div class=\"progress-help\">\n")
	output.WriteString("    <button id=\"progress-help-toggle\" type=\"button\" aria-label=\"Progress info\" aria-expanded=\"false\" aria-controls=\"progress-help-panel\">?</button>\n")
	output.WriteString("    <aside id=\"progress-help-panel\" aria-labelledby=\"progress-help-title\" hidden>\n")
	output.WriteString("      <div class=\"progress-help-header\">\n")
	output.WriteString("        <h2 id=\"progress-help-title\" tabindex=\"-1\">About saved progress</h2>\n")
	output.WriteString("        <button id=\"progress-help-close\" type=\"button\">Close</button>\n")
	output.WriteString("      </div>\n")
	output.WriteString("      <p>Progress is saved in this browser for this HTML file. The HTML file itself is not changed. Opening it in another browser or moving it to a different location may show a separate checklist.</p>\n")
	output.WriteString("    </aside>\n")
	output.WriteString("  </div>\n")
	output.WriteString(htmlPersistenceScript)
	output.WriteString("</body>\n")
	output.WriteString("</html>\n")

	return []byte(output.String()), nil
}

func writeHTMLNode(output *strings.Builder, node *tree.Node, depth int, parentPath string) {
	indent := strings.Repeat("  ", depth)
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

	output.WriteString(indent)
	output.WriteString("<li><label data-relative-path=\"")
	output.WriteString(html.EscapeString(relativePath))
	output.WriteString("\"><input type=\"checkbox\"> ")
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
		writeHTMLNode(output, child, depth+2, relativePath)
	}
	output.WriteString(indent)
	output.WriteString("  </ul>\n")
	output.WriteString(indent)
	output.WriteString("</li>\n")
}
