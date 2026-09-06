package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

func TestHTTPServerDurations(t *testing.T) {
	var cfg Config
	if err := cleanenv.ReadConfig("../../config/local.yaml", &cfg); err != nil {
		t.Fatal(err)
	}
	if cfg.HTTPServer.Timeout != 4*time.Second || cfg.HTTPServer.IdleTimeout != time.Minute {
		t.Fatalf("unexpected durations: %+v", cfg.HTTPServer)
	}

	path := filepath.Join(t.TempDir(), "invalid.yaml")
	if err := os.WriteFile(path, []byte("http_server:\n  timeout: invalid\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := cleanenv.ReadConfig(path, &cfg); err == nil {
		t.Fatal("expected invalid duration to be rejected")
	}
}
