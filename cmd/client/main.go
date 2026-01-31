package main

import (
	"os"
	"path/filepath"

	"github.com/Popolzen/go_final_project/internal/client/cli"
)

func main() {
	tokenPath := filepath.Join(configDir(), "token")

	cmd := cli.NewRootCmd(tokenPath)
	cli.Execute(cmd)
}

func configDir() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "."
	}
	return filepath.Join(home, ".config", "gophkeeper")
}
