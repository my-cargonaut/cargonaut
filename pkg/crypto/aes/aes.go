package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// Encrypt uses AES to securly encrypt a plaintext using the provided key.
// The key must be 16, 24 or 32 byte broad to select either AES-128, AES-192 or
// AES-256.
func Encrypt(key, plaintext []byte) (string, error) {
	// Create the AES blockcipher and the GCM block cipher mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Pad plaintext, reserve memory for the ciphertext and create the nonce,
	// which gets stored before the ciphertext.
	nonce := make([]byte, gcm.NonceSize(), gcm.NonceSize()+len(plaintext))
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Finally encrypt the plaintext with the block cipher mode. The nonce is
	// used as destination because seal will append the ciphertext to it. The
	// nonce is used for decrypting the ciphertext so it must be stored with the
	// actual ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts an AES encrypted message using the provided key. The key
// must be 16, 24 or 32 byte broad to select either AES-128, AES-192 or AES-256.
func Decrypt(key []byte, ciphertext string) (string, error) {
	// Create the AES blockcipher and the GCM block cipher mode.
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Add base64 padding to the message and decode it.
	cipherStr, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}

	// Verify the ciphertext integrity. The ciphertext must be greater than the
	// nonce size, because the nonce prepends the actual ciphertext.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	// separate nonce and actual ciphertext.
	nonce, cipherStr := cipherStr[:nonceSize], cipherStr[nonceSize:]

	// Decrypt the actual ciphertext with the block cipher mode and the nonce.
	plaintext, err := gcm.Open(nil, nonce, cipherStr, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}
