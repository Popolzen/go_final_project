package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Popolzen/go_final_project/internal/models"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// Register регистрирует нового пользователя
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Register(r.Context(), req.Username, req.Password)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidUsername):
			respondError(w, http.StatusBadRequest, "username must be at least 3 characters")
		case errors.Is(err, models.ErrInvalidPassword):
			respondError(w, http.StatusBadRequest, "password must be at least 8 characters")
		case errors.Is(err, models.ErrUserExists):
			respondError(w, http.StatusConflict, "username already exists")
		default:
			respondError(w, http.StatusInternalServerError, "failed to create user")
		}
		return
	}

	respondJSON(w, http.StatusCreated, AuthResponse{Token: token})
}

// Login аутентифицирует пользователя
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	token, err := h.authService.Login(r.Context(), req.Username, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCredentials) {
			respondError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		respondError(w, http.StatusInternalServerError, "login failed")
		return
	}

	respondJSON(w, http.StatusOK, AuthResponse{Token: token})
}
