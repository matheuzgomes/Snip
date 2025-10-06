package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/snip/internal/validation"
	"github.com/spf13/cobra"
)

var isAsc bool
var verbose bool
var listTag string

func init() {
	listCmd.Flags().BoolVarP(&isAsc, "asc", "a", false, "List notes in chronological order (oldest first)")
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show more information about the notes")
	listCmd.Flags().StringVarP(&listTag, "tag", "t", "", "List notes by tag")
}

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "Display all notes with sorting and detail options",
	Long: `Display all your notes in a organized list format with customizable sorting and verbosity.

By default, notes are displayed with newest first (descending order by creation date).
You can control the output format and sorting to match your workflow preferences.

You can also list notes by tag using the --tag flag.

Flags:
  --asc, -a      Sort chronologically (oldest first)
  --verbose, -v  Show detailed information including timestamps and IDs

Examples:
  snip list                    # Show newest notes first (default)
  snip ls                      # Same as above (alias)
  snip list --asc              # Show oldest notes first
  snip list -v                 # Show detailed note information
  snip list --asc --verbose    # Oldest first with full details
  snip list --tag "tag"        # List notes by tag`,
	Run: func(cmd *cobra.Command, args []string) {
		validator := validation.NewValidator()
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.ListNotes(isAsc, verbose, validator.CheckString(listTag))
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
