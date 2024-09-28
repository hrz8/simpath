package jwt

import (
	"crypto/rsa"
	"errors"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(tokenString string, publicKey *rsa.PublicKey) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("error parsing token: %v", err)
	}

	return token, nil
}

func GetClaimsJWT(token *jwt.Token) (jwt.MapClaims, error) {
	claims, claimOk := token.Claims.(jwt.MapClaims)
	if !claimOk || !token.Valid {
		return nil, errors.New("token invalid")
	}

	return claims, nil
}
