package history

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/adk/model"
	"google.golang.org/adk/session"
	"google.golang.org/genai"
)

func TestMariaDBSessionLifecycle(t *testing.T) {
	dsn := os.Getenv("AEOLYZER_TEST_DB_DSN")
	if dsn == "" {
		t.Skip("AEOLYZER_TEST_DB_DSN is not set")
	}
	ctx, cancel := context.WithTimeout(t.Context(), 30*time.Second)
	defer cancel()

	store, err := Open(ctx, dsn, bytes.Repeat([]byte{3}, 32), DefaultConfig())
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	suffix := strings.ReplaceAll(uuid.NewString()[:8], "-", "")
	appName := "testapp-" + suffix
	userID := "testuser-" + suffix
	created, err := store.CreateConversation(ctx, appName, userID)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Delete(context.Background(), &session.DeleteRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: created.ID,
	})

	getResponse, err := store.Get(ctx, &session.GetRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	fileData := []byte("attachment evidence")
	digest := sha256.Sum256(fileData)
	ref, err := store.SaveAttachment(ctx, appName, userID, created.ID, AttachmentInput{
		Name:        "evidence.txt",
		ContentType: "text/plain",
		Data:        fileData,
		SHA256:      digest,
	})
	if err != nil {
		t.Fatal(err)
	}
	userEvent := session.NewEvent("invocation-1")
	userEvent.Author = "user"
	userEvent.LLMResponse = model.LLMResponse{
		Content: &genai.Content{
			Role: genai.RoleUser,
			Parts: []*genai.Part{
				genai.NewPartFromText("Review this evidence"),
				genai.NewPartFromBytes(fileData, "text/plain"),
			},
		},
	}
	if err := store.AppendEvent(WithAttachmentRefs(ctx, []AttachmentRef{ref}), getResponse.Session, userEvent); err != nil {
		t.Fatal(err)
	}
	modelEvent := session.NewEvent("invocation-1")
	modelEvent.Author = "audit_agent"
	modelEvent.LLMResponse = model.LLMResponse{
		Content: genai.NewContentFromText("The evidence is readable.", genai.RoleModel),
	}
	if err := store.AppendEvent(ctx, getResponse.Session, modelEvent); err != nil {
		t.Fatal(err)
	}

	resumed, err := store.Get(ctx, &session.GetRequest{
		AppName:   appName,
		UserID:    userID,
		SessionID: created.ID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resumed.Session.Events().Len() != 2 {
		t.Fatalf("Get() event count = %d, want 2", resumed.Session.Events().Len())
	}
	hydrated := resumed.Session.Events().At(0).Content.Parts[1].InlineData
	if hydrated == nil || !bytes.Equal(hydrated.Data, fileData) {
		t.Fatal("Get() did not hydrate the attachment for ADK context")
	}
	if resumed.Session.Events().At(0).Content.Parts[1].PartMetadata != nil {
		t.Fatal("Get() exposed internal attachment metadata")
	}

	messages, err := store.ListMessages(ctx, appName, userID, created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(messages) != 2 || len(messages[0].Attachments) != 1 {
		t.Fatalf("ListMessages() returned %#v", messages)
	}

	var storedCiphertext []byte
	if err := store.db.QueryRowContext(
		ctx,
		`SELECT data_ciphertext FROM aeolyzer_attachments
		  WHERE app_name = ? AND user_id = ? AND session_id = ? AND attachment_id = ?`,
		appName,
		userID,
		created.ID,
		ref.ID,
	).Scan(&storedCiphertext); err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(storedCiphertext, fileData) {
		t.Fatal("attachment plaintext was present in database storage")
	}

	claim, err := store.ClaimMessageRequest(ctx, appName, userID, created.ID, "request-1234567890")
	if err != nil || !claim.Claimed {
		t.Fatalf("ClaimMessageRequest() = %#v, %v", claim, err)
	}
	if err := store.CompleteMessageRequest(
		ctx,
		appName,
		userID,
		created.ID,
		"request-1234567890",
		[]byte(`{"ok":true}`),
	); err != nil {
		t.Fatal(err)
	}
	cached, err := store.ClaimMessageRequest(ctx, appName, userID, created.ID, "request-1234567890")
	if err != nil || string(cached.CachedResponse) != `{"ok":true}` {
		t.Fatalf("ClaimMessageRequest() cached = %q, %v", cached.CachedResponse, err)
	}

	_, err = store.GetConversation(ctx, appName, "different-user", created.ID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("cross-tenant GetConversation() error = %v, want %v", err, ErrNotFound)
	}
}
