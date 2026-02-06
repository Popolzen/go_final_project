package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *CLI) newLoginCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login <username> <password>",
		Short: "Enter existing account",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			password := args[1]

			if err := c.authService.Login(cmd.Context(), username, password); err != nil {
				return fmt.Errorf("login: %w", err)
			}

			fmt.Println("login success")
			return nil
		},
	}

	return cmd
}
