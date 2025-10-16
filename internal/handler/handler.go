package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/dkotsyuruba/go-shortener/internal/model"
)

type Service interface {
	Shorten(url string) (string, error)
	Get(id string) (string, error)
	Ping() error
	ShortenBatch(batch []*model.BatchShortenRequest) ([]*model.BatchShortenResponse, error)
}

type Handler struct {
	service Service
	logger  *zap.Logger
}

func NewHandler(service Service, logger *zap.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	originalURL, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil || len(originalURL) == 0 {
		h.logger.Error(err.Error())
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	shortenURL, err := h.service.Shorten(string(originalURL))
	if err != nil && err != model.ErrDuplicatedURL {
		h.logger.Error(err.Error())
		http.Error(w, "Error shortening URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Header().Set("Content-Length", strconv.Itoa(len(shortenURL)))
	if err == model.ErrDuplicatedURL {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	w.Write([]byte(string(shortenURL)))
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	originalURL, err := h.service.Get(id)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", originalURL)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (h *Handler) ShortenJSON(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	input := model.ShortenRequest{}
	err := decoder.Decode(&input)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Malformed input data", http.StatusBadRequest)
		return
	}

	shortenedURL, err := h.service.Shorten(input.URL)
	if err != nil && err != model.ErrDuplicatedURL {
		h.logger.Error(err.Error())
		http.Error(w, "Error shortening URL", http.StatusInternalServerError)
		return
	}

	output := model.ShortenResponse{
		Result: shortenedURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err == model.ErrDuplicatedURL {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
	json.NewEncoder(w).Encode(output)
}

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	var batch []*model.BatchShortenRequest
	err := json.NewDecoder(r.Body).Decode(&batch)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Malformed input data", http.StatusBadRequest)
		return
	}

	result, err := h.service.ShortenBatch(batch)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Error while saving links", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(result)
}

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.service.Ping()
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}
