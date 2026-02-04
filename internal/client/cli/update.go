package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/spf13/cobra"
)

func (c *CLI) newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Обновить секрет",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(c.newUpdateLoginCmd())
	cmd.AddCommand(c.newUpdateTextCmd())
	cmd.AddCommand(c.newUpdateBinaryCmd())
	cmd.AddCommand(c.newUpdateCardCmd())

	return cmd
}

func (c *CLI) newUpdateLoginCmd() *cobra.Command {
	var (
		login    string
		password string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "login <n>",
		Short: "Update login/password",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType("login")
			name := args[0]

			data, err := json.Marshal(models.LoginData{
				Login:    login,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := c.secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
				return err
			}

			fmt.Println("обновлено")
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

func (c *CLI) newUpdateTextCmd() *cobra.Command {
	var metadata string

	cmd := &cobra.Command{
		Use:   "text <n> <content>",
		Short: "Update text note",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType("text")
			name := args[0]
			content := args[1]

			data, err := json.Marshal(models.TextData{
				Content: content,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := c.secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
				return err
			}

			fmt.Println("обновлено")
			return nil
		},
	}

	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")

	return cmd
}

func (c *CLI) newUpdateBinaryCmd() *cobra.Command {
	var (
		file     string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "binary <n>",
		Short: "Update binary file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType("binary")
			name := args[0]

			fileData, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("read file: %w", err)
			}

			data, err := json.Marshal(models.BinaryData{
				Filename: file,
				Content:  base64.StdEncoding.EncodeToString(fileData),
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := c.secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
				return err
			}

			fmt.Println("обновлено")
			return nil
		},
	}

	cmd.Flags().StringVar(&file, "file", "", "путь к файлу (обязательно)")
	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")
	cmd.MarkFlagRequired("file")

	return cmd
}

func (c *CLI) newUpdateCardCmd() *cobra.Command {
	var (
		number   string
		holder   string
		cvv      string
		expiry   string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "card <n>",
		Short: "Update bank card",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType("card")
			name := args[0]

			data, err := json.Marshal(models.CardData{
				Number: number,
				Holder: holder,
				CVV:    cvv,
				Expiry: expiry,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := c.secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
				return err
			}

			fmt.Println("обновлено")
			return nil
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
