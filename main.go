package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ByteFork/payloadbox/internal/logutil"
	"github.com/ByteFork/payloadbox/internal/middleware"
)

var (
	BuildTime string
	BuildSha  string
	Version   = "0.0.0"
)

// shutdownTimeout bounds how long the server waits for in-flight requests
// to complete after receiving SIGINT or SIGTERM before forcing close.
const shutdownTimeout = 5 * time.Second

func main() {
	settings := NewSettings()

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logutil.Level(settings.LogLevel),
	})))

	server := NewServer(*settings)

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", server.Health)
	mux.HandleFunc("/version", server.Version)
	mux.HandleFunc("/api/v1/settings", server.Settings)
	mux.Handle("/api/v1/history", middleware.Gzip(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			server.ClearRecords(w, r)
			return
		}

		server.ListRecords(w, r)
	})))
	mux.HandleFunc("/api/v1/events", server.Events)
	mux.HandleFunc("/", server.Record)

	srv := &http.Server{
		Addr:              server.settings.Address,
		Handler:           middleware.WithCORS(mux),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	slog.Info("starting server", "address", server.settings.Address)

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("failed to start server", "error", err.Error())
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	sig := <-quit

	slog.Warn("received shutdown signal", "signal", sig.String())

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server shutdown", "error", err.Error())
	}
}
