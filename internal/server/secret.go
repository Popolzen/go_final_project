package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CreateSecretRequest struct {
	Type     models.SecretType `json:"type"`
	Name     string            `json:"name"`
	Data     json.RawMessage   `json:"data"`
	Metadata string            `json:"metadata,omitempty"`
}

type SecretResponse struct {
	ID        string            `json:"id"`
	Type      models.SecretType `json:"type"`
	Name      string            `json:"name"`
	Metadata  string            `json:"metadata,omitempty"`
	Version   int               `json:"version"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type GetSecretResponse struct {
	ID        string            `json:"id"`
	Type      models.SecretType `json:"type"`
	Name      string            `json:"name"`
	Data      any               `json:"data"`
	Metadata  string            `json:"metadata,omitempty"`
	Version   int               `json:"version"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {

	var req CreateSecretRequest

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// TODO: userID из JWT
	userUUID, err := uuid.Parse("7c30d978-b23e-466e-bbf4-7c6d07d2fbcc")
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	data, err := parseSecretData(req.Type, req.Data)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrInvalidSecretType):
			respondError(w, http.StatusBadRequest, "invalid secret type")
		case errors.Is(err, models.ErrInvalidSecretData):
			respondError(w, http.StatusBadRequest, "invalid secret data")
		default:
			respondError(w, http.StatusBadRequest, "invalid request")
		}
		return
	}

	secret, err := h.secretService.CreateSecret(
		r.Context(),
		userUUID,
		req.Type,
		req.Name,
		data,
		req.Metadata,
	)
	if err != nil {
		handleSecretServiceError(w, err)
		return
	}

	response := SecretResponse{
		ID:        secret.ID.String(),
		Type:      secret.Type,
		Name:      secret.Name,
		Metadata:  secret.Metadata,
		Version:   secret.Version,
		CreatedAt: secret.CreatedAt.Format(time.RFC3339),
		UpdatedAt: secret.UpdatedAt.Format(time.RFC3339),
	}

	respondJSON(w, http.StatusCreated, response)
}

func (h Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	secretIDStr := chi.URLParam(r, "id")
	secretID, err := uuid.Parse(secretIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid secret id")
		return
	}
	// TODO
	userUUID, err := uuid.Parse("7c30d978-b23e-466e-bbf4-7c6d07d2fbcc")
	secret, data, err := h.secretService.GetSecret(r.Context(), secretID, userUUID)
	if err != nil {
		handleSecretServiceError(w, err)
		return
	}

	resp := GetSecretResponse{
		ID:        secret.ID.String(),
		Type:      secret.Type,
		Name:      secret.Name,
		Data:      data,
		Metadata:  secret.Metadata,
		Version:   secret.Version,
		CreatedAt: secret.CreatedAt.Format(time.RFC3339),
		UpdatedAt: secret.UpdatedAt.Format(time.RFC3339),
	}

	respondJSON(w, http.StatusOK, resp)
}

func parseSecretData(
	secretType models.SecretType,
	raw json.RawMessage,
) (any, error) {

	switch secretType {

	case models.SecretTypeLogin:
		var d models.LoginData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, models.ErrInvalidSecretData
		}
		return d, nil

	case models.SecretTypeText:
		var d models.TextData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, models.ErrInvalidSecretData
		}
		return d, nil

	case models.SecretTypeBinary:
		var d models.BinaryData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, models.ErrInvalidSecretData
		}
		return d, nil

	case models.SecretTypeCard:
		var d models.CardData
		if err := json.Unmarshal(raw, &d); err != nil {
			return nil, models.ErrInvalidSecretData
		}
		return d, nil

	default:
		return nil, models.ErrInvalidSecretType
	}
}

func handleSecretServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrInvalidSecretType):
		respondError(w, http.StatusBadRequest, "invalid secret type")

	case errors.Is(err, models.ErrInvalidSecretData):
		respondError(w, http.StatusBadRequest, "invalid secret data")

	case errors.Is(err, models.ErrEmptySecretName):
		respondError(w, http.StatusBadRequest, "secret name cannot be empty")

	case errors.Is(err, models.ErrSecretNotFound):
		respondError(w, http.StatusNotFound, "secret not found")

	case errors.Is(err, models.ErrSecretAlreadyExists):
		respondError(w, http.StatusConflict, "secret already exists")

	default:
		respondError(w, http.StatusInternalServerError, "internal server error")
	}
}
