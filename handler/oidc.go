package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/hrz8/simpath/config"
	"github.com/hrz8/simpath/jwt"
)

type JWK struct {
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	N   string `json:"n"`
	E   string `json:"e"`
}

type JWKs struct {
	Keys []JWK `json:"keys"`
}

func (h *Handler) JWKSHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	modulus, err := config.JWKSModulus()
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	exponent, err := config.JWKSExponent()
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	kid, err := config.JWKSKid()
	if err != nil {
		http.Error(w, "server_error", http.StatusInternalServerError)
		return
	}

	jwks := &JWKs{
		Keys: []JWK{
			{
				Alg: "RS256",
				Kty: "RSA",
				Kid: kid,
				Use: "sig",
				N:   modulus,
				E:   exponent,
			},
		},
	}
	json.NewEncoder(w).Encode(jwks)
}

type OIDCConfiguration struct {
	Issuer                     string   `json:"issuer"`
	AuthorizationEndpoint      string   `json:"authorization_endpoint"`
	TokenEndpoint              string   `json:"token_endpoint"`
	UserInfoEndpoint           string   `json:"userinfo_endpoint"`
	JwksURI                    string   `json:"jwks_uri"`
	ResponseTypesSupported     []string `json:"response_types_supported"`
	SubjectTypesSupported      []string `json:"subject_types_supported"`
	IDTokenSigningAlgSupported []string `json:"id_token_signing_alg_values_supported"`
}

func (h *Handler) OIDCConfig(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conf := &OIDCConfiguration{
		Issuer:                     config.JWTIssuer,
		AuthorizationEndpoint:      "http://localhost:5001/v1/oauth2/authorize",
		TokenEndpoint:              "http://localhost:5001/v1/oauth2/token",
		UserInfoEndpoint:           "http://localhost:5001/v1/oauth2/userinfo",
		JwksURI:                    "http://localhost:5001/v1/.well-known/jwks.json",
		ResponseTypesSupported:     []string{"code", "token", "id_token"},
		SubjectTypesSupported:      []string{"public"},
		IDTokenSigningAlgSupported: []string{"RS256"},
	}
	json.NewEncoder(w).Encode(conf)
}

func (h *Handler) UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}

	bearer := strings.TrimPrefix(authHeader, "Bearer ")

	pubKey, err := jwt.ReadPublicKey()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	token, err := jwt.VerifyJWT(bearer, pubKey)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	claims, err := jwt.GetClaimsJWT(token)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	subRaw := claims["sub"]
	sub, ok := subRaw.(string)
	if !ok || sub == "" {
		fmt.Println("ok", ok)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	usr, err := h.userSvc.FindUserByPublicID(sub)
	if err != nil {
		fmt.Println("err", err)
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(usr)
}
