package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// Encrypt encrypts plain text using AES-GCM with the provided key.
func Encrypt(plainText, key string) (string, error) {
	// Convert key to byte array (key length must be 16, 24, or 32 bytes for AES)
	keyBytes := []byte(key)
	if len(keyBytes) != 32 { // Example uses a 256-bit key (32 bytes)
		return "", errors.New("key must be 32 bytes long")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create GCM (Galois/Counter Mode) cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Generate a random nonce (must be unique for each encryption)
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %v", err)
	}

	// Encrypt the data
	cipherText := aesGCM.Seal(nonce, nonce, []byte(plainText), nil)

	// Return the result as a base64-encoded string
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// Decrypt decrypts base64-encoded ciphertext using AES-GCM with the provided key.
func Decrypt(cipherText, key string) (string, error) {
	// Convert key to byte array
	keyBytes := []byte(key)
	if len(keyBytes) != 32 {
		return "", errors.New("key must be 32 bytes long")
	}

	// Decode the base64-encoded ciphertext
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %v", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// Create GCM (Galois/Counter Mode) cipher
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %v", err)
	}

	// Extract the nonce (first part of the ciphertext)
	nonceSize := aesGCM.NonceSize()
	if len(cipherTextBytes) < nonceSize {
		return "", errors.New("invalid ciphertext")
	}
	nonce, cipherTextBytes := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]

	// Decrypt the data
	plainText, err := aesGCM.Open(nil, nonce, cipherTextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %v", err)
	}

	return string(plainText), nil
}
