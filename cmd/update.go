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
		"If you want to update the title of the note, you can use this flag e.g. --title 'New Title'",
	)
}

var updateCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a note",
	Long:  `Update a note with a content.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.UpdateNote(args[0], title)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
