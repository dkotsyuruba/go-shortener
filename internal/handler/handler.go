package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	Shorten(url string) (string, error)
	Get(id string) (string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	originalURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(originalURL) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortenURL, err := h.service.Shorten(string(originalURL))
	if err != nil {
		http.Error(w, "Error shortening URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortenURL)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenURL))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, err := h.service.Get(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
