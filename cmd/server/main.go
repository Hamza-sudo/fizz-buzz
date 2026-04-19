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
	defaultStatsDBDSN = "file:fizzbuzz_stats.db"
	defaultDBTimeout  = 200 * time.Millisecond
	shutdownTimeout   = 10 * time.Second
	readHeaderTimeout = 5 * time.Second
	readTimeout       = 10 * time.Second
	writeTimeout      = 15 * time.Second
	idleTimeout       = 60 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	statsStore, err := stats.NewSQLiteStoreWithTimeout(statsDBDSN(), statsDBTimeout(logger))
	if err != nil {
		logger.Error("failed to initialize stats store", "error", err)
		os.Exit(1)
	}
	defer func() {
		if closeErr := statsStore.Close(); closeErr != nil {
			logger.Error("failed to close stats store", "error", closeErr)
		}
	}()
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

func statsDBDSN() string {
	if dsn := os.Getenv("STATS_DB_DSN"); dsn != "" {
		return dsn
	}

	return defaultStatsDBDSN
}

func statsDBTimeout(logger *slog.Logger) time.Duration {
	value := os.Getenv("STATS_DB_TIMEOUT_MS")
	if value == "" {
		return defaultDBTimeout
	}

	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		logger.Warn("invalid STATS_DB_TIMEOUT_MS, falling back to default", "value", value, "default_ms", defaultDBTimeout.Milliseconds())
		return defaultDBTimeout
	}

	timeout := time.Duration(parsed) * time.Millisecond
	logger.Info("using custom stats db timeout", "timeout_ms", parsed)
	return timeout
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
