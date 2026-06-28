// Package vectorragstore enforces tenant-partitioned vector retrieval.
package vectorragstore

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
)

// NamespaceForTenant derives a non-enumerable tenant namespace.
func NamespaceForTenant(key []byte, tenantID string) (string, error) {
	if len(key) < 32 || tenantID == "" {
		return "", errors.New("tenant namespace key and tenant id are required")
	}
	mac := hmac.New(sha256.New, key)
	_, _ = mac.Write([]byte("aeolyzer-vector-namespace-v1\x00"))
	_, _ = mac.Write([]byte(tenantID))
	return "tenant_" + hex.EncodeToString(mac.Sum(nil)), nil
}

// VerifyTenantNamespace checks an exact tenant-to-namespace binding.
func VerifyTenantNamespace(key []byte, tenantID, namespace string) error {
	expected, err := NamespaceForTenant(key, tenantID)
	if err != nil {
		return err
	}
	if len(expected) != len(namespace) ||
		!hmac.Equal([]byte(expected), []byte(namespace)) {
		return errors.New("vector namespace tenant binding mismatch")
	}
	return nil
}
