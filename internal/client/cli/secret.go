package cli

import (
	"encoding/json"
	"fmt"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавить секрет",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newAddLoginCmd())

	return cmd
}

func newAddLoginCmd() *cobra.Command {
	var (
		login    string
		password string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "login <n>",
		Short: "Add login/pass",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			data, err := json.Marshal(models.LoginData{
				Login:    login,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			result, err := secretService.Create(cmd.Context(), models.SecretTypeLogin, name, metadata, data)
			if err != nil {
				return err
			}

			fmt.Printf("создан: id=%s name=%s\n", result.ID, result.Name)
			return nil
		},
	}

	cmd.Flags().StringVar(&login, "login", "", "логин (обязательно)")
	cmd.Flags().StringVar(&password, "password", "", "пароль (обязательно)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")
	cmd.MarkFlagRequired("login")
	cmd.MarkFlagRequired("password")

	return cmd
}
