// middleware/cors.go
package middleware

import (
	"net/http"
	"os"
	"strings"
)

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := strings.Split(os.Getenv("CORS_ORIGIN"), ",")
		origin := r.Header.Get("Origin")

		// Check if the request origin is in the allowed origins
		allowedOrigin := ""
		for _, allowed := range allowedOrigins {
			if allowed == origin || allowed == "*" {
				allowedOrigin = origin
				break
			}
		}

		// If no match found and no origins configured, use development default
		if allowedOrigin == "" && len(allowedOrigins) == 0 {
			allowedOrigin = "http://localhost:3000"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
