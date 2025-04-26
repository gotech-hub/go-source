package utils

import (
	"crypto/rsa"
	"crypto/x509"
)

func GetRSAPublicKeyFromPrivateKey(privateKey *rsa.PrivateKey) (*rsa.PublicKey, error) {
	// Serialize the private key to DER format
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)

	// Parse the private key to get the public key
	parsedPrivateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	publicKey := &parsedPrivateKey.PublicKey
	return publicKey, nil
}
