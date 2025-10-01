package handler

import "fmt"


func (h *handler) ListNotes(isAsc, verbose bool) error {
	notes, err := h.noteRepo.GetAll(isAsc)
	if err != nil {
		return fmt.Errorf("failed to fetch notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	const contentLimit = 50

	fmt.Printf("Found %d note(s):\n\n", len(notes))

	for i, note := range notes {
		fmt.Printf("[%d] Title: %s\n", note.ID, note.Title)

		content := note.Content
		if len(content) > contentLimit {
			content = content[:contentLimit] + "..."
		}
		fmt.Printf("Content: %s\n", content)

		if verbose {
			fmt.Printf("Created: %s\n", note.CreatedAt.Format(h.dateFormat))
			fmt.Printf("Updated: %s\n", note.UpdatedAt.Format(h.dateFormat))
		}

		if i < len(notes)-1 {
			fmt.Println("---")
		}
	}

	return nil
}