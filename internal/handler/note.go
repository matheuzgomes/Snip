package handler

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/snip/internal/note"
	"github.com/snip/internal/repository"
	"github.com/snip/internal/validation"
)

type Handler interface {
	CreateNote(title string, message *string, tag *string) error
	ListNotes(isAsc, verbose bool, tag *string) error
	GetNote(idStr string, verbose bool) error
	FindNotes(term string) error
	UpdateNote(idStr string, title string) error
	DeleteNote(idStr string) error
}

type handler struct {
	noteRepo      repository.NoteRepository
	tagRepo       repository.TagRepository
	validator     *validation.Validator
	editorHandler *EditorHandler
	dateFormat    string
}

func NewHandler(noteRepo repository.NoteRepository, tagRepo repository.TagRepository) Handler {
	return &handler{
		noteRepo:      noteRepo,
		tagRepo:       tagRepo,
		validator:     validation.NewValidator(),
		dateFormat:    "2006-01-02 15:04:05",
		editorHandler: NewEditorHandler(),
	}
}

func (h *handler) CreateNote(title string, message *string, tag *string) error {
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

	if tag != nil && *tag != "" {
		if err := h.AssociateTagsWithNote(tag, newNote.ID); err != nil {
			return fmt.Errorf("failed to associate tags with note: %w", err)
		}
	}

	fmt.Printf("Note created successfully!\n")
	fmt.Printf("● #%d  %s\n", newNote.ID, newNote.Title)

	return nil
}

func (h *handler) ListNotes(isAsc, verbose bool, tag *string) error {
	tagID := 0

	if tag != nil && *tag != "" {
		tagObj, err := h.tagRepo.GetByName(*tag)
			if err != nil {
				return fmt.Errorf("no note found for this tag: %s", *tag)
			}
			tagID = tagObj.ID
	}

	notes, err := h.noteRepo.GetAll(isAsc, tagID)
	if err != nil {
		return fmt.Errorf("failed to fetch notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	fmt.Printf("Found %d note(s):\n\n", len(notes))

	writer := bufio.NewWriter(os.Stdout)
	const contentLimit = 50

	for _, note := range notes {
		fmt.Fprintf(writer, "● #%d  %s\n", note.ID, note.Title)

		content := note.Content
		if len(content) > contentLimit {
			content = content[:contentLimit] + "..."
		}

		fmt.Fprintf(writer, "  └─ %s\n", content)

		if verbose {
			fmt.Fprintf(writer, "  └─ Created: %s\n", note.CreatedAt.Format(h.dateFormat))
			fmt.Fprintf(writer, "  └─ Updated: %s\n", note.UpdatedAt.Format(h.dateFormat))
		}

		fmt.Fprintln(writer)
	}

	defer writer.Flush()
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

	fmt.Printf("● #%d  %s\n", note.ID, note.Title)

	if note.Content != "" {
		fmt.Printf("  └─ %s\n", note.Content)
	}

	if verbose {
		fmt.Printf("  └─ Created: %s\n", note.CreatedAt.Format(h.dateFormat))
		fmt.Printf("  └─ Updated: %s\n", note.UpdatedAt.Format(h.dateFormat))
	}

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

func (h *handler) AssociateTagsWithNote(tag *string, noteID int) error {
	for tag := range strings.SplitSeq(*tag, " ") {
		tagObj, err := h.tagRepo.GetOrCreate(tag)
		if err != nil {
				return err
			}

		if err := h.noteRepo.AddTagToNote(noteID, tagObj.ID); err != nil {
			return err
		}
	}

	return nil
}
