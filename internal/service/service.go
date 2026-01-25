package service

import (
	"context"
	"time"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AuthService interface {
	Register(ctx context.Context, username, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(tokenString string) (*Claims, error)
}

type SecretService interface {
	CreateSecret(ctx context.Context, userID uuid.UUID, secretType models.SecretType, name string, data interface{}, metadata string) (*models.Secret, error)
	GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, interface{}, error)
	GetSecretByName(ctx context.Context, name string, secretType models.SecretType, userID uuid.UUID) (*models.Secret, interface{}, error)
	ListSecrets(ctx context.Context, userID uuid.UUID) ([]*models.Secret, error)
	UpdateSecret(ctx context.Context, secretID, userID uuid.UUID, data interface{}, metadata string) error
	DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error
	GetSecretsAfter(ctx context.Context, userID uuid.UUID, after time.Time) ([]*models.Secret, error)
}

// Claims структура для JWT токена
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}
