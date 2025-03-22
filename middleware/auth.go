// middleware/auth.go
package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/BodaciousX/RVParkBackend/user"
)

type ContextKey string

const UserContextKey ContextKey = "user"

type AuthMiddleware struct {
	userService user.Service
}

func (m *AuthMiddleware) RevokeUserTokens(userID string) error {
	return m.userService.RevokeAllTokens(userID)
}

func NewAuthMiddleware(userService user.Service) *AuthMiddleware {
	return &AuthMiddleware{
		userService: userService,
	}
}

// RequireAuth middleware ensures that requests have a valid authentication token
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Validate token
		user, err := m.userService.ValidateToken(parts[1])
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Add user to request context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
