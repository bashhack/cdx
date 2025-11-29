package cli

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	tests := []struct {
		name      string
		version   string
		commit    string
		buildDate string
		wantOut   []string // substrings that should appear in output
	}{
		{
			name:      "dev version",
			version:   "dev",
			commit:    "unknown",
			buildDate: "unknown",
			wantOut:   []string{"cdx dev"},
		},
		{
			name:      "release version with metadata",
			version:   "1.0.0",
			commit:    "abc123",
			buildDate: "2024-01-15",
			wantOut:   []string{"cdx 1.0.0", "commit: abc123", "built:  2024-01-15"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			origVersion, origCommit, origBuildDate := Version, Commit, BuildDate
			defer func() {
				Version, Commit, BuildDate = origVersion, origCommit, origBuildDate
			}()

			Version = tt.version
			Commit = tt.commit
			BuildDate = tt.buildDate

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs([]string{"version"})

			if err := rootCmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			out := buf.String()
			for _, want := range tt.wantOut {
				if !strings.Contains(out, want) {
					t.Errorf("output = %q, want substring %q", out, want)
				}
			}
		})
	}
}

func TestRootCommand_Help(t *testing.T) {
	buf := new(bytes.Buffer)
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}

	out := buf.String()

	wantSubstrings := []string{
		"cdx",
		"Code Explorer",
		"Fast",
		"--output",
		"--no-color",
		"version",
	}

	for _, want := range wantSubstrings {
		if !strings.Contains(out, want) {
			t.Errorf("help output missing %q", want)
		}
	}
}

func TestRootCommand_Flags(t *testing.T) {
	tests := []struct {
		name       string
		args       []string
		wantFormat string
		wantColor  bool
	}{
		{
			name:       "default values",
			args:       []string{},
			wantFormat: "auto",
			wantColor:  false,
		},
		{
			name:       "json output",
			args:       []string{"--output", "json"},
			wantFormat: "json",
			wantColor:  false,
		},
		{
			name:       "short output flag",
			args:       []string{"-o", "plain"},
			wantFormat: "plain",
			wantColor:  false,
		},
		{
			name:       "no color",
			args:       []string{"--no-color"},
			wantFormat: "auto",
			wantColor:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			outputFormat = "auto"
			noColor = false

			buf := new(bytes.Buffer)
			rootCmd.SetOut(buf)
			rootCmd.SetErr(buf)
			rootCmd.SetArgs(tt.args)

			err := rootCmd.Execute()
			if err != nil {
				t.Errorf("Execute() resulted in unexpected err = %v", err)
			}

			if got := GetOutputFormat(); got != tt.wantFormat {
				t.Errorf("GetOutputFormat() = %q, want %q", got, tt.wantFormat)
			}
			if got := GetNoColor(); got != tt.wantColor {
				t.Errorf("GetNoColor() = %v, want %v", got, tt.wantColor)
			}
		})
	}
}

func TestRootCommand_HasVersionSubcommand(t *testing.T) {
	found := false
	for _, cmd := range rootCmd.Commands() {
		if cmd.Use == "version" {
			found = true
			break
		}
	}
	if !found {
		t.Error("root command should have 'version' subcommand")
	}
}
