package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func (c *CLI) newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Список секретов",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			secrets, err := c.secretService.List(cmd.Context())
			if err != nil {
				return err
			}

			if len(secrets) == 0 {
				fmt.Println("нет секретов")
				return nil
			}

			fmt.Printf("%-36s %-10s %-20s %-20s\n", "ID", "TYPE", "NAME", "UPDATED")
			for _, s := range secrets {
				fmt.Printf("%-36s %-10s %-20s %-20s\n", s.ID, s.Type, s.Name, s.UpdatedAt)
			}

			return nil
		},
	}

	return cmd
}
