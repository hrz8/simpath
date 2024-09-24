package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var realm = "simpath"

func WriteJSON(w http.ResponseWriter, v interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(v)
}

func Error(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err})
}

func UnauthorizedError(w http.ResponseWriter, err string) {
	w.Header().Set("WWW-Authenticate", fmt.Sprintf("Bearer realm=%s", realm))
	Error(w, err, http.StatusUnauthorized)
}
