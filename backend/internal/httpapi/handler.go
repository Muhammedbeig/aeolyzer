package httpapi

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"aeolyzer/layer_02_intake"
	"aeolyzer/layer_03_orchestration"
	"aeolyzer/layer_05_extensions"
	"aeolyzer/layer_06_runtime"
	"aeolyzer/layer_08_observability"
)

const maxRequestBytes = 64 << 10 // BOUNDS: Hard limit payload memory allocation to mitigate DoS vectors.

type Handler struct {
	intake        *intake.Service
	orchestrator  *orchestration.Service
	executor      *runtime.Executor
	events        *observability.Sink
	logger        *slog.Logger
	allowedOrigin string
	now           func() time.Time
}

func NewHandler(
	intakeService *intake.Service,
	orchestrator *orchestration.Service,
	executor *runtime.Executor,
	events *observability.Sink,
	logger *slog.Logger,
	allowedOrigin string,
) *Handler {
	return &Handler{
		intake:        intakeService,
		orchestrator:  orchestrator,
		executor:      executor,
		events:        events,
		logger:        loggerOrDefault(logger),
		// INVARIANT: Normalize origin path to strip trailing slashes, neutralizing basic spoofing attacks.
		allowedOrigin: strings.TrimRight(allowedOrigin, "/"),
		now:           time.Now,
	}
}

func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("POST /v1/onboarding/inspect", h.inspectSite)
	mux.HandleFunc("POST /v1/onboarding/complete", h.completeOnboarding)
	return h.withMiddleware(mux)
}

func (h *Handler) health(response http.ResponseWriter, _ *http.Request) {
	writeJSON(response, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) inspectSite(response http.ResponseWriter, request *http.Request) {
	var input intake.SiteInspectionInput
	if err := decodeJSON(response, request, &input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid_request", "Check the website URL and try again.")
		return
	}

	decision, err := h.intake.InspectSite(input)
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid_website", "Enter a valid public website URL.")
		return
	}
	plan, err := h.orchestrator.PlanSiteInspection(decision)
	if err != nil {
		h.record(decision.TraceID, "site_inspection", "denied")
		writeError(response, http.StatusForbidden, "inspection_denied", "This website cannot be inspected.")
		return
	}
	result, err := h.executor.Execute(request.Context(), plan.Request)
	if err != nil {
		h.record(decision.TraceID, "site_inspection", "failed")
		h.logger.WarnContext(request.Context(), "site inspection failed", "trace_id", decision.TraceID, "error", err)
		writeError(response, http.StatusUnprocessableEntity, "inspection_failed", "We could not read that website. You can continue by entering the brand details manually.")
		return
	}

	h.record(decision.TraceID, "site_inspection", "succeeded")
	writeJSON(response, http.StatusOK, map[string]any{
		"canonical_url":         decision.CanonicalURL,
		"suggested_brand_name":  result.Title,
		"description":           result.Description,
		"category":              result.Category,
		"icon_url":              result.IconURL,
		"competitor_candidates": result.CandidateCompetitors,
	})
}

func (h *Handler) completeOnboarding(response http.ResponseWriter, request *http.Request) {
	var input intake.OnboardingInput
	if err := decodeJSON(response, request, &input); err != nil {
		writeError(response, http.StatusBadRequest, "invalid_request", "Check the project details and try again.")
		return
	}

	decision, err := h.intake.CompleteOnboarding(input)
	if err != nil {
		writeError(response, http.StatusBadRequest, "invalid_profile", "Complete each required project field before continuing.")
		return
	}
	category := strings.TrimSpace(request.URL.Query().Get("category"))
	promptPlan, err := h.orchestrator.BuildPromptPlan(decision, category)
	if err != nil {
		h.record(decision.TraceID, "onboarding", "failed")
		writeError(response, http.StatusUnprocessableEntity, "prompt_plan_failed", "We could not prepare the project prompts.")
		return
	}
	frame, err := extensions.BuildDashboardFrame(extensions.DashboardIntent{
		TraceID:     decision.TraceID,
		GeneratedAt: h.now(),
		Profile:     decision.Profile,
		Prompts:     promptPlan.Prompts,
	})
	if err != nil {
		h.record(decision.TraceID, "onboarding", "failed")
		writeError(response, http.StatusInternalServerError, "presentation_failed", "We could not prepare the dashboard.")
		return
	}

	h.record(decision.TraceID, "onboarding", "succeeded")
	writeJSON(response, http.StatusOK, frame)
}

func (h *Handler) record(traceID, eventType, outcome string) {
	if h.events == nil {
		return
	}
	h.events.Record(observability.Event{
		TraceID:   traceID,
		EventType: eventType,
		Outcome:   outcome,
		At:        h.now().UTC(),
	})
}

func (h *Handler) withMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		// PERIMETER: Force explicit content interpretation and disable cache poisoning vectors.
		response.Header().Set("Content-Type", "application/json; charset=utf-8")
		response.Header().Set("X-Content-Type-Options", "nosniff")
		response.Header().Set("Cache-Control", "no-store")

		origin := strings.TrimRight(request.Header.Get("Origin"), "/")
		if origin != "" && origin != h.allowedOrigin {
			writeError(response, http.StatusForbidden, "origin_denied", "This origin is not allowed.")
			return
		}
		if origin != "" && origin == h.allowedOrigin {
			response.Header().Set("Access-Control-Allow-Origin", origin)
			response.Header().Set("Vary", "Origin")
			response.Header().Set("Access-Control-Allow-Headers", "Content-Type")
			response.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		}
		if request.Method == http.MethodOptions {
			if origin == "" || origin != h.allowedOrigin {
				writeError(response, http.StatusForbidden, "origin_denied", "This origin is not allowed.")
				return
			}
			response.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(response, request)
	})
}

func decodeJSON(response http.ResponseWriter, request *http.Request, target any) error {
	request.Body = http.MaxBytesReader(response, request.Body, maxRequestBytes)
	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields() // INVARIANT: Strict schema enforcement to prevent unmarshaling silent injections.
	if err := decoder.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		// EDGE CASE: Trap trailing payload junk to invalidate concatenated attacks.
		return err
	}
	return nil
}

func loggerOrDefault(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	return slog.Default()
}

func writeError(response http.ResponseWriter, status int, code, message string) {
	writeJSONStatus(response, status, map[string]any{
		"error": map[string]string{
			"code":    code,
			"message": message,
		},
	})
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	writeJSONStatus(response, status, value)
}

func writeJSONStatus(response http.ResponseWriter, status int, value any) {
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}
