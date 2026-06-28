package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	backendconfig "aeolyzer/config"
	"aeolyzer/internal/extensions"
	"aeolyzer/internal/httpapi"
	"aeolyzer/internal/intake"
	"aeolyzer/internal/interop"
	interopconfig "aeolyzer/internal/interop/config"
	"aeolyzer/internal/observability"
	observabilityconfig "aeolyzer/internal/observability/config"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/runtime"
	"aeolyzer/internal/skills"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := validateStartupContracts(); err != nil {
		logger.Error("startup contract validation failed", "error", err)
		os.Exit(1)
	}
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

func validateStartupContracts() error {
	if _, err := backendconfig.LoadEmbedded(); err != nil {
		return fmt.Errorf("validate Layer 2 and Layer 3 config: %w", err)
	}
	if _, err := skills.ValidateEmbeddedLibrary(); err != nil {
		return fmt.Errorf("validate Layer 4 skill library: %w", err)
	}
	if _, err := extensions.NewSchemas(); err != nil {
		return fmt.Errorf("validate Layer 5 schemas: %w", err)
	}
	if err := runtime.CompileSchemas(); err != nil {
		return fmt.Errorf("validate Layer 6 schemas: %w", err)
	}
	if _, err := interopconfig.Load(); err != nil {
		return fmt.Errorf("validate Layer 7 config: %w", err)
	}
	if _, err := observabilityconfig.LoadEmbeddedPolicies(); err != nil {
		return fmt.Errorf("validate Layer 8 policies: %w", err)
	}
	return nil
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
