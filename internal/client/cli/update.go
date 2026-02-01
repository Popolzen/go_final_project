package cli

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Popolzen/go_final_project/internal/models"
	"github.com/spf13/cobra"
)

func newUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Обновить секрет",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newUpdateLoginCmd())
	cmd.AddCommand(newUpdateTextCmd())
	cmd.AddCommand(newUpdateBinaryCmd())
	cmd.AddCommand(newUpdateCardCmd())

	return cmd
}

func newUpdateLoginCmd() *cobra.Command {
	var (
		login    string
		password string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "login <type> <n>",
		Short: "Update login/password",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType(args[0])
			name := args[1]

			data, err := json.Marshal(models.LoginData{
				Login:    login,
				Password: password,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
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

func newUpdateTextCmd() *cobra.Command {
	var metadata string

	cmd := &cobra.Command{
		Use:   "text <type> <n> <content>",
		Short: "Update text note",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType(args[0])
			name := args[1]
			content := args[2]

			data, err := json.Marshal(models.TextData{
				Content: content,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
				return err
			}

			fmt.Println("обновлено")
			return nil
		},
	}

	cmd.Flags().StringVar(&metadata, "metadata", "", "метаинформация")

	return cmd
}

func newUpdateBinaryCmd() *cobra.Command {
	var (
		file     string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "binary <type> <n>",
		Short: "Update binary file",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType(args[0])
			name := args[1]

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

			if err := secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
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

func newUpdateCardCmd() *cobra.Command {
	var (
		number   string
		holder   string
		cvv      string
		expiry   string
		metadata string
	)

	cmd := &cobra.Command{
		Use:   "card <type> <n>",
		Short: "Update bank card",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			secretType := models.SecretType(args[0])
			name := args[1]

			data, err := json.Marshal(models.CardData{
				Number: number,
				Holder: holder,
				CVV:    cvv,
				Expiry: expiry,
			})
			if err != nil {
				return fmt.Errorf("marshal data: %w", err)
			}

			if err := secretService.Update(cmd.Context(), name, secretType, metadata, data); err != nil {
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
