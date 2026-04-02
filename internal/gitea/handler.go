package gitea

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-chi/chi/v5"

	"tea.kareha.org/pot/himelink/internal/config"
	"tea.kareha.org/pot/himelink/internal/render"
)

type repoInfo struct {
	DefaultBranch string `json:"default_branch"`
}

func fetchRepoInfo(baseURL, owner, repo string) (repoInfo, error) {
	url := fmt.Sprintf("%s/api/v1/repos/%s/%s", baseURL, owner, repo)

	resp, err := http.Get(url)
	if err != nil {
		return repoInfo{}, err
	}
	defer resp.Body.Close()

	var info repoInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return repoInfo{}, err
	}

	return info, nil
}

func fetchRaw(baseURL, owner, repo, branch, path string) ([]byte, error) {
	url := fmt.Sprintf(
		"%s/api/v1/repos/%s/%s/raw/%s/%s",
		baseURL, owner, repo, branch, path,
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

var validName = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

func isValid(name string) bool {
	return validName.MatchString(name)
}

func handleRepo(
	cfg *config.Config,
	route config.Route,
	w http.ResponseWriter,
	r *http.Request,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}

	http.Redirect(w, r, fmt.Sprintf(
		"/%s/%s/%s/blob/README.md",
		route.Path,
		owner,
		repo,
	), http.StatusFound)
}

func RepoHandler(cfg *config.Config, route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handleRepo(cfg, route, w, r)
	}
}

func handlePath(
	cfg *config.Config,
	route config.Route,
	w http.ResponseWriter,
	r *http.Request,
) {
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	mode := chi.URLParam(r, "mode")
	path := chi.URLParam(r, "*")

	if !isValid(owner) || !isValid(repo) {
		http.Error(w, "invalid repo name", 400)
		return
	}
	if mode != "blob" {
		http.Error(w, "invalid mode", 400)
		return
	}

	if strings.Contains(path, "..") {
		http.Error(w, "invalid path", 500)
		return
	}
	if strings.HasPrefix(path, "/") {
		http.Error(w, "invalid path", 500)
		return
	}

	info, err := fetchRepoInfo(route.API, owner, repo)
	if err != nil {
		http.Error(w, "cannot get repo info", 500)
		return
	}
	branch := info.DefaultBranch
	if branch == "" {
		http.Error(w, "invalid branch", 500)
		return
	}

	raw, err := fetchRaw(
		route.API,
		owner,
		repo,
		branch,
		path,
	)
	if err != nil {
		http.Error(w, "cannot fetch", 500)
		return
	}

	ext := filepath.Ext(path)
	filename := filepath.Base(path)
	if ext == ".md" || ext == ".markdown" || ext == ".mdown" ||
		ext == ".mkd" || ext == ".mkdn" || ext == ".mdwn" ||
		ext == ".mdtxt" || ext == ".mdtext" {
		render.RenderMarkdown(cfg, w, filename, raw)
	} else {
		http.Error(w, "unsupported extension "+ext, 500)
		return
	}
}

func PathHandler(cfg *config.Config, route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handlePath(cfg, route, w, r)
	}
}
