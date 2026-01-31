package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRegisterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "register <username> <password>",
		Short: "Register new user",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			username := args[0]
			password := args[1]

			if err := authService.Register(cmd.Context(), username, password); err != nil {
				return fmt.Errorf("register: %w", err)
			}

			fmt.Println("register succses")
			return nil
		},
	}

	return cmd
}
