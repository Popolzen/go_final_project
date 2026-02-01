package cli

import (
	"fmt"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/spf13/cobra"
)

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <type> <n>",
		Short: "Удалить секрет",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType(args[0])
			name := args[1]

			if !secretType.Validate() {
				return fmt.Errorf("invalid type: %s (login/text/binary/card)", args[0])
			}

			if err := secretService.Delete(cmd.Context(), name, secretType); err != nil {
				return err
			}

			fmt.Println("удалено")
			return nil
		},
	}

	return cmd
}
