package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var isAsc bool
var verbose bool

func init() {
	listCmd.Flags().BoolVarP(&isAsc, "asc", "a", false, "List notes in chronological order (oldest first)")
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show more information about the notes")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all notes",
	Long: `List all notes in chronological order.

By default, notes are sorted by creation date in descending order (newest first).
Use the --asc or -a flag to sort in ascending order (oldest first).
Use the --verbose or -v flag to show more information about the notes.

Examples:
snip list or snip ls  # Show newest notes first
snip list --asc, -a     # Show oldest notes first
snip list --verbose, -v # Show more information about the notes`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.ListNotes(isAsc, verbose)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
