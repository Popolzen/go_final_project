// Package cli реализует командный интерфейс клиента GophKeeper.
package cli

import (
	"fmt"
	"net/http"
	"os"

	service "github.com/Popolzen/go_final_project/internal/client/service"
	"github.com/spf13/cobra"
)

type CLI struct {
	serverAddr string
	tokenPath  string
	encKey     []byte

	httpClient *http.Client

	authService   service.AuthService
	secretService service.SecretService
}

func NewCLI(tokenPath string, encKey []byte, httpClient *http.Client) *CLI {
	return &CLI{
		tokenPath:  tokenPath,
		encKey:     encKey,
		httpClient: httpClient,
	}
}

func (c *CLI) NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "Password manager",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},

		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			c.authService = service.NewAuthService(
				c.serverAddr,
				c.tokenPath,
				c.httpClient,
			)

			c.secretService = service.NewSecretService(
				c.serverAddr,
				c.tokenPath,
				c.encKey,
				c.httpClient,
			)

			return nil
		},
	}

	cmd.PersistentFlags().
		StringVar(&c.serverAddr, "server", "http://localhost:8080", "адрес сервера")

	cmd.AddCommand(c.newRegisterCmd())
	cmd.AddCommand(c.newLoginCmd())
	cmd.AddCommand(c.newAddCmd())
	cmd.AddCommand(c.newListCmd())
	cmd.AddCommand(c.newGetCmd())
	cmd.AddCommand(c.newUpdateCmd())
	cmd.AddCommand(c.newDeleteCmd())

	return cmd
}

func Execute(cmd *cobra.Command) {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
