package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sqlx.DB
}

func NewPostgresRepository(dsn string) (*PostgresRepository, error) {
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PostgresRepository{db: db}, nil
}

// CreateUser создаёт нового пользователя
func (r *PostgresRepository) CreateUser(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (username, password_hash)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, user.Username, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		// Проверка на дубликат username
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrUserExists
		}

		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

// GetUserByUsername получает пользователя по username
func (r *PostgresRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE username = $1`

	err := r.db.GetContext(ctx, &user, query, username)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

// GetUserByID получает пользователя по ID
func (r *PostgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *PostgresRepository) Ping(ctx context.Context) error {
	return r.db.PingContext(ctx)
}

func (r *PostgresRepository) Close() error {
	return r.db.Close()
}

func (r *PostgresRepository) CreateSecret(ctx context.Context, secret *models.Secret) error {
	query := `
		INSERT INTO secrets (user_id, type, name, data, metadata)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, version, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		secret.UserID, secret.Type, secret.Name, secret.Data, secret.Metadata).
		Scan(&secret.ID, &secret.Version, &secret.CreatedAt, &secret.UpdatedAt)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return models.ErrSecretAlreadyExists
		}

		return fmt.Errorf("failed to create secret: %w", err)
	}

	return nil
}

func (r *PostgresRepository) UpdateSecret(ctx context.Context, secret *models.Secret) error {
	query := `
		UPDATE secrets
		SET data = $1, metadata = $2, version = version + 1
		WHERE id = $3 AND user_id = $4 AND deleted_at IS NULL
		RETURNING version, updated_at`

	err := r.db.QueryRowContext(ctx, query,
		secret.Data, secret.Metadata, secret.ID, secret.UserID).
		Scan(&secret.Version, &secret.UpdatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrSecretNotFound
	}

	if err != nil {
		return fmt.Errorf("failed to update secret: %w", err)
	}
	return nil
}

func (r *PostgresRepository) DeleteSecret(ctx context.Context, id, userID uuid.UUID) error {
	query := `
		UPDATE secrets
		SET deleted_at = NOW()
		WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL
	`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete secret: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rows == 0 {
		return models.ErrSecretNotFound
	}

	return nil
}

// GetUserByUsername получает пользователя по username
func (r *PostgresRepository) GetSecret(ctx context.Context, id, userID uuid.UUID) (*models.Secret, error) {
	var secret models.Secret

	query := `SELECT * FROM secrets WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &secret, query, id, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrSecretNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get secret: %w", err)
	}

	return &secret, nil
}

func (r *PostgresRepository) GetSecretByName(ctx context.Context, name string, secretType models.SecretType, userID uuid.UUID) (*models.Secret, error) {
	var secret models.Secret

	query := `SELECT * FROM secrets WHERE name = $1 AND type = $2 AND user_id = $3 AND deleted_at IS NULL`

	err := r.db.GetContext(ctx, &secret, query, name, secretType, userID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrSecretNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get secret by name: %w", err)
	}

	return &secret, nil

}

func (r *PostgresRepository) ListSecrets(ctx context.Context, userID uuid.UUID) ([]*models.Secret, error) {
	var secrets []*models.Secret
	query := `SELECT * FROM secrets WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &secrets, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to list secrets: %w", err)
	}

	return secrets, nil
}

func (r *PostgresRepository) GetSecretsAfter(ctx context.Context, userID uuid.UUID, after time.Time) ([]*models.Secret, error) {
	var secrets []*models.Secret
	query := `SELECT * FROM secrets WHERE user_id = $1 AND updated_at > $2 ORDER BY updated_at`

	err := r.db.SelectContext(ctx, &secrets, query, userID, after)
	if err != nil {
		return nil, fmt.Errorf("failed to get secrets after: %w", err)
	}

	return secrets, nil
}
