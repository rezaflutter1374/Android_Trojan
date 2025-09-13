package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

func EncryptAES(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, readErr := io.ReadFull(rand.Reader, nonce); readErr != nil {

		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {

		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	encrypted := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func DecryptAES(encryptedStr string, key []byte) ([]byte, error) {

	encrypted, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return nil, err
	}

	if len(encrypted) < 13 {

		return nil, errors.New("invalid ciphertext length")

	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := encrypted[:12]
	ciphertext := encrypted[12:]

	return aesgcm.Open(nil, nonce, ciphertext, nil)
}
