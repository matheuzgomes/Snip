package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var exportSince string

func init() {
	exportCmd.Flags().StringVarP(&exportSince, "since", "s", "", "Export notes created since date or duration (e.g., '2025-01-01' or '30d')")
}

var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export notes to JSON format",
	Long: `Export your notes to a timestamped JSON file for migration or archival purposes.

The export creates a JSON array containing notes with their metadata, content, and tags.
Exports are stored in ~/.snip/export/

Note: For backup purposes, use 'snip backup' instead, which is faster and preserves
the complete database structure.

Export is recommended when you need to:
  - Migrate notes to another system
  - Share notes in a portable format
  - Archive notes in a human-readable format

Flags:
  --since, -s    Export only notes created since a specific date or duration
                 Examples: "2025-01-01", "30d", "7d", "1y"

Examples:
  snip export                      # Export all notes
  snip export --since 30d          # Export notes from last 30 days
  snip export --since "2025-01-01" # Export notes since Jan 1, 2025
  snip export -s 7d                # Export notes from last week`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.ExportNotesToJson(exportSince)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
