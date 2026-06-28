package llmprovider

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	evaluationengine "aeolyzer/internal/observability/evaluation_engine"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}

func TestNewGeminiClientFromEnvironmentRequiresCredential(t *testing.T) {
	t.Setenv(geminiAPIKeyEnvironment, "")
	if _, err := NewGeminiJudgeClientFromEnvironment(nil); err == nil {
		t.Fatal("NewGeminiJudgeClientFromEnvironment() error = nil, want error")
	}
}

func TestGeminiClientJudgeUsesHeaderAndStructuredOutput(t *testing.T) {
	const apiKey = "test-only-secret"
	t.Setenv(geminiAPIKeyEnvironment, apiKey)

	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		if got := request.URL.String(); got != geminiInteractionsURL {
			t.Fatalf("request URL = %q, want %q", got, geminiInteractionsURL)
		}
		if strings.Contains(request.URL.String(), apiKey) {
			t.Fatal("request URL contains credential")
		}
		if got := request.Header.Get("x-goog-api-key"); got != apiKey {
			t.Fatalf("api key header = %q, want test credential", got)
		}
		body, err := io.ReadAll(request.Body)
		if err != nil {
			t.Fatalf("read request body: %v", err)
		}
		if bytes.Contains(body, []byte(apiKey)) {
			t.Fatal("request body contains credential")
		}
		var payload geminiInteractionRequest
		if err := json.Unmarshal(body, &payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if payload.Store || payload.Background {
			t.Fatal("Gemini request enables storage or background execution")
		}
		if payload.GenerationConfig.Temperature != 0 ||
			payload.GenerationConfig.ThinkingSummaries != "none" {
			t.Fatal("Gemini request is not deterministic and reasoning-minimized")
		}
		if payload.ResponseFormat.MIMEType != "application/json" ||
			payload.ResponseFormat.Schema["additionalProperties"] != false {
			t.Fatal("Gemini request does not enforce a closed JSON schema")
		}
		return geminiResponse(
			http.StatusOK,
			`{"status":"completed","steps":[{"type":"model_output","content":[{"type":"text","text":"{\"rubric_id\":\"quality-v1\",\"dimensions\":{\"grounding\":5},\"confidence\":0.98,\"summary\":\"Grounded in supplied evidence.\"}"}]}]}`,
		), nil
	})

	client, err := NewGeminiJudgeClientFromEnvironment(transport)
	if err != nil {
		t.Fatalf("NewGeminiJudgeClientFromEnvironment() error = %v", err)
	}
	raw, err := client.Judge(context.Background(), testGeminiJudgeRequest(false))
	if err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
	if _, err := evaluationengine.ValidateJudgeOutput(raw); err != nil {
		t.Fatalf("ValidateJudgeOutput() error = %v", err)
	}
	if strings.Contains(fmt.Sprintf("%+v", client), apiKey) ||
		strings.Contains(fmt.Sprintf("%#v", client), apiKey) {
		t.Fatal("formatted client contains credential")
	}
}

func TestGeminiClientJudgePairwiseSchemaRequiresWinner(t *testing.T) {
	t.Setenv(geminiAPIKeyEnvironment, "test-only-secret")
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		var payload geminiInteractionRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		required, ok := payload.ResponseFormat.Schema["required"].([]any)
		if !ok {
			t.Fatal("required schema field has unexpected type")
		}
		foundWinner := false
		for _, field := range required {
			foundWinner = foundWinner || field == "winner"
		}
		if !foundWinner {
			t.Fatal("pairwise schema does not require winner")
		}
		return geminiResponse(
			http.StatusOK,
			`{"status":"completed","steps":[{"type":"model_output","content":[{"type":"text","text":"{\"rubric_id\":\"quality-v1\",\"dimensions\":{\"grounding\":5},\"winner\":\"A\",\"confidence\":0.98,\"summary\":\"A is better grounded.\"}"}]}]}`,
		), nil
	})
	client, err := NewGeminiJudgeClientFromEnvironment(transport)
	if err != nil {
		t.Fatalf("NewGeminiJudgeClientFromEnvironment() error = %v", err)
	}
	if _, err := client.Judge(context.Background(), testGeminiJudgeRequest(true)); err != nil {
		t.Fatalf("Judge() error = %v", err)
	}
}

func TestGeminiClientSelectSkillsUsesOpaqueClosedSchema(t *testing.T) {
	t.Setenv(geminiAPIKeyEnvironment, "test-only-secret")
	transport := roundTripFunc(func(request *http.Request) (*http.Response, error) {
		var payload geminiInteractionRequest
		if err := json.NewDecoder(request.Body).Decode(&payload); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		if strings.Contains(payload.Input, "internal_skill_id") {
			t.Fatal("skill selection request disclosed an internal skill id")
		}
		if !strings.Contains(payload.Input, "candidate_1") {
			t.Fatal("skill selection request omitted opaque candidate alias")
		}
		if payload.ResponseFormat.Schema["additionalProperties"] != false {
			t.Fatal("skill selection schema permits unknown top-level fields")
		}
		return geminiResponse(
			http.StatusOK,
			`{"status":"completed","steps":[{"type":"model_output","content":[{"type":"text","text":"{\"results\":[{\"case_id\":\"case_1\",\"selected_candidate\":\"candidate_1\",\"confidence\":0.95,\"summary\":\"Candidate one matches.\"}]}"}]}]}`,
		), nil
	})
	client, err := NewGeminiJudgeClientFromEnvironment(transport)
	if err != nil {
		t.Fatalf("NewGeminiJudgeClientFromEnvironment() error = %v", err)
	}
	raw, err := client.SelectSkills(
		context.Background(),
		evaluationengine.SkillSelectionRequest{
			Model:         "gemini-3.5-flash",
			PromptVersion: "skill-router-v1",
			EvalID:        "eval-1",
			Candidates: []evaluationengine.SkillSelectionCandidate{
				{
					Alias:        "candidate_1",
					Description:  "First bounded skill description.",
					AntiTriggers: []string{"unrelated first request"},
				},
				{
					Alias:        "candidate_2",
					Description:  "Second bounded skill description.",
					AntiTriggers: []string{"unrelated second request"},
				},
				{
					Alias:        "candidate_3",
					Description:  "Third bounded skill description.",
					AntiTriggers: []string{"unrelated third request"},
				},
				{
					Alias:        "candidate_4",
					Description:  "Fourth bounded skill description.",
					AntiTriggers: []string{"unrelated fourth request"},
				},
				{
					Alias:        "candidate_5",
					Description:  "Fifth bounded skill description.",
					AntiTriggers: []string{"unrelated fifth request"},
				},
			},
			Cases: []evaluationengine.SkillSelectionCase{{
				CaseID: "case_1",
				Input:  "Use the first bounded skill.",
			}},
		},
	)
	if err != nil {
		t.Fatalf("SelectSkills() error = %v", err)
	}
	if !bytes.Contains(raw, []byte(`"selected_candidate":"candidate_1"`)) {
		t.Fatalf("SelectSkills() output = %s", raw)
	}
}

func TestGeminiClientJudgeRejectsProviderFailureWithoutLeakingBody(t *testing.T) {
	const secret = "test-only-secret"
	t.Setenv(geminiAPIKeyEnvironment, secret)
	transport := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return geminiResponse(
			http.StatusUnauthorized,
			`{"error":{"code":401,"status":"UNAUTHENTICATED","message":"do not echo `+secret+`"}}`,
		), nil
	})
	client, err := NewGeminiJudgeClientFromEnvironment(transport)
	if err != nil {
		t.Fatalf("NewGeminiJudgeClientFromEnvironment() error = %v", err)
	}
	_, err = client.Judge(context.Background(), testGeminiJudgeRequest(false))
	if err == nil {
		t.Fatal("Judge() error = nil, want error")
	}
	if strings.Contains(err.Error(), secret) || strings.Contains(err.Error(), "do not echo") {
		t.Fatalf("Judge() error leaked provider body: %v", err)
	}
}

func TestDecodeGeminiResponseRejectsUnsafeShapes(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{name: "invalid json", body: `{`},
		{name: "not completed", body: `{"status":"requires_action"}`},
		{name: "no output", body: `{"status":"completed","steps":[]}`},
		{
			name: "multiple outputs",
			body: `{"status":"completed","steps":[` +
				`{"type":"model_output","content":[{"type":"text","text":"{}"}]},` +
				`{"type":"model_output","content":[{"type":"text","text":"{}"}]}` +
				`]}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if _, err := decodeGeminiResponse(strings.NewReader(test.body)); err == nil {
				t.Fatal("decodeGeminiResponse() error = nil, want error")
			}
		})
	}
}

func TestGeminiClientJudgeHonorsContext(t *testing.T) {
	t.Setenv(geminiAPIKeyEnvironment, "test-only-secret")
	transport := roundTripFunc(func(_ *http.Request) (*http.Response, error) {
		return nil, errors.New("transport should not be called")
	})
	client, err := NewGeminiJudgeClientFromEnvironment(transport)
	if err != nil {
		t.Fatalf("NewGeminiJudgeClientFromEnvironment() error = %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := client.Judge(ctx, testGeminiJudgeRequest(false)); !errors.Is(err, context.Canceled) {
		t.Fatalf("Judge() error = %v, want context.Canceled", err)
	}
}

func geminiResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func testGeminiJudgeRequest(pairwise bool) evaluationengine.JudgeRequest {
	input := evaluationengine.JudgeInput{
		EvalID:      "eval-1",
		Instruction: "Assess grounding.",
		CandidateA:  "Candidate A",
		Pairwise:    pairwise,
		Redacted:    true,
	}
	if pairwise {
		input.CandidateB = "Candidate B"
	}
	return evaluationengine.JudgeRequest{
		Model:         "gemini-3.5-flash",
		PromptVersion: "judge-v1",
		Temperature:   0,
		Input:         input,
		Rubric: evaluationengine.Rubric{
			ID:      "quality-v1",
			Version: "1.0.0",
			Dimensions: []evaluationengine.RubricDimension{{
				ID:          "grounding",
				Description: "Uses supplied evidence.",
				Weight:      1,
				Minimum:     1,
				Maximum:     5,
			}},
			PassScore: 4,
		},
	}
}
