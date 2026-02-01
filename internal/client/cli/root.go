// Package cli реализует командный интерфейс клиента GophKeeper.
package cli

import (
	"fmt"
	"os"

	service "github.com/Popolzen/go_final_project/internal/client/service"
	"github.com/spf13/cobra"
)

var (
	serverAddr    string
	tokenPath     string
	authService   service.AuthService
	secretService service.SecretService
)

func NewRootCmd(tp string) *cobra.Command {
	tokenPath = tp

	cmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "Password manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			authService = service.NewAuthService(serverAddr, tokenPath)
			secretService = service.NewSecretService(serverAddr, tokenPath)
			return nil
		},
	}

	cmd.PersistentFlags().StringVar(&serverAddr, "server", "http://localhost:8080", "адрес сервера")

	cmd.AddCommand(newRegisterCmd())
	cmd.AddCommand(newLoginCmd())
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())

	return cmd
}

func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
