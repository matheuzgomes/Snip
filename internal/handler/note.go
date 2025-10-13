package handler

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/snip/internal/note"
	"github.com/snip/internal/repository"
	"github.com/snip/internal/validation"
)

const contentLimit = 50

type Handler interface {
	CreateNote(title string, message *string, tag *string) error
	ListNotes(isAsc, verbose bool, tag *string) error
	GetNote(idStr string, verbose bool) error
	FindNotes(term string) error
	UpdateNote(idStr string, title string) error
	DeleteNote(idStr string) error
	PatchNote(idStr string, title *string, tag *string) error
	GetRecentNotes(limit int) error
	ExportNotesToJson(since string) error
	BackupDatabase() error
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

	contentStr, err := HandleMessage(message, h)
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

	for _, note := range notes {
		tags := strings.Join(note.Tags, ", ")
		fmt.Fprintf(writer, "● #%d  %s [%s]\n", note.ID, note.Title, tags)

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
	tags := strings.Join(note.Tags, ", ")

	fmt.Printf("● #%d  %s [%s]\n", note.ID, note.Title, tags)

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

func (h *handler) PatchNote(idStr string, title *string, tag *string) error {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid note ID: %d", id)
	}

	err = h.noteRepo.CheckByID(id)
	if err != nil {
		return fmt.Errorf("failed to fetch note: %w", err)
	}

	if title != nil && *title != "" {
		if err := h.noteRepo.Patch(id, *title); err != nil {
			return fmt.Errorf("failed to update note: %w", err)
		}
	}

	if tag != nil && *tag != "" {
		if err := h.noteRepo.RemoveTagFromNote(id); err != nil {
			return fmt.Errorf("failed to remove tag from note: %w", err)
		}
		if err := h.AssociateTagsWithNote(tag, id); err != nil {
			return fmt.Errorf("failed to add tag to note: %w", err)
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

func HandleMessage(message *string, h *handler) (string, error) {
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

func (h *handler) GetRecentNotes(limit int) error {
	notes, err := h.noteRepo.GetRecent(limit)
	if err != nil {
		return fmt.Errorf("failed to get recent notes: %w", err)
	}

	if len(notes) == 0 {
		fmt.Println("No notes found.")
		return nil
	}

	fmt.Printf("Found %d note(s):\n\n", len(notes))

	for _, note := range notes {
		tags := strings.Join(note.Tags, ", ")
		fmt.Printf("● #%d  %s [%s]\n", note.ID, note.Title, tags)

		content := note.Content
		if len(content) > contentLimit {
			content = content[:contentLimit] + "..."
		}

		fmt.Printf("  └─ %s\n", content)

		fmt.Println()
	}

	return nil
}

func (h *handler) ExportNotesToJson(since string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	exportDir := filepath.Join(homeDir, ".snip", "export")
	if err := os.MkdirAll(exportDir, 0755); err != nil {
		return fmt.Errorf("failed to create export directory: %w", err)
	}

	// Parse since filter if provided
	var sinceTime *time.Time
	if since != "" {
		parsed, err := parseSinceFilter(since)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		sinceTime = &parsed
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("notes_%s.json", timestamp)
	filepath := filepath.Join(exportDir, filename)

	tempFile := filepath + ".tmp"
	f, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create export file: %w", err)
	}

	if err := h.noteRepo.ExportNotesToJsonStream(f, sinceTime); err != nil {
		f.Close()
		os.Remove(tempFile)
		return fmt.Errorf("failed to export notes: %w", err)
	}

	if err := f.Close(); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to close export file: %w", err)
	}

	if err := os.Rename(tempFile, filepath); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to finalize export: %w", err)
	}

	if sinceTime != nil {
		fmt.Printf("✓ Notes exported successfully (since %s)!\n", sinceTime.Format("2006-01-02"))
	} else {
		fmt.Printf("✓ Notes exported successfully!\n")
	}
	fmt.Printf("  Location: %s\n", filepath)
	return nil
}

func (h *handler) BackupDatabase() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	snipDir := filepath.Join(homeDir, ".snip")
	sourceDB := filepath.Join(snipDir, "notes.db")

	if _, err := os.Stat(sourceDB); os.IsNotExist(err) {
		return fmt.Errorf("database not found at %s", sourceDB)
	}

	backupDir := filepath.Join(snipDir, "backups")
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("notes_%s.db", timestamp)
	destDB := filepath.Join(backupDir, filename)

	tempFile := destDB + ".tmp"

	if err := copyFile(sourceDB, tempFile); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to backup database: %w", err)
	}

	if err := os.Rename(tempFile, destDB); err != nil {
		os.Remove(tempFile)
		return fmt.Errorf("failed to finalize backup: %w", err)
	}

	fmt.Printf("✓ Database backed up successfully!\n")
	fmt.Printf("  Location: %s\n", destDB)
	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	buf := make([]byte, 1024*1024)
	for {
		n, readErr := sourceFile.Read(buf)
		if n > 0 {
			if _, writeErr := destFile.Write(buf[:n]); writeErr != nil {
				return writeErr
			}
		}
		if readErr != nil {
			if readErr.Error() == "EOF" {
				break
			}
			return readErr
		}
	}

	return destFile.Sync()
}

func parseSinceFilter(since string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", since); err == nil {
		return t, nil
	}

	if len(since) < 2 {
		return time.Time{}, fmt.Errorf("invalid format: %s (use '2025-01-01' or '30d')", since)
	}

	unit := since[len(since)-1:]
	valueStr := since[:len(since)-1]

	var value int
	if _, err := fmt.Sscanf(valueStr, "%d", &value); err != nil {
		return time.Time{}, fmt.Errorf("invalid number in duration: %s", since)
	}

	var duration time.Duration
	switch unit {
	case "d":
		duration = time.Duration(value) * 24 * time.Hour
	case "w":
		duration = time.Duration(value) * 7 * 24 * time.Hour
	case "m":
		duration = time.Duration(value) * 30 * 24 * time.Hour
	case "y":
		duration = time.Duration(value) * 365 * 24 * time.Hour
	default:
		return time.Time{}, fmt.Errorf("invalid duration unit: %s (use d, w, m, or y)", unit)
	}

	return time.Now().Add(-duration), nil
}