package api

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecryptToken(t *testing.T) {
	// Generate a valid AES-256 key (32 bytes)
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	// Helper function to encrypt a token for testing
	encryptTokenForTest := func(key []byte, token string) (string, error) {
		block, err := aes.NewCipher(key)
		if err != nil {
			return "", err
		}

		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return "", err
		}

		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return "", err
		}

		ciphertext := gcm.Seal(nonce, nonce, []byte(token), nil)
		return base64.StdEncoding.EncodeToString(ciphertext), nil
	}

	tests := []struct {
		name          string
		encryptedToken string
		expectedToken string
		expectedError error
		setupKey      func() []byte
	}{
		{
			name:          "success - decrypt valid token",
			encryptedToken: "", // Will be set in test
			expectedToken: "test-token-123",
			setupKey:      func() []byte { return validKey },
		},
		{
			name:          "success - decrypt empty token",
			encryptedToken: "", // Will be set in test
			expectedToken: "",
			setupKey:      func() []byte { return validKey },
		},
		{
			name:          "success - decrypt token with special characters",
			encryptedToken: "", // Will be set in test
			expectedToken: "token-with-!@#$%^&*()",
			setupKey:      func() []byte { return validKey },
		},
		{
			name:          "error - invalid base64 encoding",
			encryptedToken: "invalid-base64!@#",
			expectedError: errors.New("illegal base64 data"),
			setupKey:      func() []byte { return validKey },
		},
		{
			name:          "error - ciphertext too short",
			encryptedToken: base64.StdEncoding.EncodeToString([]byte("short")),
			expectedError: errors.New("ciphertext too short"),
			setupKey:      func() []byte { return validKey },
		},
		{
			name:          "error - invalid key size",
			encryptedToken: "", // Will be set in test - encrypted with valid key
			expectedError: errors.New("crypto/aes: invalid key size"),
			setupKey:      func() []byte { return []byte("too-short") }, // Invalid key size (only 9 bytes, need 16/24/32)
		},
		{
			name:          "error - wrong key",
			encryptedToken: "", // Will be set in test
			expectedError: errors.New("cipher: message authentication failed"),
			setupKey:      func() []byte {
				wrongKey := make([]byte, 32)
				for i := range wrongKey {
					wrongKey[i] = byte(i + 100) // Different key
				}
				return wrongKey
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := tt.setupKey()
			api := NewApi(nil, key)

			// Encrypt token if needed for success cases
			var encryptedToken string
			if tt.encryptedToken == "" && tt.expectedError == nil {
				var err error
				encryptedToken, err = encryptTokenForTest(key, tt.expectedToken)
				if err != nil {
					t.Fatalf("Failed to encrypt token for test: %v", err)
				}
			} else if tt.encryptedToken == "" && tt.expectedError != nil {
				// For error cases that need encrypted token, encrypt with correct key first
				if tt.name == "error - wrong key" || tt.name == "error - invalid key size" {
					encryptedToken, _ = encryptTokenForTest(validKey, "test-token")
				} else {
					encryptedToken = tt.encryptedToken
				}
			} else {
				encryptedToken = tt.encryptedToken
			}

			decryptedToken, err := api.DecryptToken(encryptedToken)

			if tt.expectedError != nil {
				assert.Error(t, err)
				// Check if error message contains expected substring
				if tt.expectedError.Error() == "cipher: message authentication failed" {
					assert.Contains(t, err.Error(), "message authentication failed")
				} else if tt.expectedError.Error() == "crypto/aes: invalid key size" {
					assert.Contains(t, err.Error(), "invalid key size")
				} else {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				}
				assert.Empty(t, decryptedToken)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedToken, decryptedToken)
			}
		})
	}
}

// Test encryption/decryption round-trip
func TestEncryptDecryptRoundTrip(t *testing.T) {
	validKey := make([]byte, 32)
	for i := range validKey {
		validKey[i] = byte(i)
	}

	api := NewApi(nil, validKey)

	testTokens := []string{
		"simple-token",
		"token-with-special-chars-!@#$%^&*()",
		"token-with-newline\nand-tab\t",
		"very-long-token-" + string(make([]byte, 1000)),
		"",
	}

	for _, originalToken := range testTokens {
		t.Run("round-trip-"+originalToken[:min(len(originalToken), 20)], func(t *testing.T) {
			// Encrypt
			encrypted, err := api.encryptToken(originalToken)
			assert.NoError(t, err)
			assert.NotEmpty(t, encrypted)
			assert.NotEqual(t, originalToken, encrypted)

			// Decrypt
			decrypted, err := api.DecryptToken(encrypted)
			assert.NoError(t, err)
			assert.Equal(t, originalToken, decrypted)
		})
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
