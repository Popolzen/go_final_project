package service

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"io"
	"time"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/Popolzen/go_final_project/internal/storage"
	"github.com/google/uuid"
)

type secretService struct {
	repo   storage.Repository
	encKey []byte
}

func NewSecretService(repo storage.Repository, encKey []byte) SecretService {
	return &secretService{
		repo:   repo,
		encKey: encKey,
	}
}

func (s secretService) CreateSecret(ctx context.Context,
	userID uuid.UUID,
	secretType models.SecretType,
	name string,
	data any,
	metadata string,
) (*models.Secret, error) {

	if !secretType.Validate() {
		return nil, models.ErrInvalidSecretType
	}

	if name == "" {
		return nil, models.ErrEmptySecretName
	}

	raw, err := s.serializeSecretData(secretType, data)
	if err != nil {
		return nil, err
	}
	encrypted, err := s.encrypt(raw)
	if err != nil {
		return nil, err
	}

	secret := &models.Secret{
		ID:       uuid.New(),
		UserID:   userID,
		Type:     secretType,
		Name:     name,
		Data:     encrypted,
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

func (s secretService) GetSecret(ctx context.Context, secretID, userID uuid.UUID) (*models.Secret, interface{}, error) {

	secret, err := s.repo.GetSecret(ctx, secretID, userID)
	if err != nil {
		return nil, nil, err
	}
	raw, err := s.decrypt(secret.Data)
	if err != nil {
		return nil, nil, err
	}

	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, nil, err
	}

	return secret, data, nil

}

func (s secretService) GetSecretByName(ctx context.Context,
	name string, secretType models.SecretType, userID uuid.UUID) (*models.Secret, interface{}, error) {

	secret, err := s.repo.GetSecretByName(ctx, name, secretType, userID)
	if err != nil {
		return nil, nil, err
	}
	raw, err := s.decrypt(secret.Data)
	if err != nil {
		return nil, nil, err
	}

	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, nil, err
	}

	return secret, data, nil
}

func (s secretService) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*models.Secret, error) {
	return s.repo.ListSecrets(ctx, userID)
}

func (s *secretService) UpdateSecret(ctx context.Context, secretID, userID uuid.UUID, data interface{}, metadata string) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	encrypted, err := s.encrypt(raw)
	if err != nil {
		return err
	}
	secret := &models.Secret{
		ID:       secretID,
		UserID:   userID,
		Data:     encrypted,
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

func (s secretService) encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.encKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (s secretService) decrypt(ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(s.encKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce := ciphertext[:nonceSize]
	ciphertext = ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func (s secretService) serializeSecretData(secretType models.SecretType, data any) ([]byte, error) {
	switch secretType {

	case models.SecretTypeLogin:
		d, ok := data.(models.LoginData)
		if !ok {
			return nil, models.ErrInvalidSecretData
		}
		return json.Marshal(d)

	case models.SecretTypeText:
		d, ok := data.(models.TextData)
		if !ok {
			return nil, models.ErrInvalidSecretData
		}
		return json.Marshal(d)

	case models.SecretTypeBinary:
		d, ok := data.(models.BinaryData)
		if !ok {
			return nil, models.ErrInvalidSecretData
		}
		return json.Marshal(d)

	case models.SecretTypeCard:
		d, ok := data.(models.CardData)
		if !ok {
			return nil, models.ErrInvalidSecretData
		}
		return json.Marshal(d)

	default:
		return nil, models.ErrInvalidSecretType
	}
}
