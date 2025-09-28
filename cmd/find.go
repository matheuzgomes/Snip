package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find [text]",
	Short: "Search for notes containing specific text in title or content",
	Long: `Search through all your notes to find ones containing the specified text.

The search looks through both note titles and content, returning all matches.
Multiple search terms can be provided and will be joined together. The search
is case-insensitive and matches partial words.

Examples:
  snip find meeting            # Find notes containing "meeting"
  snip find "project ideas"    # Find notes with "project ideas"
  snip find TODO urgent        # Find notes containing "TODO urgent"
  snip find golang             # Find notes about golang

Tip: Use quotes for exact phrases, or separate words for broader matching.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.FindNotes(strings.Join(args, " "))
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
