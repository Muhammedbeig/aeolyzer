package middleware

import (
	"aeolyzer/layer_02_intake/contracts"
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
