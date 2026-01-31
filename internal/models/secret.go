package models

import (
	"time"

	"github.com/google/uuid"
)

type Secret struct {
	ID        uuid.UUID  `db:"id" json:"id"`
	UserID    uuid.UUID  `db:"user_id" json:"user_id"`
	Type      SecretType `db:"type" json:"type"`
	Name      string     `db:"name" json:"name"`
	Data      []byte     `db:"data" json:"data"`
	Metadata  string     `db:"metadata" json:"metadata,omitempty"`
	Version   int        `db:"version" json:"version"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt time.Time  `db:"updated_at" json:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}

func (s *Secret) Validate() error {
	if !s.Type.Validate() {
		return ErrInvalidSecretType
	}
	if len(s.Name) == 0 {
		return ErrEmptySecretName
	}
	if len(s.Data) == 0 {
		return ErrEmptySecretData
	}
	return nil
}

func (s *Secret) IsDeleted() bool {
	return s.DeletedAt != nil
}

type LoginData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TextData struct {
	Content string `json:"content"`
}

type BinaryData struct {
	Filename string `json:"filename"`
	Content  string `json:"content"` // base64
}

type CardData struct {
	Number string `json:"number"`
	Holder string `json:"holder"`
	CVV    string `json:"cvv"`
	Expiry string `json:"expiry"`
}
