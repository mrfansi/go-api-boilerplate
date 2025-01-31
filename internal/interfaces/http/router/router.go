package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/mrfansi/go-api-boilerplate/internal/infrastructure/container"
	"github.com/mrfansi/go-api-boilerplate/internal/interfaces/http/handler"
	"github.com/mrfansi/go-api-boilerplate/internal/interfaces/http/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func NewRouter(c *container.Container) (http.Handler, error) {
	// Initialize router
	r := chi.NewRouter()

	// Inject dependencies
	var (
		authHandler      *handler.AuthHandler
		userHandler      *handler.UserHandler
		authMiddleware   *middleware.AuthMiddleware
		loggerMiddleware *middleware.LoggerMiddleware
		corsMiddleware   *middleware.CorsMiddleware
	)

	if err := c.Resolve(func(
		ah *handler.AuthHandler,
		uh *handler.UserHandler,
		am *middleware.AuthMiddleware,
		lm *middleware.LoggerMiddleware,
		cm *middleware.CorsMiddleware,
	) {
		authHandler = ah
		userHandler = uh
		authMiddleware = am
		loggerMiddleware = lm
		corsMiddleware = cm
	}); err != nil {
		return nil, err
	}

	// Global middleware
	r.Use(corsMiddleware.Cors)
	r.Use(loggerMiddleware.Logger)
	r.Use(httprate.LimitByIP(100, 1*60)) // 100 requests per minute per IP

	// Health check endpoint
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Metrics endpoint
	r.Handle("/metrics", promhttp.Handler())

	// Swagger documentation
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", authHandler.Login)
			r.Post("/users", userHandler.CreateUser) // Moved to public routes
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware.Authenticate)

			// Auth routes
			r.Post("/auth/refresh", authHandler.RefreshToken)

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Get("/", userHandler.ListUsers)

				r.Route("/{id}", func(r chi.Router) {
					r.Get("/", userHandler.GetUser)
					r.Put("/", userHandler.UpdateUser)
					r.Delete("/", userHandler.DeleteUser)
					r.Put("/password", userHandler.ChangePassword)

					// Admin only routes
					r.Group(func(r chi.Router) {
						r.Use(authMiddleware.RequireRole("admin"))
						r.Put("/role", userHandler.UpdateRole)
					})
				})
			})
		})
	})

	return r, nil
}
