package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Popolzen/go_final_project/internal/models"
)

func (c *CLI) createSecret(
	ctx context.Context,
	secretType models.SecretType,
	name string,
	metadata string,
	data any,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	result, err := c.secretService.Create(ctx, secretType, name, metadata, payload)
	if err != nil {
		return err
	}

	fmt.Printf("создан: id=%s name=%s\n", result.ID, result.Name)
	return nil
}

func (c *CLI) updateSecret(
	ctx context.Context,
	secretType models.SecretType,
	name string,
	metadata string,
	data any,
) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}

	if err := c.secretService.Update(ctx, name, secretType, metadata, payload); err != nil {
		return err
	}

	fmt.Println("обновлено")
	return nil
}
