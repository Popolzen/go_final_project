package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

type authService struct {
	serverAddr string
	tokenPath  string
	httpClient *http.Client
}

func NewAuthService(
	serverAddr string,
	tokenPath string,
	httpClient *http.Client,
) AuthService {
	return &authService{
		serverAddr: serverAddr,
		tokenPath:  tokenPath,
		httpClient: httpClient,
	}
}

type authRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Token string `json:"token"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func (s *authService) Register(ctx context.Context, username, password string) error {
	token, err := s.callAuth(ctx, "/api/v1/auth/register", username, password)
	if err != nil {
		return err
	}

	return s.saveToken(token)
}

func (s *authService) Login(ctx context.Context, username, password string) error {
	token, err := s.callAuth(ctx, "/api/v1/auth/login", username, password)
	if err != nil {
		return err
	}

	return s.saveToken(token)
}

// callAuth выполняет POST-запрос к указанному эндпойнту и возвращает токен
func (s *authService) callAuth(ctx context.Context, path, username, password string) (string, error) {
	body, err := json.Marshal(authRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.serverAddr+path, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request to server: %w", err)
	}
	defer resp.Body.Close()

	// Успех
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return "", fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		return "", fmt.Errorf("server error: %s", errResp.Error)
	}

	var authResp authResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	return authResp.Token, nil
}

// saveToken сохраняет токен в файл
func (s *authService) saveToken(token string) error {
	dir := filepath.Dir(s.tokenPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("create token dir: %w", err)
	}

	if err := os.WriteFile(s.tokenPath, []byte(token), 0600); err != nil {
		return fmt.Errorf("save token: %w", err)
	}

	return nil
}
