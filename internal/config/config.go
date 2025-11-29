// Package config handles configuration loading for cdx.
package config

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	// Output format: "auto", "human", "json", "plain"
	OutputFormat string `mapstructure:"output_format"`
	// Default context lines for search results
	ContextLines int `mapstructure:"context_lines"`
	// Whether to use color output (auto-detected if not set)
	Color *bool `mapstructure:"color"`
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() *Config {
	return &Config{
		OutputFormat: "auto",
		ContextLines: 2,
		Color:        nil, // auto-detect
	}
}

// Load reads configuration from files and environment.
// It looks for config in:
// 1. .cdx.yaml in current directory
// 2. .cdx.yaml in home directory
// 3. Environment variables prefixed with CDX_
func Load() (*Config, error) {
	cfg := DefaultConfig()

	v := viper.New()
	v.SetConfigName(".cdx")
	v.SetConfigType("yaml")

	// Look in current directory first
	v.AddConfigPath(".")

	// Then home directory
	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(home)
	}

	// Set defaults so Viper knows about the keys
	v.SetDefault("output_format", cfg.OutputFormat)
	v.SetDefault("context_lines", cfg.ContextLines)

	// Environment variables (CDX_OUTPUT_FORMAT, CDX_CONTEXT_LINES, etc.)
	v.SetEnvPrefix("CDX")
	v.AutomaticEnv()

	// Read config file (ignore if not found)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	if err := v.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ConfigDir returns the path to the user's cdx config directory.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "cdx"), nil
}
