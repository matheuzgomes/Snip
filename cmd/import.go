package cmd

import (
	"github.com/matheuzgomes/Snip/internal/handler"
	"github.com/spf13/cobra"
)

var importDir string

func init() {
	importCmd.Flags().StringVarP(&importDir, "dir", "d", "", "Directory to import notes from")
}

var importCmd = &cobra.Command{
	Use:   "import",
	Short: "Import markdown notes from a directory",
	Long: `Import markdown notes from a directory into the database.

Flags:
  --dir, -d    Directory to import notes from starting from your home directory
Examples:
  snip import                      # Import all markdown notes from the current directory
  snip import --dir notes          # Snip will look for notes starting from your home directory so in this example it will look for notes in ~/notes
  snip import -d notes/work        # Snip will look for notes starting from your home directory so in this example it will look for notes in ~/notes/work`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeWithHandler(func(h handler.Handler) error {
			return h.ImportNotes(importDir)
		})
	},
}
