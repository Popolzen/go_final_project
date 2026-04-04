package server

import (
	"context"
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
	Data     []byte            `json:"data"`
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
	Data      []byte            `json:"data"`
	Metadata  string            `json:"metadata,omitempty"`
	Version   int               `json:"version"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}

type UpdateSecretRequest struct {
	Data     []byte `json:"data"`
	Metadata string `json:"metadata,omitempty"`
}

type SecretSyncResponse struct {
	ID        string            `json:"id"`
	Type      models.SecretType `json:"type"`
	Name      string            `json:"name"`
	Data      []byte            `json:"data"`
	Metadata  string            `json:"metadata,omitempty"`
	Version   int               `json:"version"`
	UpdatedAt string            `json:"updated_at"`
	DeletedAt *string           `json:"deleted_at,omitempty"`
}

func (h *Handler) CreateSecret(w http.ResponseWriter, r *http.Request) {

	var req CreateSecretRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	userUUID, err := userIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	secret, err := h.secretService.CreateSecret(
		r.Context(),
		userUUID,
		req.Type,
		req.Name,
		req.Data,
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

func (h *Handler) GetSecret(w http.ResponseWriter, r *http.Request) {
	secretIDStr := chi.URLParam(r, "id")
	secretID, err := uuid.Parse(secretIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid secret id")
		return
	}
	userUUID, err := userIDFromContext(r.Context())
	secret, err := h.secretService.GetSecret(r.Context(), secretID, userUUID)
	if err != nil {
		handleSecretServiceError(w, err)
		return
	}

	resp := GetSecretResponse{
		ID:        secret.ID.String(),
		Type:      secret.Type,
		Name:      secret.Name,
		Data:      secret.Data,
		Metadata:  secret.Metadata,
		Version:   secret.Version,
		CreatedAt: secret.CreatedAt.Format(time.RFC3339),
		UpdatedAt: secret.UpdatedAt.Format(time.RFC3339),
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *Handler) GetSecretByName(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		respondError(w, http.StatusBadRequest, "secret name is required")
		return
	}
	secretType := models.SecretType(r.URL.Query().Get("type"))
	if !secretType.Validate() {
		respondError(w, http.StatusBadRequest, "invalid secret type")
		return
	}
	userUUID, err := userIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}
	secret, err := h.secretService.GetSecretByName(r.Context(), name, secretType, userUUID)

	if err != nil {
		handleSecretServiceError(w, err)
		return
	}
	resp := GetSecretResponse{
		ID:        secret.ID.String(),
		Type:      secret.Type,
		Name:      secret.Name,
		Data:      secret.Data,
		Metadata:  secret.Metadata,
		Version:   secret.Version,
		CreatedAt: secret.CreatedAt.Format(time.RFC3339),
		UpdatedAt: secret.UpdatedAt.Format(time.RFC3339),
	}

	respondJSON(w, http.StatusOK, resp)
}

func (h *Handler) ListSecrets(w http.ResponseWriter, r *http.Request) {
	userUUID, err := userIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}
	secrets, err := h.secretService.ListSecrets(r.Context(), userUUID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to list secrets")
		return
	}
	// Формируем ответ
	response := make([]SecretResponse, 0, len(secrets))
	for _, secret := range secrets {
		response = append(response, SecretResponse{
			ID:        secret.ID.String(),
			Type:      secret.Type,
			Name:      secret.Name,
			Metadata:  secret.Metadata,
			Version:   secret.Version,
			CreatedAt: secret.CreatedAt.Format(time.RFC3339),
			UpdatedAt: secret.UpdatedAt.Format(time.RFC3339),
		})
	}

	respondJSON(w, http.StatusOK, response)

}

func (h *Handler) UpdateSecret(w http.ResponseWriter, r *http.Request) {
	secretIDStr := chi.URLParam(r, "id")
	secretID, err := uuid.Parse(secretIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid secret id")
		return
	}
	// TODO
	userUUID, err := userIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	var req UpdateSecretRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Data) == 0 {
		respondError(w, http.StatusBadRequest, "data is required")
		return
	}

	err = h.secretService.UpdateSecret(
		r.Context(),
		secretID,
		userUUID,
		req.Data,
		req.Metadata,
	)
	if err != nil {
		handleSecretServiceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) DeleteSecret(w http.ResponseWriter, r *http.Request) {
	secretIDStr := chi.URLParam(r, "id")
	secretID, err := uuid.Parse(secretIDStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid secret id")
		return
	}
	userUUID, err := userIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}
	err = h.secretService.DeleteSecret(r.Context(), secretID, userUUID)
	if err != nil {
		handleSecretServiceError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) GetSecretsAfter(w http.ResponseWriter, r *http.Request) {
	userID, err := userIDFromContext(r.Context())
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid user")
		return
	}

	afterStr := r.URL.Query().Get("after")
	if afterStr == "" {
		respondError(w, http.StatusBadRequest, "`after` query param is required")
		return
	}

	after, err := time.Parse(time.RFC3339, afterStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid `after` time format")
		return
	}

	secrets, err := h.secretService.GetSecretsAfter(r.Context(), userID, after)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to get secrets")
		return
	}

	resp := make([]SecretSyncResponse, 0, len(secrets))
	for _, s := range secrets {
		item := SecretSyncResponse{
			ID:        s.ID.String(),
			Type:      s.Type,
			Name:      s.Name,
			Data:      s.Data,
			Metadata:  s.Metadata,
			Version:   s.Version,
			UpdatedAt: s.UpdatedAt.Format(time.RFC3339),
		}

		if s.DeletedAt != nil {
			t := s.DeletedAt.Format(time.RFC3339)
			item.DeletedAt = &t
		}

		resp = append(resp, item)
	}

	respondJSON(w, http.StatusOK, resp)
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

func userIDFromContext(ctx context.Context) (uuid.UUID, error) {
	raw, ok := ctx.Value(userIDKey).(string)
	if !ok {
		return uuid.Nil, errors.New("user id not found in context") //TODO
	}
	return uuid.Parse(raw)
}
