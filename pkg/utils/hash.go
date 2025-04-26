package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
)

// Sha256Hash hash data, return base64 of hash
func Sha256Hash(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	hashedBytes := hasher.Sum(nil)
	return base64.StdEncoding.EncodeToString(hashedBytes)
}

// Sha256HashHex hash data, return hex string of hash
func Sha256HashHex(data []byte) string {
	hasher := sha256.New()
	hasher.Write(data)
	hashedBytes := hasher.Sum(nil)

	// return hex string
	return hex.EncodeToString(hashedBytes)
}
