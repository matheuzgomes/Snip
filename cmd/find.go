package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var findCmd = &cobra.Command{
	Use:   "find [text]",
	Short: "Find a note",
	Long:  `Find a note with the specified text.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.FindNotes(strings.Join(args, " "))
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
