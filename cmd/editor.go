package cmd

import (
	"fmt"
	"runtime"

	"github.com/snip/internal/handler"
	"github.com/spf13/cobra"
)

var editorCmd = &cobra.Command{
	Use:   "editor",
	Short: "Show editor information and available options",
	Long: `Display information about the currently detected editor and list all available editors on your system.

This command helps you understand which editor Snip will use for editing notes and shows alternatives you can configure.`,
	Run: func(cmd *cobra.Command, args []string) {
		showEditorInfo()
	},
}

func showEditorInfo() {
	currentEditor := getCurrentEditor()
	fmt.Printf("┌─┐ Platform: %s\n", runtime.GOOS)
	fmt.Printf("└─┘ Current Editor: %s\n\n", currentEditor)

	fmt.Println("┌─ Available Editors:")
	availableEditors := handler.ListAvailableEditors()

	for _, editor := range availableEditors {
		status := "├─"
		if editor == currentEditor {
			status = "└─ [CURRENT]"
		}
		fmt.Printf("  %s %s\n", status, editor)
	}

	fmt.Println("\n┌─ Configuration:")
	fmt.Println("├─ Set EDITOR environment variable to override")

	switch runtime.GOOS {
	case "windows":
		fmt.Println("├─ Windows: set EDITOR=code")
		fmt.Println("└─ PowerShell: $env:EDITOR='code'")
	default:
		fmt.Println("├─ Linux/macOS: export EDITOR=code")
	}
}

func getCurrentEditor() string {
	editorHandler := handler.NewEditorHandler()
	return editorHandler.GetDetectedEditor()
}

func init() {
	rootCmd.AddCommand(editorCmd)
}
