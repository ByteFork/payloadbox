package ui

import (
	"net/http"
	"path/filepath"
	"strings"
)

var Handler http.Handler = http.HandlerFunc(serve)

func ShouldServe(path string) bool {
	return path == "/" || path == "/index.html" || strings.HasPrefix(path, "/assets/")
}

func serve(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	if path == "/" {
		path = "/index.html"
	}

	data, err := uiFS.ReadFile("dist" + path)
	if err != nil {
		if filepath.Ext(path) != "" {
			http.NotFound(w, r)
			return
		}

		data, err = uiFS.ReadFile("dist/index.html")
		if err != nil {
			http.Error(w, "UI missing", http.StatusInternalServerError)
			return
		}
	}

	switch {
	case strings.HasSuffix(path, ".js"):
		w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	case strings.HasSuffix(path, ".css"):
		w.Header().Set("Content-Type", "text/css; charset=utf-8")
	case strings.HasSuffix(path, ".html"):
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	case strings.HasSuffix(path, ".svg"):
		w.Header().Set("Content-Type", "image/svg+xml")
	default:
		if ct := http.DetectContentType(data); ct != "" {
			w.Header().Set("Content-Type", ct)
		}
	}

	_, _ = w.Write(data)
}
