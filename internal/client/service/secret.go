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
	encKey     []byte
}

func NewSecretService(serverAddr, tokenPath string, encKey []byte) SecretService {
	return &secretService{
		serverAddr: serverAddr,
		tokenPath:  tokenPath,
		encKey:     encKey,
	}
}

type createRequest struct {
	Type     models.SecretType `json:"type"`
	Name     string            `json:"name"`
	Data     []byte            `json:"data"`
	Metadata string            `json:"metadata,omitempty"`
}

type updateRequest struct {
	Data     []byte `json:"data"`
	Metadata string `json:"metadata,omitempty"`
}

func (s *secretService) Create(ctx context.Context, secretType models.SecretType, name, metadata string, data []byte) (*SecretInfo, error) {
	token, err := s.readToken()
	if err != nil {
		return nil, err
	}

	encrypted, err := s.encrypt(data)
	if err != nil {
		return nil, fmt.Errorf("encrypt data request: %w", err)
	}

	body, err := json.Marshal(createRequest{
		Type:     secretType,
		Name:     name,
		Data:     encrypted,
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

func (s *secretService) List(ctx context.Context) ([]*SecretInfo, error) {
	token, err := s.readToken()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.serverAddr+"/api/v1/secrets/", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("server error: %s", errResp.Error)
	}

	var result []*SecretInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return result, nil
}

func (s *secretService) GetByName(ctx context.Context, name string, secretType models.SecretType) (*SecretDetail, error) {
	token, err := s.readToken()
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/v1/secrets/by-name/%s?type=%s", s.serverAddr, name, secretType)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("server error: %s", errResp.Error)
	}

	var result SecretDetail
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	decrypted, err := s.decrypt(result.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}
	result.Data = decrypted
	return &result, nil
}

func (s *secretService) Update(ctx context.Context, name string, secretType models.SecretType, metadata string, data []byte) error {
	token, err := s.readToken()
	if err != nil {
		return err
	}

	secret, err := s.GetByName(ctx, name, secretType)
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	encrypted, err := s.encrypt(data)
	if err != nil {
		return fmt.Errorf("encrypt data request: %w", err)
	}
	body, err := json.Marshal(updateRequest{
		Data:     encrypted,
		Metadata: metadata,
	})
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/secrets/%s", s.serverAddr, secret.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("server error: %s", errResp.Error)
	}

	return nil
}

func (s *secretService) Delete(ctx context.Context, name string, secretType models.SecretType) error {
	token, err := s.readToken()
	if err != nil {
		return err
	}

	secret, err := s.GetByName(ctx, name, secretType)
	if err != nil {
		return fmt.Errorf("get secret: %w", err)
	}

	url := fmt.Sprintf("%s/api/v1/secrets/%s", s.serverAddr, secret.ID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("request to server: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		var errResp errorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return fmt.Errorf("server error: status %d", resp.StatusCode)
		}
		return fmt.Errorf("server error: %s", errResp.Error)
	}

	return nil
}

// readToken читает токен из файла.
func (s *secretService) readToken() (string, error) {
	data, err := os.ReadFile(s.tokenPath)
	if err != nil {
		return "", fmt.Errorf("не авторизован, выполните login")
	}
	return string(data), nil
}
