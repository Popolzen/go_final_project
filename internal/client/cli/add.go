package cli

import (
	"encoding/base64"
	"os"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/spf13/cobra"
)

func (c *CLI) newAddCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Добавить секрет",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(c.newAddLoginCmd())
	cmd.AddCommand(c.newAddTextCmd())
	cmd.AddCommand(c.newAddBinaryCmd())
	cmd.AddCommand(c.newAddCardCmd())

	return cmd
}

func (c *CLI) newAddLoginCmd() *cobra.Command {
	var (
		login    string
		password string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "login <name>",
		Short: "Add login/password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.createSecret(
				cmd.Context(),
				models.SecretTypeLogin,
				args[0],
				metadata,
				models.LoginData{
					Login:    login,
					Password: password,
				},
			)
		},
	}

	cmd.Flags().StringVar(&login, "login", "", "логин (обязательно)")
	cmd.Flags().StringVar(&password, "password", "", "пароль (обязательно)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")
	cmd.MarkFlagRequired("login")
	cmd.MarkFlagRequired("password")

	return cmd
}

func (c *CLI) newAddTextCmd() *cobra.Command {
	var metadata string

	cmd := &cobra.Command{
		Use:   "text <name> <content>",
		Short: "Add text note",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.createSecret(
				cmd.Context(),
				models.SecretTypeText,
				args[0],
				metadata,
				models.TextData{
					Content: args[1],
				},
			)
		},
	}

	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")

	return cmd
}

func (c *CLI) newAddBinaryCmd() *cobra.Command {
	var (
		file     string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "binary <name>",
		Short: "Add binary file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileData, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			return c.createSecret(
				cmd.Context(),
				models.SecretTypeBinary,
				args[0],
				metadata,
				models.BinaryData{
					Filename: file,
					Content:  base64.StdEncoding.EncodeToString(fileData),
				},
			)
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "путь к файлу (обязательно)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")
	cmd.MarkFlagRequired("file")

	return cmd
}

func (c *CLI) newAddCardCmd() *cobra.Command {
	var (
		number   string
		holder   string
		cvv      string
		expiry   string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "card <name>",
		Short: "Add bank card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return c.createSecret(
				cmd.Context(),
				models.SecretTypeCard,
				args[0],
				metadata,
				models.CardData{
					Number: number,
					Holder: holder,
					CVV:    cvv,
					Expiry: expiry,
				},
			)
		},
	}

	cmd.Flags().StringVar(&number, "number", "", "номер карты (обязательно)")
	cmd.Flags().StringVar(&holder, "holder", "", "владелец (обязательно)")
	cmd.Flags().StringVar(&cvv, "cvv", "", "CVV (обязательно)")
	cmd.Flags().StringVar(&expiry, "expiry", "", "срок действия MM/YY (обязательно)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")
	cmd.MarkFlagRequired("number")
	cmd.MarkFlagRequired("holder")
	cmd.MarkFlagRequired("cvv")
	cmd.MarkFlagRequired("expiry")

	return cmd
}
