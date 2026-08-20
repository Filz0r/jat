package api

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed all:webdist
var webDist embed.FS

func (s *Server) staticHandler() http.Handler {
	staticFS, err := fs.Sub(webDist, "webdist")
	if err != nil {
		// This should never happen with a valid embedded tree.
		s.logger.Printf("error creating static file sub-tree: %s", err)
		return http.NotFoundHandler()
	}
	fileServer := http.FileServer(http.FS(staticFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If the requested file exists, serve it; otherwise fall back to
		// index.html so React Router can handle the path.
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		if _, err := fs.Stat(staticFS, path[1:]); err != nil {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
