package service

import (
	"context"
	"time"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/Popolzen/go_final_project/internal/storage"
	"github.com/google/uuid"
)

type secretService struct {
	repo storage.Repository
}

func NewSecretService(repo storage.Repository) SecretService {
	return &secretService{
		repo: repo,
	}
}

func (s secretService) CreateSecret(ctx context.Context,
	userID uuid.UUID,
	secretType models.SecretType,
	name string,
	data []byte,
	metadata string,
) (*models.Secret, error) {

	if !secretType.Validate() {
		return nil, models.ErrInvalidSecretType
	}

	if name == "" {
		return nil, models.ErrEmptySecretName
	}

	secret := &models.Secret{
		ID:       uuid.New(),
		UserID:   userID,
		Type:     secretType,
		Name:     name,
		Data:     data,
		Metadata: metadata,
		Version:  1,
	}
	if err := secret.Validate(); err != nil {
		return nil, err
	}

	if err := s.repo.CreateSecret(ctx, secret); err != nil {
		return nil, err
	}

	return secret, nil
}

func (s secretService) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, error) {

	secret, err := s.repo.GetSecret(ctx, secretID, userID)
	if err != nil {
		return nil, err
	}

	return secret, nil

}

func (s secretService) GetSecretByName(ctx context.Context,
	name string, secretType models.SecretType, userID uuid.UUID) (*models.Secret, error) {

	secret, err := s.repo.GetSecretByName(ctx, name, secretType, userID)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func (s secretService) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*models.Secret, error) {
	return s.repo.ListSecrets(ctx, userID)
}

func (s *secretService) UpdateSecret(ctx context.Context, secretID, userID uuid.UUID, data []byte, metadata string) error {

	secret := &models.Secret{
		ID:       secretID,
		UserID:   userID,
		Data:     data,
		Metadata: metadata,
	}

	return s.repo.UpdateSecret(ctx, secret)
}

func (s *secretService) DeleteSecret(ctx context.Context, secretID, userID uuid.UUID) error {
	return s.repo.DeleteSecret(ctx, secretID, userID)
}

func (s *secretService) GetSecretsAfter(ctx context.Context, userID uuid.UUID, after time.Time) ([]*models.Secret, error) {
	return s.repo.GetSecretsAfter(ctx, userID, after)
}
