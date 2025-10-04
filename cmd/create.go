package cmd

import (
	"fmt"
	"strings"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var message string

var tag string

func init() {
	createCmd.Flags().StringVarP(&message, "message", "m", "", "Content of the note")
	createCmd.Flags().StringVarP(&tag, "tag", "t", "", "Tag of the note")
}

var createCmd = &cobra.Command{
	Use:   "create [title]",
	Short: "Create a new note with title and optional content",
	Long: `Create a new note by providing a title and optional content.

The title can be multiple words and will be joined together. You can provide content
in two ways:
1. Use the --message flag to provide content directly
2. If no message is provided, your default editor will open for interactive content editing

Examples:
  snip create "My Daily Notes"                    # Opens editor for content
  snip create "Quick Note" --message "Hello!"     # Uses provided message
  snip create Meeting Notes                       # Opens editor for content
  snip create TODO --message "Buy milk"           # Uses provided message`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			stringPtr := func(s string) *string {
				if s == "" {
					return nil
				}
				return &s
			}
			return h.CreateNote(strings.Join(args, " "), stringPtr(message), stringPtr(tag))
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
