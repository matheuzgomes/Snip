package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Create a backup of your notes database",
	Long: `Create a timestamped backup of your notes database.

The backup is a complete copy of the SQLite database file, preserving all notes,
tags, relationships, and metadata. Backups are stored in ~/.snip/backups/

This is the recommended method for backing up your notes as it:
  - Preserves the complete database structure
  - Is fast and reliable
  - Can be easily restored by copying back to ~/.snip/notes.db
  - Takes less space than JSON exports

Examples:
  snip backup          # Create a backup with current timestamp`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.BackupDatabase()
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
