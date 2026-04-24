// Package middleware houses HTTP middleware used by the server.
package middleware

import "net/http"

const (
	corsAllowedMethods = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	corsAllowedHeaders = "Content-Type, Authorization, X-Requested-With"
	corsExposedHeaders = "Content-Type"
)

func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == "" {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")
		w.Header().Add("Vary", "Access-Control-Request-Headers")
		w.Header().Set("Access-Control-Allow-Methods", corsAllowedMethods)
		w.Header().Set("Access-Control-Allow-Headers", corsAllowedHeaders)
		w.Header().Set("Access-Control-Expose-Headers", corsExposedHeaders)

		if r.Method == http.MethodOptions &&
			r.Header.Get("Access-Control-Request-Method") != "" {
			w.Header().Set("Access-Control-Max-Age", "600")
			w.WriteHeader(http.StatusNoContent)

			return
		}

		next.ServeHTTP(w, r)
	})
}
