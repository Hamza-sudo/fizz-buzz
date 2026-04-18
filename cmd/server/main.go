package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"fizz-buzz/internal/httpapi"
	"fizz-buzz/internal/stats"
)

const (
	defaultAddr       = ":8080"
	defaultMaxLimit   = 100000
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	statsStore := stats.NewStore()
	maxLimit := serverMaxLimit(logger)
	handler := httpapi.NewHandler(statsStore, maxLimit)
	addr := serverAddr()

	// Configure conservative defaults so the server behaves safely in production-like environments.
	server := &http.Server{
		Addr:              addr,
		Handler:           handler.Routes(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	go func() {
		logger.Info("server starting", "addr", server.Addr)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server failed", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	logger.Info("shutdown signal received")

	// Give in-flight requests a bounded amount of time to complete before exiting.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}

	logger.Info("server stopped")
}

func serverAddr() string {
	if port := os.Getenv("PORT"); port != "" {
		return ":" + port
	}

	// Default to a conventional local development port.
	return defaultAddr
}

func serverMaxLimit(logger *slog.Logger) int {
	value := os.Getenv("MAX_LIMIT")
	if value == "" {
		return defaultMaxLimit
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		logger.Warn("invalid MAX_LIMIT, falling back to default", "value", value, "default", defaultMaxLimit)
		return defaultMaxLimit
	}

	logger.Info("using custom max limit", "max_limit", parsed)
	return parsed
}
