package handler

import (
	"fmt"

	"github.com/snip/internal/note"
)


func (h *handler) CreateNote(title string, message *string) error {
	if err := h.validator.ValidateNote(title); err != nil {
		return err
	}

	contentStr, err := CheckMessage(message, h)
	
	if err != nil {
		return err
	}

	newNote := note.NewNote(title, contentStr)
	if err := h.noteRepo.Create(newNote); err != nil {
		return fmt.Errorf("failed to create note: %w", err)
	}

	fmt.Printf("Note created successfully!\n")
	fmt.Printf("[%d] Title: %s\n", newNote.ID, newNote.Title)

	return nil
}



func CheckMessage(message *string, h *handler) (string, error) {
	if message != nil && *message != "" {
		return *message, nil
	}

	tempFile, err := h.editorHandler.HandleEditor("")
	if err != nil {
		return "", err
	}
	defer h.editorHandler.RemoveTempFile(tempFile)

	content, err := h.editorHandler.ReadTempFile(tempFile)
	if err != nil {
		return "", err
	}

	contentStr := string(content)
	return contentStr, nil
}