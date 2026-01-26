package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/caarlos0/env/v9"
)

const (
	DefaultServerAddr = ":8080"
	DefaultDBDSN      = "host=localhost port=5432 user=postgres password=123456 dbname=gophkeeper sslmode=disable"
)

// Config holds application configuration
type Config struct {
	ServerAddr  string `json:"server_address" env:"SERVER_ADDRESS"`
	EnableHTTPS bool   `json:"enable_https" env:"ENABLE_HTTPS"`
	CertFile    string `json:"cert_file" env:"CERT_FILE"`
	KeyFile     string `json:"key_file" env:"KEY_FILE"`

	JWTSecret  string `env:"JWT_SECRET"`
	EncryptKey string `env:"ENCRYPT_KEY"`

	DatabaseDSN string `json:"database_dsn" env:"DATABASE_DSN"`
}

// defaults < file < env < cli
func New() (*Config, error) {
	cfg := defaultConfig()

	path, err := parseFlags()
	if err != nil {
		return nil, err
	}

	if err := cfg.loadFromFile(path); err != nil {
		return nil, err
	}

	if err := cfg.loadFromEnv(); err != nil {
		return nil, err
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		ServerAddr:  DefaultServerAddr,
		DatabaseDSN: DefaultDBDSN,
	}
}

func parseFlags() (string, error) {
	var configPath string

	flag.StringVar(&configPath, "c", "", "config file path")
	flag.StringVar(&configPath, "config", "", "config file path")

	flag.Parse()

	if configPath == "" {
		configPath = os.Getenv("CONFIG")
	}

	return configPath, nil
}

func (c *Config) loadFromFile(path string) error {
	if path == "" {
		return nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config file: %w", err)
	}

	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("parse config file: %w", err)
	}

	return nil
}

func (c *Config) loadFromEnv() error {
	if err := env.Parse(c); err != nil {
		return fmt.Errorf("parse env: %w", err)
	}
	return nil
}

func (c *Config) validate() error {
	if c.JWTSecret == "" {
		return errors.New("JWT_SECRET is required")
	}

	if c.EncryptKey == "" {
		return errors.New("ENCRYPT_KEY is required")
	}

	if c.EnableHTTPS {
		if c.CertFile == "" || c.KeyFile == "" {
			return errors.New("CERT_FILE and KEY_FILE are required when HTTPS is enabled")
		}
	}

	return nil
}
