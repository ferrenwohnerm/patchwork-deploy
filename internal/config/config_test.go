package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/patchwork-deploy/internal/config"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "patchwork.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write temp config: %v", err)
	}
	return path
}

func TestLoad_Valid(t *testing.T) {
	content := `
version: "1"
project: my-app
environments:
  - name: staging
    base_url: https://staging.example.com
    variables:
      LOG_LEVEL: debug
  - name: production
    base_url: https://example.com
    variables:
      LOG_LEVEL: info
`
	path := writeTemp(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Project != "my-app" {
		t.Errorf("expected project 'my-app', got %q", cfg.Project)
	}
	if len(cfg.Environments) != 2 {
		t.Errorf("expected 2 environments, got %d", len(cfg.Environments))
	}
}

func TestLoad_MissingFile(t *testing.T) {
	_, err := config.Load("/nonexistent/patchwork.yaml")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestLoad_DuplicateEnvironment(t *testing.T) {
	content := `
project: dup-test
environments:
  - name: staging
    base_url: https://staging.example.com
  - name: staging
    base_url: https://staging2.example.com
`
	path := writeTemp(t, content)
	_, err := config.Load(path)
	if err == nil {
		t.Fatal("expected error for duplicate environment, got nil")
	}
}

func TestGetEnvironment(t *testing.T) {
	content := `
project: env-test
environments:
  - name: production
    base_url: https://prod.example.com
`
	path := writeTemp(t, content)
	cfg, err := config.Load(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	env, err := cfg.GetEnvironment("production")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env.BaseURL != "https://prod.example.com" {
		t.Errorf("unexpected base_url: %q", env.BaseURL)
	}
	_, err = cfg.GetEnvironment("staging")
	if err == nil {
		t.Fatal("expected error for missing environment, got nil")
	}
}
