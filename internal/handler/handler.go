package handler

import (
	"fmt"
	"strconv"

	"github.com/snip/internal/note"
	"github.com/snip/internal/validation"
)

type Handler interface {
	CreateNote(title string) error
	ListNotes(isAsc, verbose bool) error
	GetNote(idStr string, verbose bool) error
	FindNotes(term string) error
	UpdateNote(idStr string, title string) error
	DeleteNote(idStr string) error
}

type handler struct {
	noteRepo  note.Repository
	validator *validation.Validator
	editorHandler *EditorHandler
	dateFormat string
}

func NewHandler(repo note.Repository) Handler {
	return &handler{
		noteRepo:  repo,
		validator: validation.NewValidator(),
		dateFormat: "2006-01-02 15:04:05",
		editorHandler: NewEditorHandler(),
	}
}

func (h *handler) CreateNote(title string) error {
    if err := h.validator.ValidateNote(title); err != nil {
        return err
    }

    tempFile, err := h.editorHandler.HandleEditor("")
    if err != nil {
        return err
    }
    defer h.editorHandler.RemoveTempFile(tempFile)
    
    content, err := h.editorHandler.ReadTempFile(tempFile)
    if err != nil {
        return err
    }
    
    contentStr := string(content)
    
    newNote := note.NewNote(title, contentStr)
    if err := h.noteRepo.Create(newNote); err != nil {
        return fmt.Errorf("failed to create note: %w", err)
    }
    
    fmt.Printf("Note created successfully!\n")
    fmt.Printf("[%d] Title: %s\n", newNote.ID, newNote.Title)
    
    return nil
}

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

func (h *handler) GetNote(idStr string, verbose bool) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid note ID: %s", idStr)
	}

	note, err := h.noteRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch note -> %w", err)
	}


	fmt.Println("---")
	fmt.Printf("[%d] Title: %s\n", note.ID, note.Title)
	fmt.Printf("Content: %s\n", note.Content)

	if verbose {
		fmt.Printf("Created: %s\n", note.CreatedAt.Format(h.dateFormat))
		fmt.Printf("Updated: %s\n", note.UpdatedAt.Format(h.dateFormat))
	}
	fmt.Println("---")

	return nil
}

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
        return fmt.Errorf("failed to create note: %w", err)
    }
    
    fmt.Printf("Note updated successfully!\n")
    
    return nil
}

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
