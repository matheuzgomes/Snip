package cmd

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "snip",
	Short: "A fast and lightweight note-taking CLI application",
	Long: `Snip is a terminal-based note management application.
It allows you to create, edit, view, and delete notes quickly and efficiently.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(showCmd)
	rootCmd.AddCommand(findCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(patchCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(editorCmd)
	rootCmd.AddCommand(recentCmd)
	rootCmd.AddCommand(backupCmd)
	rootCmd.AddCommand(exportCmd)
}
