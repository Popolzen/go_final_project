package server

import (
	"encoding/json"
	"net/http"

	"github.com/Popolzen/go_final_project/internal/storage"
)

// Handler содержит зависимости для всех handlers
type Handler struct {
	repo      storage.Repository
	jwtSecret []byte
	encKey    []byte
}

// NewHandler создаёт новый Handler
func NewHandler(repo storage.Repository, jwtSecret, encKey []byte) *Handler {
	return &Handler{
		repo:      repo,
		jwtSecret: jwtSecret,
		encKey:    encKey,
	}
}

// respondJSON отправляет JSON ответ
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// respondError отправляет ошибку
func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"error": message})
}
