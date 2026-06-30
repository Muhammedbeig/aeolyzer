// Package a2atransport mounts official A2A HTTP handlers for ADK agents.
package a2atransport

import (
	"context"
	"crypto/subtle"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"aeolyzer/internal/extensions/a2a_server"

	"github.com/a2aproject/a2a-go/v2/a2a"
	"github.com/a2aproject/a2a-go/v2/a2asrv"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/runner"
	"google.golang.org/adk/server/adka2a/v2"
	"google.golang.org/adk/session"
	"google.golang.org/api/idtoken"
)

const (
	defaultJSONRPCPath       = "/a2a"
	defaultMaxRequestBytes   = 128 << 10
	defaultInactivityTimeout = 30 * time.Second
)

// Identity is the authenticated A2A caller identity safe for downstream policy.
type Identity struct {
	Subject string
	Tenant  string
}

// TokenVerifier validates bearer tokens for non-card A2A calls.
type TokenVerifier interface {
	VerifyBearer(ctx context.Context, token string) (Identity, error)
}

// Config contains the dependencies required to mount the A2A transport.
type Config struct {
	Agent              agent.Agent
	Card               *a2a.AgentCard
	Verifier           TokenVerifier
	JSONRPCPath        string
	MaxRequestBytes    int64
	InactivityTimeout  time.Duration
	SessionService     session.Service
	AuthenticatedAttrs map[string]any
}

// Server exposes the public Agent Card and authenticated A2A JSON-RPC endpoint.
type Server struct {
	cardHandler http.Handler
	rpcHandler  http.Handler
	rpcPath     string
	maxBytes    int64
}

// NewServer creates an A2A transport mount backed by Google ADK.
func NewServer(config Config) (*Server, error) {
	if config.Agent == nil {
		return nil, errors.New("a2a transport requires an adk agent")
	}
	if config.Verifier == nil {
		return nil, errors.New("a2a transport requires a token verifier")
	}
	if err := a2aserver.ValidateAgentCard(config.Card); err != nil {
		return nil, fmt.Errorf("validate a2a agent card: %w", err)
	}
	rpcPath, err := cleanJSONRPCPath(config.JSONRPCPath)
	if err != nil {
		return nil, err
	}
	maxBytes := config.MaxRequestBytes
	if maxBytes == 0 {
		maxBytes = defaultMaxRequestBytes
	}
	if maxBytes < 1 || maxBytes > 1<<20 {
		return nil, errors.New("a2a max request bytes is invalid")
	}
	timeout := config.InactivityTimeout
	if timeout == 0 {
		timeout = defaultInactivityTimeout
	}
	if timeout < time.Second || timeout > 2*time.Minute {
		return nil, errors.New("a2a inactivity timeout is invalid")
	}
	sessionService := config.SessionService
	if sessionService == nil {
		sessionService = session.InMemoryService()
	}

	executor := adka2a.NewExecutor(adka2a.ExecutorConfig{
		RunnerConfig: runner.Config{
			AppName:        config.Agent.Name(),
			Agent:          config.Agent,
			SessionService: sessionService,
		},
	})
	requestHandler := a2asrv.NewHandler(
		executor,
		a2asrv.WithCapabilityChecks(&config.Card.Capabilities),
		a2asrv.WithAgentInactivityTimeout(timeout),
		a2asrv.WithCallInterceptors(&authInterceptor{
			verifier: config.Verifier,
			attrs:    config.AuthenticatedAttrs,
		}),
	)
	return &Server{
		cardHandler: a2asrv.NewStaticAgentCardHandler(config.Card),
		rpcHandler:  a2asrv.NewJSONRPCHandler(requestHandler),
		rpcPath:     rpcPath,
		maxBytes:    maxBytes,
	}, nil
}

// Routes returns an HTTP handler for the public card and authenticated JSON-RPC endpoint.
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.Handle(a2asrv.WellKnownAgentCardPath, s.cardHandler)
	mux.Handle(s.rpcPath, s.limitBody(s.rpcHandler))
	return mux
}

// WellKnownAgentCardPath returns the canonical A2A discovery path.
func WellKnownAgentCardPath() string {
	return a2asrv.WellKnownAgentCardPath
}

func (s *Server) limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		request.Body = http.MaxBytesReader(response, request.Body, s.maxBytes)
		next.ServeHTTP(response, request)
	})
}

type authInterceptor struct {
	a2asrv.PassthroughCallInterceptor
	verifier TokenVerifier
	attrs    map[string]any
}

func (a *authInterceptor) Before(
	ctx context.Context,
	callCtx *a2asrv.CallContext,
	_ *a2asrv.Request,
) (context.Context, any, error) {
	values, ok := callCtx.ServiceParams().Get("authorization")
	if !ok || len(values) != 1 {
		return ctx, nil, a2a.ErrUnauthenticated
	}
	token, ok := strings.CutPrefix(values[0], "Bearer ")
	if !ok || strings.TrimSpace(token) == "" {
		return ctx, nil, a2a.ErrUnauthenticated
	}
	identity, err := a.verifier.VerifyBearer(ctx, strings.TrimSpace(token))
	if err != nil || strings.TrimSpace(identity.Subject) == "" {
		return ctx, nil, a2a.ErrUnauthenticated
	}
	attrs := make(map[string]any, len(a.attrs)+1)
	for key, value := range a.attrs {
		attrs[key] = value
	}
	if identity.Tenant != "" {
		attrs["tenant"] = identity.Tenant
	}
	callCtx.User = a2asrv.NewAuthenticatedUser(identity.Subject, attrs)
	return ctx, nil, nil
}

// GoogleIDTokenVerifier validates Google-issued OpenID Connect identity tokens.
type GoogleIDTokenVerifier struct {
	Audience    string
	TenantClaim string
}

// VerifyBearer validates one bearer token and returns only safe identity facts.
func (v GoogleIDTokenVerifier) VerifyBearer(ctx context.Context, token string) (Identity, error) {
	if strings.TrimSpace(v.Audience) == "" {
		return Identity{}, errors.New("google id token verifier audience is required")
	}
	payload, err := idtoken.Validate(ctx, token, v.Audience)
	if err != nil {
		return Identity{}, fmt.Errorf("validate google id token: %w", err)
	}
	if payload.Subject == "" {
		return Identity{}, errors.New("google id token subject is empty")
	}
	identity := Identity{Subject: payload.Subject}
	if claim := strings.TrimSpace(v.TenantClaim); claim != "" {
		if tenant, ok := payload.Claims[claim].(string); ok {
			identity.Tenant = tenant
		}
	}
	return identity, nil
}

// StaticBearerVerifier supports local deterministic tests and single-tenant dev deployments.
type StaticBearerVerifier struct {
	Token   string
	Subject string
	Tenant  string
}

// VerifyBearer validates a configured opaque token without logging or exposing it.
func (v StaticBearerVerifier) VerifyBearer(_ context.Context, token string) (Identity, error) {
	if v.Token == "" || subtle.ConstantTimeCompare([]byte(token), []byte(v.Token)) != 1 {
		return Identity{}, errors.New("static bearer token is invalid")
	}
	subject := strings.TrimSpace(v.Subject)
	if subject == "" {
		subject = "a2a-caller"
	}
	return Identity{Subject: subject, Tenant: strings.TrimSpace(v.Tenant)}, nil
}

// RejectingVerifier keeps A2A execution fail-closed when auth is not configured.
type RejectingVerifier struct{}

// VerifyBearer always rejects execution.
func (RejectingVerifier) VerifyBearer(context.Context, string) (Identity, error) {
	return Identity{}, errors.New("a2a token verifier is not configured")
}

func cleanJSONRPCPath(value string) (string, error) {
	cleaned := strings.TrimSpace(value)
	if cleaned == "" {
		cleaned = defaultJSONRPCPath
	}
	if !strings.HasPrefix(cleaned, "/") ||
		strings.Contains(cleaned, "..") ||
		strings.ContainsAny(cleaned, "?#") {
		return "", errors.New("a2a json-rpc path is invalid")
	}
	return cleaned, nil
}

func drainAndClose(body io.ReadCloser) {
	_, _ = io.Copy(io.Discard, io.LimitReader(body, 4<<10))
	_ = body.Close()
}
