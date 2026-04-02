package render

import (
	"net/http"

	"tea.kareha.org/pot/himelink/internal/config"
	"tea.kareha.org/pot/himelink/internal/templates"
)

func RenderText(
	cfg *config.Config, w http.ResponseWriter,
	filename string,
	raw []byte,
) {
	tmpl := templates.New("text.html")

	tmpl.Execute(w, struct {
		SiteName string
		Title    string
		Text     string
	}{
		SiteName: cfg.Site.Name,
		Title:    filename,
		Text:     string(raw),
	})
}
