package mcptransportplane

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"encoding/json"
	"errors"
)

// SchemaHash returns the hash of a compact canonical JSON representation.
func SchemaHash(schema []byte) (string, error) {
	if len(schema) == 0 || len(schema) > 2<<20 {
		return "", errors.New("mcp schema size is invalid")
	}
	var document any
	if err := json.Unmarshal(schema, &document); err != nil {
		return "", errors.New("mcp schema is invalid json")
	}
	canonical, err := json.Marshal(document)
	if err != nil {
		return "", errors.New("mcp schema cannot be canonicalized")
	}
	digest := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(digest[:]), nil
}

// VerifySchemaHash requires an exact pinned schema digest.
func VerifySchemaHash(schema []byte, expected string) error {
	actual, err := SchemaHash(schema)
	if err != nil {
		return err
	}
	if len(actual) != len(expected) ||
		subtle.ConstantTimeCompare([]byte(actual), []byte(expected)) != 1 {
		return errors.New("mcp schema hash mismatch")
	}
	return nil
}
