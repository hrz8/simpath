package config

import (
	"errors"
	"os"
	"strings"
)

const (
	DatabaseURL = "postgres://postgres:toor@localhost:5432/simpath?sslmode=disable"
	AutoMigrate = true
	AllowClient = "http://localhost:8088"
	// token
	AccessTokenLifetime  = 7200    // 2 hours
	RefreshTokenLifetime = 1209600 // 14 days
	IDTokenLifetime      = 7200
	AuthCodeLifetime     = 7200
	// session
	SessionSecretKey = "kcp?l2Qh39{89Wq2"
	SessionPath      = "/"
	// SessionMaxAge       = 604800 // 7 days
	SessionMaxAge       = 3600 // 1 hour
	SessionHttpOnly     = true
	SessionName         = "simpath_session"
	UserSessionName     = "simpath_session_user"
	UserDataSessionKey  = "userdata"
	CSRFTokenSessionKey = "csrftoken"
	// JWT
	JWTIssuer          = "http://localhost:5001"
	JWTAccessTokenAud  = "resource-server-xyz"
	JWTRefreshTokenAud = "http://localhost:5001"
)

func JWTPrivateKey() (string, error) {
	pkRaw := os.Getenv("JWT_PRIVATE_KEY")
	if pkRaw == "" {
		return "", errors.New("missing private key")
	}
	pkRaw = restoreNewlines(pkRaw)
	return pkRaw, nil
}

func JWTPublicKey() (string, error) {
	pubRaw := os.Getenv("JWT_PUBLIC_KEY")
	if pubRaw == "" {
		return "", errors.New("missing public key")
	}
	pubRaw = restoreNewlines(pubRaw)
	return pubRaw, nil
}

func JWKSKid() (string, error) {
	kid := os.Getenv("JWKS_KID")
	if kid == "" {
		return "", errors.New("missing jwks kid")
	}
	return kid, nil
}

func JWKSModulus() (string, error) {
	modulus := os.Getenv("JWKS_MODULUS")
	if modulus == "" {
		return "", errors.New("missing jwks modulus")
	}
	return modulus, nil
}

func JWKSExponent() (string, error) {
	n := os.Getenv("JWKS_EXPONENT")
	if n == "" {
		return "", errors.New("missing jwks exponent")
	}
	return n, nil
}

func restoreNewlines(s string) string {
	return strings.ReplaceAll(s, "\\n", "\n")
}
