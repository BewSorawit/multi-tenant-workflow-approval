package platform

import (
	"testing"
	"time"
)

func TestLoadConfig_DefaultValues(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("APP_ENV", "")
	t.Setenv("APP_PORT", "")
	t.Setenv("SHUTDOWN_TIMEOUT", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.AppEnv != "local" {
		t.Errorf("want local got %s", cfg.AppEnv)
	}

	if cfg.AppPort != "8080" {
		t.Errorf("want 8080 got %s", cfg.AppPort)
	}

	if cfg.ShutdownTimeout != 10*time.Second {
		t.Errorf("unexpected timeout")
	}
}

func TestLoadConfig_MissingDatabaseURL(t *testing.T) {
	t.Setenv("DATABASE_URL", "")

	_, err := LoadConfig()

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestLoadConfig_ShutdownTimeout(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("SHUTDOWN_TIMEOUT", "30s")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}

	if cfg.ShutdownTimeout != 30*time.Second {
		t.Fatal("timeout mismatch")
	}
}

func TestLoadConfig_InvalidTimeout(t *testing.T) {
	t.Setenv("DATABASE_URL", "postgres://localhost/test")
	t.Setenv("SHUTDOWN_TIMEOUT", "abc")

	_, err := LoadConfig()

	if err == nil {
		t.Fatal("expected error")
	}
}
