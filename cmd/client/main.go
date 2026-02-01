package main

import (
	"log"
	"os"
	"path/filepath"

	config "github.com/Popolzen/go_final_project/configs"
	"github.com/Popolzen/go_final_project/internal/client/cli"
)

func main() {
	tokenPath := filepath.Join(configDir(), "token")
	cfg, err := config.New()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	cmd := cli.NewRootCmd(tokenPath, []byte(cfg.EncryptKey))
	cli.Execute(cmd)
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "gophkeeper")
}
