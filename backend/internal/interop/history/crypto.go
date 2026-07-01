package history

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
)

const (
	encryptionKeyBytes = 32
	cipherVersion      = byte(1)
)

var ErrDecrypt = errors.New("decrypt protected data")

type Cipher struct {
	aead cipher.AEAD
}

func NewCipher(key []byte) (*Cipher, error) {
	if len(key) != encryptionKeyBytes {
		return nil, fmt.Errorf("encryption key must be %d bytes", encryptionKeyBytes)
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create aes cipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("create gcm cipher: %w", err)
	}
	return &Cipher{aead: aead}, nil
}

func ParseKey(value string) ([]byte, error) {
	decoders := []func(string) ([]byte, error){
		base64.StdEncoding.DecodeString,
		base64.RawStdEncoding.DecodeString,
		base64.RawURLEncoding.DecodeString,
		hex.DecodeString,
	}
	for _, decode := range decoders {
		key, err := decode(value)
		if err == nil && len(key) == encryptionKeyBytes {
			return key, nil
		}
	}
	return nil, fmt.Errorf("data encryption key must encode exactly %d bytes", encryptionKeyBytes)
}

func (c *Cipher) Encrypt(plaintext, additionalData []byte) ([]byte, error) {
	if c == nil || c.aead == nil {
		return nil, errors.New("cipher is not configured")
	}
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, fmt.Errorf("generate encryption nonce: %w", err)
	}
	output := make([]byte, 1, 1+len(nonce)+len(plaintext)+c.aead.Overhead())
	output[0] = cipherVersion
	output = append(output, nonce...)
	output = c.aead.Seal(output, nonce, plaintext, additionalData)
	return output, nil
}

func (c *Cipher) Decrypt(ciphertext, additionalData []byte) ([]byte, error) {
	if c == nil || c.aead == nil || len(ciphertext) < 1+c.aead.NonceSize()+c.aead.Overhead() {
		return nil, ErrDecrypt
	}
	if ciphertext[0] != cipherVersion {
		return nil, ErrDecrypt
	}
	nonceEnd := 1 + c.aead.NonceSize()
	plaintext, err := c.aead.Open(nil, ciphertext[1:nonceEnd], ciphertext[nonceEnd:], additionalData)
	if err != nil {
		return nil, ErrDecrypt
	}
	return plaintext, nil
}
