package cli

import (
	"github.com/spf13/cobra"
)

// Version information (set at build time via ldflags)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildDate = "unknown"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version, commit hash, and build date of cdx.`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("cdx %s\n", Version)
		if Commit != "unknown" || BuildDate != "unknown" {
			cmd.Printf("  commit: %s\n", Commit)
			cmd.Printf("  built:  %s\n", BuildDate)
		}
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
