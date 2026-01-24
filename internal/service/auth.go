package service

import (
	"context"
	"time"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/Popolzen/go_final_project/internal/storage"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo      storage.Repository
	jwtSecret []byte
}

func NewAuthService(repo storage.Repository, jwtSecret []byte) AuthService {
	return authService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

func (s authService) Login(ctx context.Context, username, password string) (string, error) {
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil || user == nil {
		return "", models.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", models.ErrInvalidCredentials
	}

	return s.generateToken(user)
}

func (s authService) Register(ctx context.Context, username, password string) (string, error) {

	if len(username) < 3 {
		return "", models.ErrInvalidUsername
	}
	if len(password) < 8 {
		return "", models.ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username:     username,
		PasswordHash: string(hash),
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return "", err
	}

	return s.generateToken(user)
}

func (s *authService) generateToken(user *models.User) (string, error) {
	claims := Claims{
		UserID:   user.ID.String(),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
