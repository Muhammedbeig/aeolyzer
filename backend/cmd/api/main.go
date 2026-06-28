package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"aeolyzer/internal/httpapi"
	"aeolyzer/internal/intake"
	"aeolyzer/internal/interop"
	"aeolyzer/internal/observability"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/runtime"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	address := envOrDefault("AEOLYZER_ADDRESS", "127.0.0.1:8080")
	frontendOrigin := envOrDefault("AEOLYZER_FRONTEND_ORIGIN", "http://localhost:3000")
	// Process-bound ephemeral key avoids persisted credential management.
	// Limits blast radius of compromised keys to a single process lifecycle.
	signingKey := newSigningKey()

	intakeService := intake.NewService(newTraceID, signingKey, time.Now)
	orchestrator := orchestrator.NewService()
	connector := interop.NewSiteClient(8 * time.Second)
	executor := runtime.NewExecutor(net.DefaultResolver, connector, signingKey, time.Now)
	// Bounding the channel to 500 limits memory footprint during sudden telemetry bursts.
	// Dropping events on overflow is preferred over OOM cascading failures.
	events := observability.NewSink(500)
	handler := httpapi.NewHandler(
		intakeService,
		orchestrator,
		executor,
		events,
		logger,
		frontendOrigin,
	)

	server := &http.Server{
		Addr:    address,
		Handler: handler.Routes(),
		// Aggressive connection timeouts mitigate Slowloris attacks and socket descriptor exhaustion.
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Pre-allocating a buffered channel prevents goroutine leaks in the server error path.
	serverErrors := make(chan error, 1)
	go func() {
		logger.Info("AEOlyzer API listening", "address", address)
		serverErrors <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		if err != nil && err != http.ErrServerClosed {
			logger.Error("API stopped", "error", err)
			os.Exit(1)
		}
	case <-ctx.Done():
		shutdownContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownContext); err != nil {
			logger.Error("API shutdown failed", "error", err)
			os.Exit(1)
		}
	}
}

func newSigningKey() []byte {
	value := make([]byte, 32)
	// Panic on PRNG exhaustion. Operating with zero-entropy keys silently invalidates execution isolation.
	if _, err := rand.Read(value); err != nil {
		slog.Error("generate execution signing key", "error", err)
		os.Exit(1)
	}
	return value
}

func newTraceID() string {
	var value [16]byte
	// Graceful degradation on PRNG failure; prioritizing request throughput over trace fidelity.
	if _, err := rand.Read(value[:]); err != nil {
		return "trace-unavailable"
	}
	return hex.EncodeToString(value[:])
}

func envOrDefault(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}
