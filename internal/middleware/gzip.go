package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"
)

var gzPool = sync.Pool{
	New: func() any {
		w, _ := gzip.NewWriterLevel(io.Discard, gzip.BestSpeed)
		return w
	},
}

type gzipResponseWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *gzipResponseWriter) WriteHeader(status int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(status)
}

func (w *gzipResponseWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func Gzip(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, ok := gzPool.Get().(*gzip.Writer)
		if !ok {
			// Pool returned an unexpected type; fall back to uncompressed.
			next.ServeHTTP(w, r)
			return
		}
		defer gzPool.Put(gz)

		w.Header().Set("Content-Encoding", "gzip")
		gz.Reset(w)

		defer func() { _ = gz.Close() }()

		next.ServeHTTP(&gzipResponseWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
