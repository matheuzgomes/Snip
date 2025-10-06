package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

func init() {
	showCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Show more information about the notes")
}

var showCmd = &cobra.Command{
	Use:   "show [id]",
	Short: "Display the content of a specific note by ID",
	Long: `Display the full content of a specific note identified by its unique ID.

This command shows the note's title, content and tags in a readable format. Use the verbose
flag to see additional metadata like creation and modification timestamps.

Flags:
  --verbose, -v  Show detailed metadata (timestamps, ID, etc.)

Examples:
  snip show 1              # Display note with ID 1
  snip show 42             # Display note with ID 42  
  snip show 1 --verbose    # Show note 1 with full metadata
  snip show 1 -v           # Same as above (short flag)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.GetNote(args[0], verbose)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
