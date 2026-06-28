package middleware

import (
	"aeolyzer/internal/intake/contracts"
)

func GuardOutboundResponse(rawResponse string, intent contracts.Intent) (string, error) {
	if ContainsProtectedMetadata(rawResponse) {
		redacted, _, err := RedactProtectedMetadata(rawResponse)
		if err != nil {
			return "", err
		}
		return redacted, nil
	}

	return rawResponse, nil
}
