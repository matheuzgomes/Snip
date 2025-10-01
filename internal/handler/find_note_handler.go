package handler

import "fmt"



func (h *handler) FindNotes(term string) error {
	notes, err := h.noteRepo.Search(term)
	if err != nil {
		return fmt.Errorf("failed to search notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	for i, note := range notes {
		fmt.Printf("[%d] Title: %s\n", note.ID, note.Title)
		fmt.Printf("Content: %s\n", note.Content)
		if i < len(notes)-1 {
			fmt.Println("---")
		}
	}

	return nil
}