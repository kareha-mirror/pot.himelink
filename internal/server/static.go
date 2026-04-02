package server

import (
	"io"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"tea.kareha.org/pot/himelink/internal/config"
)

func handleStatic(cfg *config.Config, w http.ResponseWriter, r *http.Request) bool {
	baseDir, err := filepath.Abs(cfg.Site.Static)
	if err != nil {
		log.Fatalf("Invalid static dir: %v", err)
	}

	cleanPath := filepath.Clean(r.URL.Path)
	staticPath := filepath.Join(baseDir, cleanPath)
	if !strings.HasPrefix(staticPath, baseDir+string(os.PathSeparator)) {
		return false
	}

	st, err := os.Stat(staticPath)
	if err != nil || st.IsDir() {
		return false
	}

	f, err := os.Open(staticPath)
	if err != nil {
		http.Error(w, "static file read error", http.StatusInternalServerError)
		return true
	}
	defer f.Close()

	// 21600 secs = 6 hours
	w.Header().Set("Cache-Control", "public, max-age=21600")

	mimeType := mime.TypeByExtension(filepath.Ext(staticPath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", mimeType)

	io.Copy(w, f)
	return true
}
