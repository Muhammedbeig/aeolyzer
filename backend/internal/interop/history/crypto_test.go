package history

import (
	"bytes"
	"encoding/base64"
	"errors"
	"testing"
)

func TestCipherRoundTrip(t *testing.T) {
	t.Parallel()

	cipher, err := NewCipher(bytes.Repeat([]byte{7}, encryptionKeyBytes))
	if err != nil {
		t.Fatal(err)
	}
	aad := []byte("tenant|session|event")
	ciphertext, err := cipher.Encrypt([]byte("protected"), aad)
	if err != nil {
		t.Fatal(err)
	}
	plaintext, err := cipher.Decrypt(ciphertext, aad)
	if err != nil {
		t.Fatal(err)
	}
	if string(plaintext) != "protected" {
		t.Fatalf("Decrypt() = %q, want protected", plaintext)
	}
	if bytes.Contains(ciphertext, plaintext) {
		t.Fatal("ciphertext contains plaintext")
	}
}

func TestCipherRejectsWrongTenantContext(t *testing.T) {
	t.Parallel()

	cipher, err := NewCipher(bytes.Repeat([]byte{8}, encryptionKeyBytes))
	if err != nil {
		t.Fatal(err)
	}
	ciphertext, err := cipher.Encrypt([]byte("protected"), []byte("tenant-a"))
	if err != nil {
		t.Fatal(err)
	}
	_, err = cipher.Decrypt(ciphertext, []byte("tenant-b"))
	if !errors.Is(err, ErrDecrypt) {
		t.Fatalf("Decrypt() error = %v, want %v", err, ErrDecrypt)
	}
}

func TestParseKey(t *testing.T) {
	t.Parallel()

	encoded := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{9}, encryptionKeyBytes))
	key, err := ParseKey(encoded)
	if err != nil {
		t.Fatal(err)
	}
	if len(key) != encryptionKeyBytes {
		t.Fatalf("ParseKey() length = %d, want %d", len(key), encryptionKeyBytes)
	}
	if _, err := ParseKey("short"); err == nil {
		t.Fatal("ParseKey() accepted an invalid key")
	}
}
