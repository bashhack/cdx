package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.OutputFormat != "auto" {
		t.Errorf("OutputFormat = %q, want %q", cfg.OutputFormat, "auto")
	}
	if cfg.ContextLines != 2 {
		t.Errorf("ContextLines = %d, want %d", cfg.ContextLines, 2)
	}
	if cfg.Color != nil {
		t.Errorf("Color = %v, want nil (auto-detect)", cfg.Color)
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	// Load should succeed even with no config file (graceful degradation)
	// Change to a temp directory with no config
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Should return defaults when no config file exists
	if cfg.OutputFormat != "auto" {
		t.Errorf("OutputFormat = %q, want %q", cfg.OutputFormat, "auto")
	}
}

func TestLoad_WithConfigFile(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	// Create a config file
	configContent := `output_format: json
context_lines: 5
`
	configPath := filepath.Join(tmp, ".cdx.yaml")
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatal(err)
	}

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.OutputFormat != "json" {
		t.Errorf("OutputFormat = %q, want %q", cfg.OutputFormat, "json")
	}
	if cfg.ContextLines != 5 {
		t.Errorf("ContextLines = %d, want %d", cfg.ContextLines, 5)
	}
}

func TestLoad_EnvOverride(t *testing.T) {
	tmp := t.TempDir()
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(origDir); err != nil {
			t.Errorf("failed to restore working directory: %v", err)
		}
	})

	if err := os.Chdir(tmp); err != nil {
		t.Fatal(err)
	}

	if err := os.Setenv("CDX_OUTPUT_FORMAT", "plain"); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := os.Unsetenv("CDX_OUTPUT_FORMAT"); err != nil {
			t.Errorf("failed to unset CDX_OUTPUT_FORMAT: %v", err)
		}
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.OutputFormat != "plain" {
		t.Errorf("OutputFormat = %q, want %q (from env)", cfg.OutputFormat, "plain")
	}
}

func TestConfigDir(t *testing.T) {
	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir() error = %v", err)
	}

	// Should end with .config/cdx
	if filepath.Base(dir) != "cdx" {
		t.Errorf("ConfigDir() = %q, want path ending in 'cdx'", dir)
	}
	if filepath.Base(filepath.Dir(dir)) != ".config" {
		t.Errorf("ConfigDir() = %q, want path containing '.config'", dir)
	}
}
