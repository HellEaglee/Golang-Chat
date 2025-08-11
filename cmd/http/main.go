package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/HellEaglee/Golang-Chat/docs"
	jwt "github.com/HellEaglee/Golang-Chat/internal/adapter/auth/JWT"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/config"
	httphandler "github.com/HellEaglee/Golang-Chat/internal/adapter/handler/http"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/logger"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres"
	"github.com/HellEaglee/Golang-Chat/internal/adapter/storage/postgres/repository"
	"github.com/HellEaglee/Golang-Chat/internal/core/service"
)

// @title			Golang Chat API
// @version		1.0
// @description	This is a chat application API.
// @termsOfService	http://swagger.io/terms/
//
// @host			localhost:8080
// @BasePath		/v1
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

	// DI
	tokenRepo := repository.NewTokenRepository(db)
	token, err := jwt.New(config.Token, tokenRepo)
	if err != nil {
		slog.Error("Error initializing token service", "error", err)
		os.Exit(1)
	}

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	userHandler := httphandler.NewUserHandler(userService)

	authService := service.NewAuthService(userRepo, token)
	authHandler := httphandler.NewAuthHandler(authService)

	router, err := httphandler.NewRouter(
		config.HTTP,
		*authHandler,
		*userHandler,
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
