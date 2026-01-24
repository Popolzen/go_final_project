package server

import (
	"encoding/json"
	"net/http"

	"github.com/Popolzen/go_final_project/internal/service"
)

// Handler содержит зависимости для всех handlers
type Handler struct {
	authService   service.AuthService
	secretService service.SecretService
}

// NewHandler создаёт новый Handler с сервисами
func NewHandler(authService service.AuthService, secretService service.SecretService) *Handler {
	return &Handler{
		authService:   authService,
		secretService: secretService,
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
