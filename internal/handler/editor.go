package handler

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

type EditorHandler struct {
	detectedEditor string
}

func NewEditorHandler() *EditorHandler {
	return &EditorHandler{
		detectedEditor: detectEditor(),
	}
}

func (e *EditorHandler) HandleEditor(content string) (*os.File, error) {
	tempFile, err := e.CreateTempFile()
	if err != nil {
		return nil, err
	}

	if content != "" {
		if _, err := tempFile.WriteString(content); err != nil {
			return nil, fmt.Errorf("failed to write content to temp file: %w", err)
		}
		tempFile.Close()
	}

	editor, args := e.GetEditor()
	cmd := exec.Command(editor, append(args, tempFile.Name())...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to open editor: %w", err)
	}

	return tempFile, nil
}

func detectEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	switch runtime.GOOS {
	case "windows":
		return detectWindowsEditor()
	case "darwin":
		return detectMacEditor()
	default:
		return detectUnixEditor()
	}
}

func detectWindowsEditor() string {
	editors := []string{
		"code",      // Visual Studio Code
		"notepad++", // Notepad++
		"subl",      // Sublime Text
		"atom",      // Atom
		"micro",     // Micro editor
		"nano",      // Nano (via WSL/Chocolatey)
		"vim",       // Vim (via Git Bash/Chocolatey)
		"notepad",   // Windows Notepad (fallback)
	}

	for _, editor := range editors {
		if isEditorAvailable(editor) {
			return editor
		}
	}

	return "notepad"
}

func detectMacEditor() string {
	editors := []string{
		"code", // Visual Studio Code
		"subl", // Sublime Text
		"atom", // Atom
		"nano", // Nano
		"vim",  // Vim
		"vi",   // Vi
		"open", // macOS open command
	}

	for _, editor := range editors {
		if isEditorAvailable(editor) {
			return editor
		}
	}

	return "vi"
}

func detectUnixEditor() string {
	editors := []string{
		"nano",  // Nano
		"vim",   // Vim
		"vi",    // Vi
		"micro", // Micro editor
		"code",  // Visual Studio Code
	}

	for _, editor := range editors {
		if isEditorAvailable(editor) {
			return editor
		}
	}

	return "vi"
}

func isEditorAvailable(editor string) bool {
	if editor == "notepad" && runtime.GOOS == "windows" {
		return true
	}

	_, err := exec.LookPath(editor)
	return err == nil
}

func (e *EditorHandler) GetEditor() (string, []string) {
	editor := e.detectedEditor

	switch editor {
	case "code":
		return "code", []string{"--wait"}
	case "subl":
		return "subl", []string{"--wait"}
	case "atom":
		return "atom", []string{"--wait"}
	case "open":
		return "open", []string{"-t"}
	default:
		return editor, []string{}
	}
}

func (e *EditorHandler) CreateTempFile() (*os.File, error) {
	tempFile, err := os.CreateTemp("", "snip-note-*.md")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	return tempFile, nil
}

func (e *EditorHandler) ReadTempFile(tempFile *os.File) (string, error) {
	content, err := os.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read temp file: %w", err)
	}
	return string(content), nil
}

func (e *EditorHandler) RemoveTempFile(tempFile *os.File) error {
	return os.Remove(tempFile.Name())
}

func (e *EditorHandler) GetDetectedEditor() string {
	return e.detectedEditor
}


func (e *EditorHandler) ListAvailableEditors() []string {
	var editors []string

	switch runtime.GOOS {
	case "windows":
		editors = []string{"code", "notepad++", "subl", "atom", "micro", "nano", "vim", "notepad"}
	case "darwin":
		editors = []string{"code", "subl", "atom", "nano", "vim", "vi", "open"}
	default:
		editors = []string{"nano", "vim", "vi", "micro", "code"}
	}

	var available []string
	for _, editor := range editors {
		if isEditorAvailable(editor) {
			available = append(available, editor)
		}
	}

	return available
}


func (e *EditorHandler) ShowEditorInfo() {
	currentEditor := e.GetDetectedEditor()
	fmt.Printf("┌─┐ Platform: %s\n", runtime.GOOS)
	fmt.Printf("└─┘ Current Editor: %s\n\n", currentEditor)

	fmt.Println("┌─ Available Editors:")
	availableEditors := e.ListAvailableEditors()

	for _, editor := range availableEditors {
		status := "├─"

		if editor == currentEditor {
			status = "├─ [CURRENT]"
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