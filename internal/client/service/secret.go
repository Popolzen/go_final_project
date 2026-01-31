package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Popolzen/go_final_project/internal/models"
)

type secretService struct {
	serverAddr string
	tokenPath  string
}

func NewSecretService(serverAddr, tokenPath string) SecretService {
	return &secretService{
		serverAddr: serverAddr,
		tokenPath:  tokenPath,
	}
}

type createRequest struct {
	Type     models.SecretType `json:"type"`
	Name     string            `json:"name"`
	Data     []byte            `json:"data"`
	Metadata string            `json:"metadata,omitempty"`
}

func (s *secretService) Create(ctx context.Context, secretType models.SecretType, name, metadata string, data []byte) (*SecretInfo, error) {
	token, err := s.readToken()
	if err != nil {
		return nil, err
	}

	body, err := json.Marshal(createRequest{
		Type:     secretType,
		Name:     name,
		Data:     data,
		Metadata: metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.serverAddr+"/api/v1/secrets/", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("server error: %s", errResp.Error)
	}

	var result SecretInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &result, nil
}

// readToken читает токен из файла.
func (s *secretService) readToken() (string, error) {
	data, err := os.ReadFile(s.tokenPath)
	if err != nil {
		return "", fmt.Errorf("не авторизован, выполните login")
	}
	return string(data), nil
}
