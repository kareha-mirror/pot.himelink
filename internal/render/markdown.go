package render

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/text"

	"tea.kareha.org/pot/himelink/internal/config"
	"tea.kareha.org/pot/himelink/internal/templates"
)

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM, // GitHub-like
	),
)

func markdownToHTML(filename string, src []byte) (string, []byte, error) {
	reader := text.NewReader(src)

	// Markdown -> AST
	doc := md.Parser().Parse(reader)

	title := extractTitle(filename, src, doc)

	var buf bytes.Buffer
	err := md.Renderer().Render(&buf, src, doc)
	if err != nil {
		return "", nil, err
	}
	return title, buf.Bytes(), nil
}

func RenderMarkdown(
	cfg *config.Config, w http.ResponseWriter,
	filename string,
	raw []byte,
) {
	title, htmlBytes, err := markdownToHTML(filename, raw)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	tmpl := templates.New("markdown.html")

	tmpl.Execute(w, struct {
		SiteName string
		Title    string
		Rendered template.HTML
	}{
		SiteName: cfg.Site.Name,
		Title:    title,
		Rendered: template.HTML(string(htmlBytes)),
	})
}
