package intake

import (
	"errors"
	"testing"

	"aeolyzer/internal/workspace"
)

func TestNormalizeChatContentType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		agent       string
		contentType workspace.ContentType
		want        workspace.ContentType
		wantErr     error
	}{
		{
			name:  "content default",
			agent: "content",
			want:  workspace.ContentTypeArticle,
		},
		{
			name:        "content selected",
			agent:       "content",
			contentType: workspace.ContentTypeLinkedInPost,
			want:        workspace.ContentTypeLinkedInPost,
		},
		{
			name:        "audit rejects content type",
			agent:       "audit",
			contentType: workspace.ContentTypeArticle,
			wantErr:     ErrInvalidContentType,
		},
		{
			name:        "unknown rejected",
			agent:       "content",
			contentType: "unknown",
			wantErr:     ErrInvalidContentType,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := NormalizeChatContentType(test.agent, test.contentType)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("NormalizeChatContentType() error = %v, want %v", err, test.wantErr)
			}
			if got != test.want {
				t.Fatalf("NormalizeChatContentType() = %q, want %q", got, test.want)
			}
		})
	}
}
