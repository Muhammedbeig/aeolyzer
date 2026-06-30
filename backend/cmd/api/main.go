package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"iter"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	backendconfig "aeolyzer/config"
	"aeolyzer/internal/extensions"
	a2aserver "aeolyzer/internal/extensions/a2a_server"
	"aeolyzer/internal/httpapi"
	"aeolyzer/internal/intake"
	"aeolyzer/internal/intake/contracts"
	"aeolyzer/internal/intake/middleware"
	"aeolyzer/internal/interop"
	a2atransport "aeolyzer/internal/interop/a2a_transport"
	interopconfig "aeolyzer/internal/interop/config"
	"aeolyzer/internal/observability"
	observabilityconfig "aeolyzer/internal/observability/config"
	"aeolyzer/internal/orchestrator"
	"aeolyzer/internal/runtime"
	"aeolyzer/internal/skills"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
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
	a2aServer, err := newA2AServer(envOrDefault("AEOLYZER_PUBLIC_BASE_URL", "https://localhost:8080"))
	if err != nil {
		logger.Error("A2A startup failed", "error", err)
		os.Exit(1)
	}
	routes := http.NewServeMux()
	a2aRoutes := a2aServer.Routes()
	routes.Handle(a2atransport.WellKnownAgentCardPath(), a2aRoutes)
	routes.Handle("/a2a", a2aRoutes)
	routes.Handle("/", handler.Routes())

	server := &http.Server{
		Addr:    address,
		Handler: routes,
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

func newA2AServer(publicBaseURL string) (*a2atransport.Server, error) {
	adkAgent, err := newPublicA2AAgent()
	if err != nil {
		return nil, fmt.Errorf("create public a2a agent: %w", err)
	}
	card, err := a2aserver.NewAgentCard(a2aserver.AgentCardConfig{
		Name:          "AEOlyzer",
		Description:   "Provides guarded website visibility, AEO audit, and content-planning capabilities through A2A.",
		PublicBaseURL: publicBaseURL,
		Skills:        a2aserver.DefaultPublicSkills(),
	})
	if err != nil {
		return nil, fmt.Errorf("build public a2a agent card: %w", err)
	}
	return a2atransport.NewServer(a2atransport.Config{
		Agent:    adkAgent,
		Card:     card,
		Verifier: a2aVerifierFromEnv(),
	})
}

func newPublicA2AAgent() (agent.Agent, error) {
	return agent.New(agent.Config{
		Name:        "public_site_guidance",
		Description: "Explains safe public website visibility and content-planning options.",
		Run: func(ctx agent.InvocationContext) iter.Seq2[*session.Event, error] {
			return func(yield func(*session.Event, error) bool) {
				text, err := publicA2AResponse(ctx.UserContent())
				if err != nil {
					yield(nil, err)
					return
				}
				event := session.NewEvent(ctx.InvocationID())
				event.Author = ctx.Agent().Name()
				event.LLMResponse = model.LLMResponse{
					Content: genai.NewContentFromText(text, genai.RoleModel),
				}
				yield(event, nil)
			}
		},
	})
}

func publicA2AResponse(content *genai.Content) (string, error) {
	input := contracts.SanitizedInput{RawText: extractA2AText(content)}
	if err := middleware.CheckForPromptInjection(input); err != nil {
		return "", err
	}
	if middleware.CheckForProtectedDisclosure(input) == contracts.DisclosureStatusDetected {
		return middleware.GuardOutboundResponse(
			"I can explain public AEOlyzer capabilities, but I cannot provide internal tools, workflows, traces, MCP endpoints, or policy details.",
			contracts.IntentProtectedDisclosure,
		)
	}
	return middleware.GuardOutboundResponse(
		"AEOlyzer can help authenticated agents request guarded website visibility, AEO audit, and content-planning support. Unsafe or internal-disclosure requests are blocked before workflow execution.",
		contracts.IntentCapabilityExplanation,
	)
}

func extractA2AText(content *genai.Content) string {
	if content == nil {
		return ""
	}
	var builder strings.Builder
	for _, part := range content.Parts {
		if part == nil || part.Text == "" {
			continue
		}
		if builder.Len() > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(part.Text)
	}
	return middleware.SanitizeContextValue(builder.String(), 4096)
}

func a2aVerifierFromEnv() a2atransport.TokenVerifier {
	if audience := os.Getenv("AEOLYZER_A2A_AUDIENCE"); audience != "" {
		return a2atransport.GoogleIDTokenVerifier{
			Audience:    audience,
			TenantClaim: os.Getenv("AEOLYZER_A2A_TENANT_CLAIM"),
		}
	}
	if token := os.Getenv("AEOLYZER_A2A_DEV_TOKEN"); token != "" {
		return a2atransport.StaticBearerVerifier{
			Token:   token,
			Subject: envOrDefault("AEOLYZER_A2A_DEV_SUBJECT", "local-a2a-caller"),
			Tenant:  os.Getenv("AEOLYZER_A2A_DEV_TENANT"),
		}
	}
	return a2atransport.RejectingVerifier{}
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
