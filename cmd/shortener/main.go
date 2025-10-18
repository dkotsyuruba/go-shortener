package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"go.uber.org/zap"

	"github.com/dkotsyuruba/go-shortener/internal/config"
	"github.com/dkotsyuruba/go-shortener/internal/handler"
	"github.com/dkotsyuruba/go-shortener/internal/middleware"
	"github.com/dkotsyuruba/go-shortener/internal/repository"
	"github.com/dkotsyuruba/go-shortener/internal/service"
	"github.com/dkotsyuruba/go-shortener/pkg/shortener"
	"github.com/go-chi/chi/v5"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	cfg := config.InitConfig()
	repo, err := repository.NewRepository(cfg)
	if err != nil && !os.IsNotExist(err) {
		logger.Fatal("repository initialization failed", zap.Error(err))
	}

	shortener := shortener.NewRealShortenerService()
	srv := service.NewService(repo, cfg.Service, shortener)
	handlers := handler.NewHandler(srv, logger)

	router := chi.NewRouter()
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.GzipMiddleware)
	router.Get("/{id}", handlers.Get)
	router.Get("/ping", handlers.Ping)
	router.Post("/", handlers.Shorten)
	router.Post("/api/shorten", handlers.ShortenJSON)
	router.Post("/api/shorten/batch", handlers.ShortenBatch)

	server := &http.Server{
		Addr:         cfg.Server.Address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("starting server error", zap.Error(err))
		}
	}()

	logger.Info("server started successfully at " + cfg.Server.Address)

	<-stopChan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("shutdown error", zap.Error(err))
	}

	if err := repo.Close(); err != nil {
		logger.Fatal("repository shut down failed", zap.Error(err))
	}

	logger.Info("server shut down successfully")
}
