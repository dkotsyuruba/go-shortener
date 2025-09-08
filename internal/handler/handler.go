package handler

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Service interface {
	Shorten(Url string) (string, error)
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
	originalUrl, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(originalUrl) == 0 {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortenUrl, err := h.service.Shorten(string(originalUrl))
	if err != nil {
		http.Error(w, "Error shortening Url", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortenUrl)))
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortenUrl))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalUrl, err := h.service.Get(id)
	if err != nil {
		http.Error(w, "Not found", http.StatusBadRequest)
		return
	}

	w.Header().Set("Location", originalUrl)
	w.WriteHeader(http.StatusTemporaryRedirect)
}
