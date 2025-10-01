package handler

import (
	"fmt"
	"strconv"
)

func (h *handler) GetNote(idStr string, verbose bool) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid note ID: %s", idStr)
	}

	note, err := h.noteRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch note -> %w", err)
	}

	fmt.Printf("● #%d  %s\n", note.ID, note.Title)

	if note.Content != "" {
		fmt.Printf("  └─ %s\n", note.Content)
	}

	if verbose {
		fmt.Printf("  Created: %s\n", note.CreatedAt.Format(h.dateFormat))
		fmt.Printf("  Updated: %s\n", note.UpdatedAt.Format(h.dateFormat))
	}

	return nil
}
