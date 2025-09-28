package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Permanently remove a note by ID",
	Long: `Permanently delete a note from your collection using its unique ID.

This action cannot be undone. The note will be completely removed from your database.
Make sure you have the correct note ID before executing this command.

Examples:
  snip delete 1        # Delete note with ID 1
  snip delete 42       # Delete note with ID 42
  
Tip: Use 'snip list' or 'snip show [id]' to verify the note before deletion.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.DeleteNote(args[0])
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
