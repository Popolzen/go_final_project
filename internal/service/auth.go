package service

import (
	"context"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/Popolzen/go_final_project/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo      storage.Repository
	jwtSecret string
}

// func NewAuthService(repo storage.Repository, jwtSecret string) AuthService {
// 	return authService{
// 		repo:      repo,
// 		jwtSecret: string,
// 	}
// }

func (s *authService) Login(ctx context.Context, username, password string) (string, error) {
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

	// TODO: generate token
	return "", nil
}

func (s *authService) Register(ctx context.Context, username, password string) (string, error) {

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

	// TODO: generate token
	return "", nil
}
