package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"

	"be-ecommerce/internal/utils"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const (
	UserIDKey   contextKey = "user_id"
	UserRoleKey contextKey = "user_role"
)

// RequireAuth is a middleware that checks if the request has a valid JWT token.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Missing Authorization header", "unauthorized")
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid Authorization header format", "unauthorized")
			return
		}

		tokenString := parts[1]
		secretKey := os.Getenv("JWT_SECRET")
		if secretKey == "" {
			secretKey = "default_secret_key"
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrAbortHandler
			}
			return []byte(secretKey), nil
		})

		if err != nil || !token.Valid {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired token", "unauthorized")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid token claims", "unauthorized")
			return
		}

		userID := claims["user_id"].(string)
		role := claims["role"].(string)

		// Set user context
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		ctx = context.WithValue(ctx, UserRoleKey, role)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole is a middleware that checks if the authenticated user has a specific role.
// It must be used AFTER RequireAuth.
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value(UserRoleKey).(string)
			if !ok || userRole != role {
				utils.RespondError(w, http.StatusForbidden, "You do not have permission to access this resource", "forbidden")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
