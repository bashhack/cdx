// Package cli implements the command-line interface for cdx.
package cli

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	outputFormat string
	noColor      bool
)

// ExitError is an error that carries a specific exit code.
// Commands can return this to signal a non-standard exit code.
type ExitError struct {
	Err  error
	Code int
}

func (e ExitError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return ""
}

func (e ExitError) Unwrap() error {
	return e.Err
}

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

// Execute runs the root command and exits on error.
// This is the main entry point for the CLI binary.
func Execute() {
	if err := ExecuteE(); err != nil {
		// Check for ExitError with custom exit code
		var exitErr ExitError
		if errors.As(err, &exitErr) {
			os.Exit(exitErr.Code)
		}
		os.Exit(1)
	}
}

// ExecuteE runs the root command and returns any error.
// This is useful for testing and programmatic use.
func ExecuteE() error {
	return rootCmd.Execute()
}

func init() {
	// Silence Cobra's built-in error output - we handle errors ourselves via formatters
	rootCmd.SilenceErrors = true
	// Don't show usage on errors - only on --help
	rootCmd.SilenceUsage = true

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
