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
	List(ctx context.Context) ([]*SecretInfo, error)
	GetByName(ctx context.Context, name string, secretType models.SecretType) (*SecretDetail, error)
	Update(ctx context.Context, name string, secretType models.SecretType, metadata string, data []byte) error
	Delete(ctx context.Context, name string, secretType models.SecretType) error
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

type SecretDetail struct {
	ID        string            `json:"id"`
	Type      models.SecretType `json:"type"`
	Name      string            `json:"name"`
	Data      []byte            `json:"data"`
	Metadata  string            `json:"metadata,omitempty"`
	Version   int               `json:"version"`
	CreatedAt string            `json:"created_at"`
	UpdatedAt string            `json:"updated_at"`
}
