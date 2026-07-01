package middleware

import (
	"errors"
	"unicode/utf8"

	"aeolyzer/internal/intake/contracts"

	"google.golang.org/adk/agent"
	"google.golang.org/adk/model"
)

const maxDisplayResponseRunes = 50_000

func GuardModelResponse(
	_ agent.CallbackContext,
	response *model.LLMResponse,
	responseError error,
) (*model.LLMResponse, error) {
	if responseError != nil {
		return nil, responseError
	}
	if response == nil || response.Content == nil {
		return nil, errors.New("model response is empty")
	}
	parts := response.Content.Parts[:0]
	for _, part := range response.Content.Parts {
		if part == nil || part.Thought {
			continue
		}
		part.ThoughtSignature = nil
		if part.Text != "" {
			guarded, err := GuardOutboundResponse(part.Text, contracts.IntentFallbackClarification)
			if err != nil {
				return nil, err
			}
			if utf8.RuneCountInString(guarded) > maxDisplayResponseRunes {
				runes := []rune(guarded)
				guarded = string(runes[:maxDisplayResponseRunes])
			}
			part.Text = guarded
		}
		parts = append(parts, part)
	}
	response.Content.Parts = parts
	return response, nil
}
