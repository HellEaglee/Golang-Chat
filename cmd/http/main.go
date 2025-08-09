package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/handler/http"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/logger"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
)

func main() {
	config, err := config.New()
	if err != nil {
		slog.Error("Error loading env variables", "error", err)
	}

	logger.Set(config.App)
	slog.Info("Starting the app", "app", config.App.Name, "env", config.App.Env)

	ctx := context.Background()
	db, err := postgres.New(ctx, config.DB)
	if err != nil {
		slog.Error("Error initializing DB connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	err = db.Migrate()
	if err != nil {
		slog.Error("Error migrating DB", "error", err)
		os.Exit(1)
	}

	router, err := http.NewRouter(
		config.HTTP,
	)
	if err != nil {
		slog.Error("Error initializing router", "error", err)
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%s", config.HTTP.URL, config.HTTP.Port)
	slog.Info("Starting the HTTP server", "listen_address", listenAddr)
	err = router.Serve(listenAddr)
	if err != nil {
		slog.Error("Error starting the HTTP server", "error", err)
		os.Exit(1)
	}
}
