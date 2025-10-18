package cmd

import (
	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var editorCmd = &cobra.Command{
	Use:   "editor",
	Short: "Show editor information and available options",
	Long: `Display information about the currently detected editor and list all available editors on your system.

This command helps you understand which editor Snip will use for editing notes and shows alternatives you can configure.`,
	Run: func(cmd *cobra.Command, args []string) {
		editorHandler := handler.NewEditorHandler()
		editorHandler.ShowEditorInfo()
	},
}