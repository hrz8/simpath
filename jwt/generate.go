package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"time"

	JWT "github.com/golang-jwt/jwt/v5"
	"github.com/hrz8/simpath/config"
)

type AccessTokenClaims struct {
	Aud         string
	Sub         string
	Scope       string
	ClientID    string
	Permissions []string
	Roles       []string
}

func GenerateAccessToken(jwtID string, data AccessTokenClaims) (string, error) {
	claims := JWT.MapClaims{
		"iss":         config.JWTIssuer,
		"sub":         data.Sub,
		"aud":         data.Aud,
		"iat":         time.Now().Unix(),
		"exp":         time.Now().Add(config.AccessTokenLifetime * time.Second).Unix(),
		"scope":       data.Scope,
		"client_id":   data.ClientID,
		"jti":         jwtID,
		"permissions": data.Permissions,
		"roles":       data.Roles,
	}
	secret, err := ReadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("error generate access token: %+v", err)
	}
	return generateJWT(claims, secret)
}

type RefreshTokenClaims struct {
	Aud      string
	Sub      string
	Scope    string
	ClientID string
}

func GenerateRefreshToken(jwtID string, data RefreshTokenClaims) (string, error) {
	claims := JWT.MapClaims{
		"iss":       config.JWTIssuer,
		"sub":       data.Sub,
		"aud":       data.Aud,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(config.RefreshTokenLifetime * time.Second).Unix(),
		"scope":     data.Scope,
		"client_id": data.ClientID,
		"jti":       jwtID,
	}
	secret, err := ReadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("error generate refresh token: %+v", err)
	}
	return generateJWT(claims, secret)
}

type IDTokenClaims struct {
	Aud      string
	Sub      string
	Scope    string
	ClientID string
	AuthTime int64
	Name     string
	Email    string
	Picture  string
}

func GenerateIDToken(jwtID string, data IDTokenClaims) (string, error) {
	claims := JWT.MapClaims{
		"iss":       config.JWTIssuer,
		"sub":       data.Sub,
		"aud":       data.Aud,
		"iat":       time.Now().Unix(),
		"exp":       time.Now().Add(config.IDTokenLifetime * time.Second).Unix(),
		"auth_time": data.AuthTime,
		"nonce":     "somenonce",
		"amr":       "password",
		"name":      data.Name,
		"email":     data.Email,
		"picture":   data.Picture,
	}
	secret, err := ReadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("error generate refresh token: %+v", err)
	}
	return generateJWT(claims, secret)
}

func generateJWT(claims JWT.MapClaims, privateKey *rsa.PrivateKey) (string, error) {
	kid := os.Getenv("JWKS_KID")
	if kid == "" {
		return "", errors.New("missing kid")
	}

	token := JWT.NewWithClaims(JWT.SigningMethodRS256, claims)
	token.Header["kid"] = kid

	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
