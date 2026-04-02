package render

import (
	"net/http"

	"tea.kareha.org/pot/himelink/internal/config"
	"tea.kareha.org/pot/himelink/internal/templates"
)

type RepoInfo struct {
	Description string
	Name        string
	ReadmeName  string
	ReadmePath  string
	URL         string
	OwnerName   string
	OwnerPath   string
}

func RenderRepo(cfg *config.Config, w http.ResponseWriter, info RepoInfo) {
	tmpl := templates.New("repo.html")

	tmpl.Execute(w, struct {
		SiteName string
		Title    string
		Info     RepoInfo
	}{
		SiteName: cfg.Site.Name,
		Title:    info.Name,
		Info:     info,
	})
}
