package cli

import (
	"context"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/bashhack/cdx/internal/output"
	"github.com/bashhack/cdx/internal/search"
)

const (
	defaultMaxResults    = 10
	defaultSearchTimeout = 30 * time.Second
)

var (
	defLang         string
	defAll          bool
	defContextLines int
)

var defCmd = &cobra.Command{
	Use:   "def <symbol>",
	Short: "Find where a symbol is defined",
	Long: `Find where a symbol (function, type, method, etc.) is defined in the codebase.

Examples:
  cdx def GetUserByID           # Find definition of GetUserByID
  cdx def GetUserByID -C 5      # Show 5 lines of context
  cdx def UserService --lang=ts # Search TypeScript files only
  cdx def Config -o json        # Output as JSON`,
	Args: cobra.ExactArgs(1),
	RunE: runDef,
}

func init() {
	defCmd.Flags().StringVarP(&defLang, "lang", "l", "", "Force language (go, ts, js, py, rust)")
	defCmd.Flags().BoolVarP(&defAll, "all", "a", false, "Include test files and show all results (no limit)")
	defCmd.Flags().IntVarP(&defContextLines, "context", "C", 0, "Lines of context around definition")

	rootCmd.AddCommand(defCmd)
}

func runDef(cmd *cobra.Command, args []string) error {
	symbol := args[0]

	// Get current directory. If Getwd fails (e.g., directory was deleted),
	// fall back to "." silently - it's semantically equivalent and this edge
	// case doesn't warrant logging infrastructure.
	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	// Create searcher
	searcher := search.NewGrepSearcher(dir)

	// Build search options
	opts := search.Options{
		Language:     defLang,
		Context:      defContextLines,
		IncludeTests: defAll,
		Directory:    dir,
	}

	if !defAll {
		opts.MaxResults = defaultMaxResults
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), defaultSearchTimeout)
	defer cancel()

	// Find definitions
	results, err := searcher.FindDefinition(ctx, symbol, opts)

	// Determine output format
	format := output.Format(outputFormat)
	formatter := output.New(format, noColor)

	// Handle output
	w := cmd.OutOrStdout()

	if err != nil {
		// Format error output - we handle all error display ourselves
		if fmtErr := formatter.FormatError(w, err); fmtErr != nil {
			return fmtErr
		}
		// Not found is a special case - exit code 3 per COMMANDS.md
		if _, ok := err.(search.ErrNotFound); ok {
			return ExitError{Code: 3, Err: err}
		}
		return err
	}

	return formatter.FormatResults(w, results)
}
