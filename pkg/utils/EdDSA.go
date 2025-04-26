package utils

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"github.com/pkg/errors"
)

func CheckPublicKeyEdDSA(publicKeyHex string) error {
	publicKeyBytes, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return err
	}

	if len(publicKeyBytes) != ed25519.PublicKeySize {
		return fmt.Errorf("publicKey must have a size 32 bytes")
	}

	return nil
}

func EdDSAVerifySignature(message, sigHex, publicKeyHex string) (bool, error) {
	publicKey, err := hex.DecodeString(publicKeyHex)
	if err != nil {
		return false, errors.Wrap(err, "decoding public key error")
	}

	signature, err := hex.DecodeString(sigHex)
	if err != nil {
		return false, errors.Wrap(err, "decoding signature error")
	}

	// Verify the signature
	if ed25519.Verify(publicKey, []byte(message), signature) {
		return true, nil
	}

	return false, nil
}
