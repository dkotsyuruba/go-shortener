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
	repo := repository.NewRepository(cfg.Service.FileStorage)
	loadedLinks, err := repo.LoadFromFile(cfg.Service.FileStorage)
	if err != nil {
		logger.Fatal("loading data from file failed", zap.Error(err))
	}
	logger.Info("loaded", zap.Int("count", len(loadedLinks)), zap.String("filename", cfg.Service.FileStorage))

	shortener := shortener.NewRealShortenerService()
	srv := service.NewService(repo, cfg.Service, shortener)
	handlers := handler.NewHandler(srv)

	router := chi.NewRouter()
	router.Use(middleware.LoggerMiddleware(logger))
	router.Use(middleware.GzipMiddleware)
	router.Post("/", handlers.Shorten)
	router.Get("/{id}", handlers.Get)
	router.Post("/api/shorten", handlers.ShortenJSON)

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

	if err := repo.Persist(cfg.Service.FileStorage); err != nil {
		logger.Fatal("saving data to file failed", zap.Error(err))
	}

	logger.Info("server shut down successfully")
}
