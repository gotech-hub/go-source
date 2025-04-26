package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
)

func Encrypt(plaintext, secretKeyHex string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	iv := secretKey[:aes.BlockSize]

	bPlaintext := PKCS5Padding([]byte(plaintext), aes.BlockSize)

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, bPlaintext)

	// encode base64 and return
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertextBase64, secretKeyHex string) (string, error) {
	if ciphertextBase64 == "" {
		return "", nil
	}

	secretKey, err := hex.DecodeString(secretKeyHex)
	if err != nil {
		return "", err
	}

	ciphertextByte, err := base64.StdEncoding.DecodeString(ciphertextBase64)
	if err != nil {
		return "", err
	}

	iv := secretKey[:aes.BlockSize]

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		return "", err
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextByte, ciphertextByte)

	return string(PKCS5UnPadding(ciphertextByte)), nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	return src[:(length - unpadding)]
}
