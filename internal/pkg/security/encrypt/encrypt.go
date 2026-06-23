package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

func EncryptAES(dataBytes []byte, keyBytes []byte) (string, error) {
	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return "", err
	}

	cipherText := aesGCM.Seal(nonce, nonce, dataBytes, nil)

	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(cipherBase64 string, keyBytes []byte) ([]byte, error) {
	cipherText, err := base64.URLEncoding.DecodeString(cipherBase64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(keyBytes)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce, cipherText := cipherText[:nonceSize], cipherText[nonceSize:]

	return aesGCM.Open(nil, nonce, cipherText, nil)
}
