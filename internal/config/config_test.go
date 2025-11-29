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
	// Change to a temp directory with no config and isolate HOME
	// to prevent picking up ~/.cdx.yaml from the developer's machine
	tmp := t.TempDir()
	t.Setenv("HOME", tmp)
	origDir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if chErr := os.Chdir(origDir); chErr != nil {
			t.Errorf("failed to restore working directory: %v", chErr)
		}
	})

	err = os.Chdir(tmp)
	if err != nil {
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
		if chErr := os.Chdir(origDir); chErr != nil {
			t.Errorf("failed to restore working directory: %v", chErr)
		}
	})

	// Create a config file
	configContent := `output_format: json
context_lines: 5
`
	configPath := filepath.Join(tmp, ".cdx.yaml")
	err = os.WriteFile(configPath, []byte(configContent), 0o600)
	if err != nil {
		t.Fatal(err)
	}

	err = os.Chdir(tmp)
	if err != nil {
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
		if chErr := os.Chdir(origDir); chErr != nil {
			t.Errorf("failed to restore working directory: %v", chErr)
		}
	})

	err = os.Chdir(tmp)
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("CDX_OUTPUT_FORMAT", "plain")

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

	// Should end with cdx (OS-specific parent varies)
	if filepath.Base(dir) != "cdx" {
		t.Errorf("ConfigDir() = %q, want path ending in 'cdx'", dir)
	}

	// Should be an absolute path
	if !filepath.IsAbs(dir) {
		t.Errorf("ConfigDir() = %q, want absolute path", dir)
	}
}
