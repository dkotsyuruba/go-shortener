package main

import (
	"log"
	"net/http"

	"github.com/dkotsyuruba/go-shortener/internal/config"
	"github.com/dkotsyuruba/go-shortener/internal/handler"
	"github.com/dkotsyuruba/go-shortener/internal/repository"
	"github.com/dkotsyuruba/go-shortener/internal/service"
	"github.com/dkotsyuruba/go-shortener/pkg/shortener"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.InitConfig()
	repo := repository.NewRepository()
	shortener := shortener.NewRealShortenerService()
	srv := service.NewService(repo, cfg.Service, shortener)
	handlers := handler.NewHandler(srv)

	router := chi.NewRouter()
	router.Post("/", handlers.Shorten)
	router.Get("/{id}", handlers.Get)

	log.Fatal(http.ListenAndServe(cfg.Server.Port, router))
}
