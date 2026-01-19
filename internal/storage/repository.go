package storage

import (
	"context"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/google/uuid"
)

type Repository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)

	// CreateSecret(ctx context.Context, secret *models.Secret) error
	// UpdateSecret(ctx context.Context, secret *models.Secret) error
	// DeleteSecret(ctx context.Context, id, userID uuid.UUID) error
	// GetSecret(ctx context.Context, id, userID uuid.UUID) (*models.Secret, error)
	// GetSecretByName(ctx context.Context, name string, secretType models.SecretType, userID uuid.UUID) (*models.Secret, error)
	// ListSecrets(ctx context.Context, userID uuid.UUID) ([]*models.Secret, error)
	// GetSecretsAfter(ctx context.Context, userID uuid.UUID, after time.Time) ([]*models.Secret, error)

	Ping(ctx context.Context) error
	Close() error
}
