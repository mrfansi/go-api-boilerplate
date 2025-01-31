package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mrfansi/go-api-boilerplate/internal/application/service"
	"github.com/mrfansi/go-api-boilerplate/internal/domain/errors"
)

// contextKey is a custom type for context keys to avoid collisions
type contextKey string

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

// Authenticate verifies the JWT token in the request header
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := extractToken(r)
		if token == "" {
			respondWithError(w, http.StatusUnauthorized, errors.ErrUnauthorized)
			return
		}

		jwtToken, err := m.authService.ValidateToken(token)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, err)
			return
		}

		claims, ok := jwtToken.Claims.(jwt.MapClaims)
		if !ok {
			respondWithError(w, http.StatusUnauthorized, errors.ErrInvalidToken)
			return
		}

		// Add claims to request context
		ctx := context.WithValue(r.Context(), contextKey("claims"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireRole checks if the authenticated user has the required role
func (m *AuthMiddleware) RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(contextKey("claims")).(jwt.MapClaims)
			if !ok {
				respondWithError(w, http.StatusUnauthorized, errors.ErrUnauthorized)
				return
			}

			userRole, ok := claims["role"].(string)
			if !ok || userRole != role {
				respondWithError(w, http.StatusForbidden, errors.ErrForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// Helper functions

func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	response := errors.NewErrorResponse(code, err.Error())
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(response)
}
