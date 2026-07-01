package orchestrator

import (
	"context"
	"strings"
	"testing"

	"aeolyzer/internal/extensions"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func TestBuildAgentContextInstruction(t *testing.T) {
	t.Parallel()

	instruction := buildAgentContextInstruction(agentModelContext{
		knowledge:   "Brand profile\n- Name: AEOlyzer",
		contentType: extensions.ContentTypeBlogPost,
	})
	for _, want := range []string{
		"Active content type: blog_post.",
		"User-approved tenant context follows.",
		"Name: AEOlyzer",
	} {
		if !strings.Contains(instruction, want) {
			t.Fatalf("buildAgentContextInstruction() = %q, want %q", instruction, want)
		}
	}
}

func TestInjectAgentModelContext(t *testing.T) {
	t.Parallel()

	userContent := genai.NewContentFromText("Draft the article", genai.RoleUser)
	callbackContext := modelCallbackContext{
		Context: withAgentModelContext(
			context.Background(),
			"Brand profile\n- Name: AEOlyzer",
			extensions.ContentTypeArticle,
		),
	}
	request := &model.LLMRequest{Contents: []*genai.Content{userContent}}

	if _, err := InjectAgentModelContext(callbackContext, request); err != nil {
		t.Fatal(err)
	}
	if request.Config == nil || request.Config.SystemInstruction == nil {
		t.Fatal("InjectAgentModelContext() did not add a system instruction")
	}
	instruction := request.Config.SystemInstruction.Parts[0].Text
	if !strings.Contains(instruction, "Name: AEOlyzer") {
		t.Fatalf("system instruction = %q, want tenant context", instruction)
	}
	if len(request.Contents) != 1 || request.Contents[0] != userContent {
		t.Fatal("InjectAgentModelContext() changed persisted user contents")
	}
}

type modelCallbackContext struct {
	context.Context
}

func (modelCallbackContext) UserContent() *genai.Content          { return nil }
func (modelCallbackContext) InvocationID() string                 { return "" }
func (modelCallbackContext) AgentName() string                    { return "" }
func (modelCallbackContext) ReadonlyState() session.ReadonlyState { return nil }
func (modelCallbackContext) UserID() string                       { return "" }
func (modelCallbackContext) AppName() string                      { return "" }
func (modelCallbackContext) SessionID() string                    { return "" }
func (modelCallbackContext) Branch() string                       { return "" }
func (modelCallbackContext) Artifacts() agent.Artifacts           { return nil }
func (modelCallbackContext) State() session.State                 { return nil }
