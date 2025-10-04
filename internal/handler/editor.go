package handler

import (
	"fmt"
	"os"
	"os/exec"
)

type EditorHandler struct {
}

func NewEditorHandler() *EditorHandler {
	return &EditorHandler{}
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

	editor := e.GetEditor()
	cmd := exec.Command(editor, tempFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to open editor: %w", err)
	}

	return tempFile, nil
}


func (e *EditorHandler) GetEditor() string {
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}
	if _, err := exec.LookPath("nano"); err == nil {
		return "nano"
	}
	if _, err := exec.LookPath("vim"); err == nil {
		return "vim"
	}
	return "vi"
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
