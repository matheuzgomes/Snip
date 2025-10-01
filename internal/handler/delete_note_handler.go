package handler

import (
	"fmt"
	"strconv"
)



func (h *handler) DeleteNote(idStr string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid note ID: %d", id)
	}

	if err := h.noteRepo.CheckByID(id); err != nil {
		return fmt.Errorf("this note does not exist: %w", err)
	}

	if err := h.noteRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete note: %w", err)
	}

	fmt.Printf("Note deleted successfully!\n")

	return nil
}