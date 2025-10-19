package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/dkotsyuruba/go-shortener/internal/model"
	"github.com/dkotsyuruba/go-shortener/internal/service"
	"github.com/dkotsyuruba/go-shortener/internal/utils"
)

type Handler struct {
	service service.Service
	logger  *zap.Logger
}

func NewHandler(service service.Service, logger *zap.Logger) *Handler {
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

	userID, _ := utils.GetUserID(r)

	shortenURL, err := h.service.Shorten(string(originalURL), userID)
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
	w.Write([]byte(shortenURL))
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
	var request model.ShortenRequest
	err := decoder.Decode(&request)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Malformed input data", http.StatusBadRequest)
		return
	}

	userID, _ := utils.GetUserID(r)

	shortenedURL, err := h.service.Shorten(request.URL, userID)
	if err != nil && err != model.ErrDuplicatedURL {
		h.logger.Error(err.Error())
		http.Error(w, "Error shortening URL", http.StatusInternalServerError)
		return
	}

	result := model.ShortenResponse{
		Result: shortenedURL,
	}

	w.Header().Set("Content-Type", "application/json")
	if err == model.ErrDuplicatedURL {
		w.WriteHeader(http.StatusConflict)
	} else {
		w.WriteHeader(http.StatusCreated)
	}

	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	var batch []*model.BatchShortenRequest
	err := json.NewDecoder(r.Body).Decode(&batch)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Malformed input data", http.StatusBadRequest)
		return
	}

	userID, _ := utils.GetUserID(r)

	result, err := h.service.ShortenBatch(batch, userID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Error while saving links", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
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

func (h *Handler) GetUserURLs(w http.ResponseWriter, r *http.Request) {
	userID, ok := utils.GetUserID(r)
	if !ok {
		http.Error(w, "Unauthorized: user ID not found", http.StatusUnauthorized)
		return
	}

	result, err := h.service.GetAllByUserID(userID)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Error while getting links", http.StatusInternalServerError)
		return
	}

	if len(result) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error(err.Error())
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
