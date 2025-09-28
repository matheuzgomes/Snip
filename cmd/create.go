package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new note with title and interactive content editing",
	Long: `Create a new note by providing a title and opening your default editor for content.

The title can be multiple words and will be joined together. After running this command,
your default editor will open allowing you to write the note content. Save and close
the editor to create the note.

Examples:
  snip create "My Daily Notes"     # Creates note with title "My Daily Notes"
  snip create Meeting Notes        # Creates note with title "Meeting Notes"
  snip create TODO                 # Creates note with title "TODO"`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.CreateNote(strings.Join(args, " "))
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
