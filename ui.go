package main

import "net/http"

func (s *Server) UI(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = w.Write([]byte("PayloadBox - UI not yet built into this binary.\n"))
}
