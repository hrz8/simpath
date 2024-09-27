package jwt

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
)

func ReadPrivateKey() (*rsa.PrivateKey, error) {
	pkRaw := os.Getenv("JWT_PRIVATE_KEY")
	if pkRaw == "" {
		return nil, errors.New("missing private key")
	}

	pkRaw = restoreNewlines(pkRaw)
	block, _ := pem.Decode([]byte(pkRaw))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("failed to decode private key")
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to decode private key: %v", err)
	}
	rsaPrivKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("failed to read public key: not an RSA public key")
	}

	return rsaPrivKey, nil
}

func restoreNewlines(s string) string {
	return strings.ReplaceAll(s, "\\n", "\n")
}
