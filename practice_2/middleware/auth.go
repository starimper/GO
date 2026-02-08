package middleware

import (
	"encoding/json"
	"net/http"
)

const API_KEY = "secret12345"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("X-API-KEY") != API_KEY {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "unauthorized",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
