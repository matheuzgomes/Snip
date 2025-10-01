package handler

import (
	"bufio"
	"fmt"
	"os"
)


const contentLimit = 50

func (h *handler) ListNotes(isAsc, verbose bool) error {
	notes, err := h.noteRepo.GetAll(isAsc)
	if err != nil {
		return fmt.Errorf("failed to fetch notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	fmt.Printf("Found %d note(s):\n\n", len(notes))

	writer := bufio.NewWriter(os.Stdout)

	for _, note := range notes {
		fmt.Fprintf(writer, "● #%d  %s\n", note.ID, note.Title)
		
		content := note.Content
        if len(content) > contentLimit {
            content = content[:contentLimit] + "..."
		}

		fmt.Fprintf(writer, "  └─ %s\n", content)

		if verbose {
			fmt.Fprintf(writer, "  Created: %s\n", note.CreatedAt.Format(h.dateFormat))
			fmt.Fprintf(writer, "  Updated: %s\n", note.UpdatedAt.Format(h.dateFormat))
		}

		fmt.Fprintln(writer)
	}

	writer.Flush()

	return nil
}