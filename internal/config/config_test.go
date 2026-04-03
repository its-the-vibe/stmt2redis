package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/its-the-vibe/stmt2redis/internal/config"
)

const testYAML = `
redis:
  host: redis.example.com
  port: 6380
  db: 1

lists:
  starling: transactions:starling
  amex: transactions:amex
  monzo: transactions:monzo
  monzo_flex: transactions:monzo-flex
`

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "config*.yaml")
	if err != nil {
		t.Fatalf("creating temp config: %v", err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatalf("writing temp config: %v", err)
	}
	f.Close()
	return f.Name()
}

func TestLoad(t *testing.T) {
	path := writeTempConfig(t, testYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Redis.Host != "redis.example.com" {
		t.Errorf("expected host redis.example.com, got %s", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6380 {
		t.Errorf("expected port 6380, got %d", cfg.Redis.Port)
	}
	if cfg.Redis.DB != 1 {
		t.Errorf("expected db 1, got %d", cfg.Redis.DB)
	}
}

func TestLoadDefaults(t *testing.T) {
	path := writeTempConfig(t, "lists:\n  starling: tx:starling\n")
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Redis.Host != "localhost" {
		t.Errorf("expected default host localhost, got %s", cfg.Redis.Host)
	}
	if cfg.Redis.Port != 6379 {
		t.Errorf("expected default port 6379, got %d", cfg.Redis.Port)
	}
}

func TestLoadFileNotFound(t *testing.T) {
	_, err := config.Load(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestListKey(t *testing.T) {
	path := writeTempConfig(t, testYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	tests := []struct {
		csvType string
		want    string
	}{
		{"starling", "transactions:starling"},
		{"amex", "transactions:amex"},
		{"monzo", "transactions:monzo"},
		{"monzo-flex", "transactions:monzo-flex"},
	}

	for _, tt := range tests {
		key, err := cfg.ListKey(tt.csvType)
		if err != nil {
			t.Errorf("ListKey(%q): unexpected error: %v", tt.csvType, err)
			continue
		}
		if key != tt.want {
			t.Errorf("ListKey(%q): got %q, want %q", tt.csvType, key, tt.want)
		}
	}
}

func TestListKeyUnsupported(t *testing.T) {
	path := writeTempConfig(t, testYAML)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, err = cfg.ListKey("unknown")
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}
