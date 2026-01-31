package service

import (
	"context"

	"github.com/Popolzen/go_final_project/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) error
	Login(ctx context.Context, username, password string) error
}

type SecretService interface {
	Create(ctx context.Context, secretType models.SecretType, name, metadata string, data []byte) (*SecretInfo, error)
}

type SecretInfo struct {
	ID        string            `json:"id"`
	Type      models.SecretType `json:"type"`
	Name      string            `json:"name"`
	Metadata  string            `json:"metadata,omitempty"`
	Version   int               `json:"version"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
