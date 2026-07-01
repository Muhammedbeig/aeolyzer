package intake

import (
	"errors"

	"aeolyzer/internal/workspace"
)

var ErrInvalidContentType = errors.New("invalid content type")

func NormalizeChatContentType(
	chatAgent string,
	contentType workspace.ContentType,
) (workspace.ContentType, error) {
	switch chatAgent {
	case "audit":
		if contentType != "" {
			return "", ErrInvalidContentType
		}
		return "", nil
	case "content":
		if contentType == "" {
			return workspace.ContentTypeArticle, nil
		}
		if !contentType.Valid() {
			return "", ErrInvalidContentType
		}
		return contentType, nil
	default:
		return "", ErrInvalidContentType
	}
}
