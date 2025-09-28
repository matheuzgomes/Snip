package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var title string

func init() {
	updateCmd.Flags().StringVarP(
		&title,
		"title",
		"t",
		"",
		"If you want to update the title, you can use this flag e.g. --title 'New Title'",
	)
}

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Modify an existing note's content and optionally its title",
	Long: `Update an existing note by opening your default editor to modify its content.
You can also change the note's title using the --title flag.

The note's content will open in your default editor where you can make changes.
Save and close the editor to apply the updates. The modification timestamp
will be automatically updated.

Flags:
  --title, -t    Update the note's title (optional)

Examples:
  snip update 1                           # Edit content of note 1
  snip update 1 --title "New Title"      # Edit content and change title
  snip update 42 -t "Updated Meeting"    # Edit note 42 with new title`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.UpdateNote(args[0], title)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
