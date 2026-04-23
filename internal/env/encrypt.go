package env

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
	"strings"
)

const encryptedPrefix = "enc:"

// Encrypt encrypts all values in the set using AES-GCM with the provided key.
// The key must be 16, 24, or 32 bytes long (AES-128, AES-192, AES-256).
// Encrypted values are stored as base64-encoded strings prefixed with "enc:".
func Encrypt(s *Set, key []byte) error {
	if s == nil {
		return errors.New("encrypt: set must not be nil")
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return errors.New("encrypt: key must be 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	for k, v := range s.Vars {
		if strings.HasPrefix(v, encryptedPrefix) {
			continue
		}
		nonce := make([]byte, gcm.NonceSize())
		if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
			return err
		}
		ciphertext := gcm.Seal(nonce, nonce, []byte(v), nil)
		s.Vars[k] = encryptedPrefix + base64.StdEncoding.EncodeToString(ciphertext)
	}
	return nil
}

// Decrypt decrypts all "enc:"-prefixed values in the set using AES-GCM.
func Decrypt(s *Set, key []byte) error {
	if s == nil {
		return errors.New("decrypt: set must not be nil")
	}
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return errors.New("decrypt: key must be 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}
	for k, v := range s.Vars {
		if !strings.HasPrefix(v, encryptedPrefix) {
			continue
		}
		data, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(v, encryptedPrefix))
		if err != nil {
			return err
		}
		if len(data) < gcm.NonceSize() {
			return errors.New("decrypt: ciphertext too short for key " + k)
		}
		nonce, ciphertext := data[:gcm.NonceSize()], data[gcm.NonceSize():]
		plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
		if err != nil {
			return err
		}
		s.Vars[k] = string(plaintext)
	}
	return nil
}
