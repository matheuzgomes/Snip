package handler

import (
	"fmt"
	"strconv"
)


func (h *handler) UpdateNote(idStr string, title string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid note ID: %d", id)
	}

	note, err := h.noteRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch note: %w", err)
	}

	tempFile, err := h.editorHandler.HandleEditor(note.Content)
	if err != nil {
		return err
	}
	defer h.editorHandler.RemoveTempFile(tempFile)

	content, err := h.editorHandler.ReadTempFile(tempFile)
	if err != nil {
		return err
	}

	contentStr := string(content)
	if err := h.noteRepo.Update(id, contentStr, title); err != nil {
		return fmt.Errorf("failed to update note: %w", err)
	}

	fmt.Printf("Note updated successfully!\n")

	return nil
}