package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBaseURL = "https://api.linkmngr.com/v1"
	dirName        = ".linkmngr"
	fileName       = "config.json"
)

type Config struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve user home dir: %w", err)
	}
	return filepath.Join(home, dirName, fileName), nil
}

func Load() (Config, error) {
	cfg := Config{BaseURL: defaultBaseURL}

	if envURL := strings.TrimRight(os.Getenv("LINKMNGR_BASE_URL"), "/"); envURL != "" {
		cfg.BaseURL = envURL
	}
	if envToken := os.Getenv("LINKMNGR_TOKEN"); envToken != "" {
		cfg.Token = envToken
	}

	p, err := configPath()
	if err != nil {
		return cfg, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return cfg, nil
		}
		return cfg, fmt.Errorf("read config file %q: %w", p, err)
	}

	var disk Config
	if err := json.Unmarshal(data, &disk); err != nil {
		return cfg, fmt.Errorf("parse config file %q: %w", p, err)
	}

	if cfg.BaseURL == defaultBaseURL && disk.BaseURL != "" {
		cfg.BaseURL = strings.TrimRight(disk.BaseURL, "/")
	}
	if cfg.Token == "" {
		cfg.Token = disk.Token
	}
	return cfg, nil
}

func Save(cfg Config) error {
	if cfg.BaseURL == "" {
		cfg.BaseURL = defaultBaseURL
	}
	cfg.BaseURL = strings.TrimRight(cfg.BaseURL, "/")

	p, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("serialize config: %w", err)
	}
	if err := os.WriteFile(p, data, 0o600); err != nil {
		return fmt.Errorf("write config file %q: %w", p, err)
	}
	return nil
}
