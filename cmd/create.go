package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new note",
	Long:  `Create a new note with the specified title and content.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.CreateNote(strings.Join(args, " "))
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
