package ui

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func Attach(r *chi.Mux, distDir string) {
	fileServer := http.FileServer(http.Dir(distDir))

	r.Handle("/assets/*", http.StripPrefix("/", fileServer))
	r.Handle("/favicon.ico", fileServer)
	r.Handle("/robots.txt", fileServer)

	r.NotFound(func(w http.ResponseWriter, req *http.Request) {
		path := filepath.Join(distDir, filepath.Clean(req.URL.Path))
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, req)
			return
		}

		if strings.HasPrefix(req.URL.Path, "/actions") || strings.HasPrefix(req.URL.Path, "/public") {
			http.NotFound(w, req)
			return
		}

		indexPath := filepath.Join(distDir, "index.html")
		http.ServeFile(w, req, indexPath)
	})
}
