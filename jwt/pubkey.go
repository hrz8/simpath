package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/hrz8/simpath/config"
)

func ReadPublicKey() (*rsa.PublicKey, error) {
	pubRaw, err := config.JWTPublicKey()
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode([]byte(pubRaw))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode public key")
	}
	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode public key: %v", err)
	}
	rsaPubKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("failed to read public key: not an RSA public key")
	}

	return rsaPubKey, nil
}
