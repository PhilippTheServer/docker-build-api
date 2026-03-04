package app // Declares this file belongs to package "app" (an importable library)

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/philipptheserver/docker-build-api/internal/httpapi"
)

// type introduces a new named type in Go.
// defines a struct named Config to hold configuration values.
type Config struct {
	Addr string // Addr is a field of type string, e.g. ":8080" or "localhost:8080"
}

// functions that start by a captial letter are exported and can be
// used by other packages that import this package.
// Run is the main function to start the application. It takes a
// Config and a logger as parameters.

func Run(cfg Config, logger *slog.Logger) error { // Function to start the app
	// Parameters:
	// - cfg: a Config struct containing configuration values (e.g., server address)
	// - logger: a pointer to a slog.Logger for logging messages
	// Returns:
	// - error: nil on success, non-nil on failure

	// := short variable declaration.
	handler := httpapi.NewRouter(logger) // Build the router + middleware used to handle HTTP requests

	// &http.Server{...} creates a new http.Server struct value and takes its address.
	// In Go, &T{...} is a very common way to allocate and get a pointer to a struct.
	httpServer := &http.Server{ // Configure the HTTP server instance
		Addr:    cfg.Addr, // Bind address/port the server listens on (e.g. ":8080")
		Handler: handler,  // The handler that serves requests (router + middleware + endpoints)

		// These timeouts protect against slow clients and resource exhaustion (important in production).
		ReadTimeout:       10 * time.Second, // Max time to read the entire request (headers + body)
		ReadHeaderTimeout: 5 * time.Second,  // Max time to read request headers (slowloris protection)
		WriteTimeout:      30 * time.Second, // Max time to write the response (protects from slow readers)
		IdleTimeout:       60 * time.Second, // Max time to keep idle keep-alive connections open
	}

	errCh := make(chan error, 1) // Channel to receive a fatal server error from the background goroutine

	// go func() { ... }() starts an anonymous function as a new goroutine.
	// A goroutine runs concurrently with the caller.
	go func() { // Run the server in the background so we can also wait for signals in the main goroutine
		logger.Info("listening", "addr", cfg.Addr) // Log that we are starting to listen (structured key/value)

		// ListenAndServe starts the HTTP server loop (accept connections, handle requests).
		// It blocks until the server stops. On normal shutdown it returns http.ErrServerClosed.
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			// "if err := ...; condition" is an if-statement with an initializer.
			// It declares err scoped to the if block.

			// errors.Is checks if err is (or wraps) a target error; safer than err == target due to wrapping.
			errCh <- err // Send the error to errCh (may block if buffer full; buffer=1 avoids most issues here)
		}
	}() // The trailing () calls the anonymous function immediately (in the new goroutine)

	stop := make(chan os.Signal, 1) // Channel that will receive shutdown signals

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM) // Ask OS to deliver SIGINT/SIGTERM notifications to 'stop'

	// select waits until one of its cases can proceed.
	// This is Go's way to wait on multiple channels at once.
	select { // Wait until either we receive a shutdown signal or the server fails unexpectedly
	case <-stop: // Receive from stop (we ignore the actual signal value here)
		// If we wanted the actual signal: sig := <-stop
	case err := <-errCh: // Receive an error from the server goroutine
		return err // Propagate the server error to caller (main), causing non-zero exit
	}

	// context.WithTimeout creates a derived context that automatically cancels after the duration.
	// We use it to bound how long graceful shutdown is allowed to take.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // Create shutdown context with 10s deadline

	// defer schedules a function call to run when the surrounding function returns.
	// This ensures we always release resources (cancel the context timer) even on early returns.
	defer cancel() // Stop the timeout timer and release context resources when Run returns

	// Shutdown gracefully stops the server:
	// - stops accepting new connections
	// - waits for in-flight requests to finish until ctx deadline
	// - then force-closes remaining connections
	return httpServer.Shutdown(ctx) // Perform graceful shutdown and return any shutdown error (nil on clean shutdown)
}
