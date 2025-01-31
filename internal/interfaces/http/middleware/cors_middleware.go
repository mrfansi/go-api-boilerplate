package middleware

import (
	"net/http"
	"strings"

	"github.com/mrfansi/go-api-boilerplate/internal/infrastructure/config"
)

type CorsMiddleware struct {
	config *config.Config
}

func NewCorsMiddleware(cfg *config.Config) *CorsMiddleware {
	return &CorsMiddleware{
		config: cfg,
	}
}

// Cors handles CORS (Cross-Origin Resource Sharing)
func (m *CorsMiddleware) Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get origin from request header
		origin := r.Header.Get("Origin")

		// Check if origin is allowed
		allowed := false
		for _, allowedOrigin := range m.config.Cors.AllowedOrigins {
			if allowedOrigin == "*" || allowedOrigin == origin {
				allowed = true
				break
			}
		}

		if allowed {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", strings.Join(m.config.Cors.AllowedMethods, ","))
			w.Header().Set("Access-Control-Allow-Headers", strings.Join(m.config.Cors.AllowedHeaders, ","))
			w.Header().Set("Access-Control-Expose-Headers", strings.Join(m.config.Cors.ExposedHeaders, ","))

			if m.config.Cors.AllowCredentials {
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			if m.config.Cors.MaxAge > 0 {
				w.Header().Set("Access-Control-Max-Age", string(m.config.Cors.MaxAge))
			}
		}

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
