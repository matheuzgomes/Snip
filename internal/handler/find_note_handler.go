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

	fmt.Printf("Found %d note(s) matching '%s':\n\n", len(notes), term)

	for _, note := range notes {
		fmt.Printf("● #%d  %s\n", note.ID, note.Title)

		if note.Content != "" {
			fmt.Printf("  └─ %s\n", note.Content)
		}

		fmt.Println()
	}

	return nil
}
