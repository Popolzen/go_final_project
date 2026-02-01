package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/spf13/cobra"
)

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <type> <n>",
		Short: "Получить секрет",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType(args[0])
			name := args[1]

			if !secretType.Validate() {
				return fmt.Errorf("invalid type: %s (login/text/binary/card)", args[0])
			}

			secret, err := secretService.GetByName(cmd.Context(), name, secretType)
			if err != nil {
				return err
			}

			fmt.Printf("ID: %s\n", secret.ID)
			fmt.Printf("Name: %s\n", secret.Name)
			fmt.Printf("Type: %s\n", secret.Type)
			fmt.Printf("Version: %d\n", secret.Version)
			fmt.Printf("Updated: %s\n", secret.UpdatedAt)
			if secret.Metadata != "" {
				fmt.Printf("Metadata: %s\n", secret.Metadata)
			}
			fmt.Println()

			switch secret.Type {
			case models.SecretTypeLogin:
				var login models.LoginData
				if err := json.Unmarshal(secret.Data, &login); err != nil {
					return fmt.Errorf("decode login: %w", err)
				}
				fmt.Printf("Login: %s\n", login.Login)
				fmt.Printf("Password: %s\n", login.Password)

			case models.SecretTypeText:
				var text models.TextData
				if err := json.Unmarshal(secret.Data, &text); err != nil {
					return fmt.Errorf("decode text: %w", err)
				}
				fmt.Printf("Content:\n%s\n", text.Content)

			case models.SecretTypeBinary:
				var binary models.BinaryData
				if err := json.Unmarshal(secret.Data, &binary); err != nil {
					return fmt.Errorf("decode binary: %w", err)
				}
				fmt.Printf("Filename: %s\n", binary.Filename)
				decoded, err := base64.StdEncoding.DecodeString(binary.Content)
				if err != nil {
					return fmt.Errorf("decode base64: %w", err)
				}
				fmt.Printf("Size: %d bytes\n", len(decoded))

			case models.SecretTypeCard:
				var card models.CardData
				if err := json.Unmarshal(secret.Data, &card); err != nil {
					return fmt.Errorf("decode card: %w", err)
				}
				fmt.Printf("Number: %s\n", card.Number)
				fmt.Printf("Holder: %s\n", card.Holder)
				fmt.Printf("CVV: %s\n", card.CVV)
				fmt.Printf("Expiry: %s\n", card.Expiry)
			}

			return nil
		},
	}

	return cmd
}
