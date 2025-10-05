package cmd

import (
	"fmt"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var patchTitle string
var patchTag string

func init() {
	patchCmd.Flags().StringVarP(
		&patchTitle,
		"title",
		"t",
		"",
		"If you want to update the title, you can use this flag e.g. --title 'New Title'",
	)
	patchCmd.Flags().StringVarP(
		&patchTag,
		"tag",
		"a",
		"",
		"If you want to update the tag, you can use this flag e.g. --tag 'Tag' or --tag 'Tag1 Tag2'",
	)
}

var patchCmd = &cobra.Command{
	Use:   "patch [id]",
	Short: "Patch an existing note's title and/or its tags",
	Long: `Patch an existing note by opening your default editor to modify its content.
You can also change the note's title using the --title flag.

Flags:
  --title, -t    Update the note's title (optional)
  --tag, -a      Update the note's tag (optional)

Examples:
  snip patch 1                           # Patch note 1
  snip patch 1 --title "New Title"       # Patch note 1 with new title
  snip patch 42 --tag "Meeting"  		 # Patch note 42 with new tag
  snip patch 42 --title "New Title" --tag "Meeting"  # Patch note 42 with new title and tag
  snip patch 42 --title "New Title" --tag "Meeting Technology"  # Patch note 42 with new title and two new tags`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := executeWithHandler(func(h handler.Handler) error {
			return h.PatchNote(args[0], &patchTitle, &patchTag)
		}); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	},
}
