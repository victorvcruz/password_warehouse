package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"os"
)

var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

type ServiceClient interface {
	Encrypt(text string) (string, error)
	Decrypt(text string) (string, error)
}

type service struct {
	secret    string
	secretApi string
}

func NewCrypto() ServiceClient {
	return &service{
		secret:    os.Getenv("CRYPTO_SECRET"),
		secretApi: os.Getenv("CRYPTO_SECRET_API"),
	}
}

func (a *service) Encrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(a.secret))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a *service) Decrypt(text string) (string, error) {
	block, err := aes.NewCipher([]byte(a.secret))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func (a *service) EncryptApi(text string) (string, error) {
	block, err := aes.NewCipher([]byte(a.secretApi))
	if err != nil {
		return "", err
	}
	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func (a *service) DecryptApi(text string) (string, error) {
	block, err := aes.NewCipher([]byte(a.secretApi))
	if err != nil {
		return "", err
	}

	cipherText, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", err
	}

	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}
