// Package cli implements the command-line interface for cdx.
package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputFormat string
	noColor      bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "cdx",
	Short: "Fast codebase exploration CLI",
	Long: `cdx (Code Explorer) is a fast, native CLI for codebase exploration
with optional LLM enhancement.

Philosophy: Fast by default, smart when needed.

Examples:
  cdx def MyFunction     # Find definition of MyFunction
  cdx refs MyFunction    # Find references to MyFunction
  cdx outline main.go    # Show structure of main.go`,
}

// Execute runs the root command.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "auto",
		"Output format: auto, human, json, plain")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false,
		"Disable color output")
}

// GetOutputFormat returns the current output format setting.
func GetOutputFormat() string {
	return outputFormat
}

// GetNoColor returns whether color is disabled.
func GetNoColor() bool {
	return noColor
}
