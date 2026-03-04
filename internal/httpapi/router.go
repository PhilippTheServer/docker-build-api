package httpapi

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(logger *slog.Logger) http.Handler {
	r := chi.NewRouter() // Create a new chi router instance

	// Add some common middleware to the router:
	r.Use(middleware.RequestID) // Generate a unique request ID for each request (useful for logging)
	r.Use(middleware.RealIP)    // Get the real client IP from X-Forwarded-For header (if behind a proxy)
	r.Use(middleware.Logger)    // Log each request with method, path, status, duration
	r.Use(middleware.Recoverer) // Recover from panics and return 500 Internal Server Error

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) { // Define a simple health check endpoint
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK")) // Write "OK" response body (ignore write errors for simplicity)
	})

	_ = logger //keep for now, we will use it later for logging in the handlers

	return r // Return the configured router as an http.Handler
}
