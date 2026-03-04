package main // declares that this is the entry point of the application

import (
	"log/slog"
	"os"

	"github.com/philipptheserver/docker-build-api/internal/app"
)

func main() { // needed as an entry point. Must be provided
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger) // set gloabl logger

	cfg := app.Config{Addr: ":8080"} // set Port of the server

	if err := app.Run(cfg, logger); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	} // if server returns an error, log and exit with code 1
}
