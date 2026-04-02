package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"tea.kareha.org/pot/himelink/internal/config"
	"tea.kareha.org/pot/himelink/internal/gitea"
	"tea.kareha.org/pot/himelink/internal/github"
)

func RegisterRoutes(cfg *config.Config, r chi.Router) {
	for _, route := range cfg.Routes {
		if route.Protocol == "gitea" {
			path := "/" + route.Path + "/{owner}/{repo}"
			r.Get(path, gitea.RepoHandler(cfg, route))
			r.Get(path+"/", gitea.RepoHandler(cfg, route))
			r.Get(path+"/{mode}/*", gitea.PathHandler(cfg, route))
		} else if route.Protocol == "github" {
			path := "/" + route.Path + "/{owner}"
			r.Get(path, github.OwnerHandler(cfg, route, false))
			r.Get(path+"/", github.OwnerHandler(cfg, route, true))
			r.Get(path+"/{repo}", github.RepoHandler(cfg, route, false))
			r.Get(path+"/{repo}/", github.RepoHandler(cfg, route, true))
			r.Get(path+"/{repo}/{mode}/*", github.PathHandler(cfg, route))
		}
	}

	fileServer := http.FileServer(http.Dir(cfg.Site.Static))
	r.Handle("/*", fileServer)
}
