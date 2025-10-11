package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var limit int

func init() {
	recentCmd.Flags().IntVarP(&limit, "limit", "l", 10, "Limit the number of notes to display")
}

var recentCmd = &cobra.Command{
	Use:     "recent",
	Aliases: []string{"rec"},
	Short:   "Display recent notes",
	Long: `Display recent notes in a organized list format.

By default, notes are displayed with newest first (descending order by updated date).

Flags:
  --limit, -l      Limit the number of notes to display

Examples:
  snip recent                    # Show newest notes first (default)
  snip rec                       # Same as above (alias)
  snip recent --limit 10         # Show 10 recent notes
  snip recent -l 10              # Same as above (short flag)`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.GetRecentNotes(limit)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
