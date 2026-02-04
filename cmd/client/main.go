package main

import (
	"log"
	"os"
	"path/filepath"

	config "github.com/Popolzen/go_final_project/configs"
	"github.com/Popolzen/go_final_project/internal/client/cli"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	tokenPath := filepath.Join(configDir(), "token")

	cliApp := cli.NewCLI(tokenPath, []byte(cfg.EncryptKey))
	cmd := cliApp.NewRootCmd()

	cli.Execute(cmd)
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "."
	}
	return filepath.Join(home, ".config", "gophkeeper")
}
