package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveCreatesPrivateConfigPath(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	cfg := Config{
		BaseURL: "https://api.example.com/v1/",
		Token:   "secret-token",
	}
	if err := Save(cfg); err != nil {
		t.Fatalf("save config: %v", err)
	}

	p, err := configPath()
	if err != nil {
		t.Fatalf("configPath: %v", err)
	}

	dirInfo, err := os.Stat(filepath.Dir(p))
	if err != nil {
		t.Fatalf("stat config dir: %v", err)
	}
	if dirInfo.Mode().Perm()&0o077 != 0 {
		t.Fatalf("config dir should not be accessible by group/others, got mode %o", dirInfo.Mode().Perm())
	}

	fileInfo, err := os.Stat(p)
	if err != nil {
		t.Fatalf("stat config file: %v", err)
	}
	if fileInfo.Mode().Perm()&0o077 != 0 {
		t.Fatalf("config file should not be accessible by group/others, got mode %o", fileInfo.Mode().Perm())
	}
}

func TestLoadUsesEnvOverrides(t *testing.T) {
	t.Setenv("HOME", t.TempDir())

	if err := Save(Config{BaseURL: "https://api.disk.example/v1", Token: "disk-token"}); err != nil {
		t.Fatalf("save config: %v", err)
	}

	t.Setenv("LINKMNGR_BASE_URL", "https://api.env.example/v1/")
	t.Setenv("LINKMNGR_TOKEN", "env-token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.BaseURL != "https://api.env.example/v1" {
		t.Fatalf("expected env base URL override, got %q", cfg.BaseURL)
	}
	if cfg.Token != "env-token" {
		t.Fatalf("expected env token override, got %q", cfg.Token)
	}
}
