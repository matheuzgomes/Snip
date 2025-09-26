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
	Use:   "show",
	Short: "Show a specific note",
	Long: `Show a specific note by ID.

Use the --verbose or -v flag to show more information about the note.

Examples:
snip show 1          # Show note with ID 1
snip show 1 --verbose, -v # Show more information about the note`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.GetNote(args[0], verbose)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
