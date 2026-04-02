package render

import (
	"net/http"

	"tea.kareha.org/pot/himelink/internal/config"
	"tea.kareha.org/pot/himelink/internal/templates"
)

type OwnerInfo struct {
	Login string
	URL   string
	Type  string
	Repos []RepoInfo
}

func RenderOwner(cfg *config.Config, w http.ResponseWriter, info OwnerInfo) {
	tmpl := templates.New("owner.html")

	tmpl.Execute(w, struct {
		SiteName string
		Title    string
		Info     OwnerInfo
	}{
		SiteName: cfg.Site.Name,
		Title:    info.Login,
		Info:     info,
	})
}
