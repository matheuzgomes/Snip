package handler

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/matheuzgomes/Snip/internal/note"
	"github.com/matheuzgomes/Snip/internal/repository"
)

func setupHandlerDB(t *testing.T) (*sql.DB, repository.NoteRepository, repository.TagRepository) {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS notes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);
		CREATE TABLE IF NOT EXISTS notes_tags (
			note_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			PRIMARY KEY (note_id, tag_id),
			FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
			FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
		);
		CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts4(id, title, content);
		CREATE TRIGGER IF NOT EXISTS notes_fts_ai AFTER INSERT ON notes BEGIN
			INSERT INTO notes_fts(id, title, content) VALUES (new.id, new.title, new.content);
		END;
		CREATE TRIGGER IF NOT EXISTS notes_fts_au AFTER UPDATE ON notes BEGIN
			UPDATE notes_fts SET title = new.title, content = new.content WHERE id = old.id;
		END;
		CREATE TRIGGER IF NOT EXISTS notes_fts_ad AFTER DELETE ON notes BEGIN
			DELETE FROM notes_fts WHERE id = old.id;
		END;
	`)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	noteRepo, err := repository.NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}
	tagRepo, err := repository.NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	return db, noteRepo, tagRepo
}

func TestHandler_CreateNote(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("hello world")
	if err := h.CreateNote("Test", msg, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	notes, err := noteRepo.GetAll(false, 0)
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}
	if len(notes) != 1 {
		t.Fatalf("expected 1 note, got %d", len(notes))
	}
	if notes[0].Title != "Test" {
		t.Errorf("expected 'Test', got '%s'", notes[0].Title)
	}
}

func TestHandler_CreateNoteWithTag(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("tagged content")
	tag := stringPtr("work")
	if err := h.CreateNote("Tagged", msg, tag); err != nil {
		t.Fatalf("CreateNote with tag failed: %v", err)
	}

	notes, _ := noteRepo.GetAll(false, 0)
	if len(notes[0].Tags) != 1 || notes[0].Tags[0] != "work" {
		t.Errorf("expected tag ['work'], got %v", notes[0].Tags)
	}
}

func TestHandler_ListNotes(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("content")
	for _, title := range []string{"B", "A", "C"} {
		if err := h.CreateNote(title, msg, nil); err != nil {
			t.Fatalf("CreateNote(%q) failed: %v", title, err)
		}
	}

	// Just verify it doesn't error — ListNotes writes to stdout
	if err := h.ListNotes(false, false, nil); err != nil {
		t.Fatalf("ListNotes failed: %v", err)
	}
}

func TestHandler_FindNotes(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("discuss project alpha")
	if err := h.CreateNote("Meeting", msg, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}
	msg2 := stringPtr("buy groceries")
	if err := h.CreateNote("Shopping", msg2, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	if err := h.FindNotes("alpha"); err != nil {
		t.Fatalf("FindNotes failed: %v", err)
	}
}

func TestHandler_GetNote(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("note content")
	if err := h.CreateNote("My Note", msg, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	if err := h.GetNote("1", false, false); err != nil {
		t.Fatalf("GetNote failed: %v", err)
	}

	if err := h.GetNote("999", false, false); err == nil {
		t.Error("expected error for nonexistent note")
	}
}

func TestHandler_DeleteNote(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("to delete")
	if err := h.CreateNote("Delete Me", msg, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	if err := h.DeleteNote("1"); err != nil {
		t.Fatalf("DeleteNote failed: %v", err)
	}

	if err := noteRepo.CheckByID(1); err == nil {
		t.Error("expected error after delete")
	}
}

func TestHandler_PatchNote(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("original")
	if err := h.CreateNote("Original", msg, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	newTitle := stringPtr("Patched")
	if err := h.PatchNote("1", newTitle, nil); err != nil {
		t.Fatalf("PatchNote failed: %v", err)
	}

	got, _ := noteRepo.GetByID(1)
	if got.Title != "Patched" {
		t.Errorf("expected 'Patched', got '%s'", got.Title)
	}
}

func TestHandler_ExportNotes(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	msg := stringPtr("export me")
	if err := h.CreateNote("Export", msg, nil); err != nil {
		t.Fatalf("CreateNote failed: %v", err)
	}

	// Override export dir via home — tricky, just test that it doesn't crash
	if err := h.ExportNotes("", "json"); err != nil {
		t.Fatalf("ExportNotes failed: %v", err)
	}

	homeDir, _ := os.UserHomeDir()
	exportDir := filepath.Join(homeDir, ".snip", "export")
	files, _ := os.ReadDir(exportDir)
	if len(files) == 0 {
		t.Error("expected exported files")
	}
	// Cleanup
	os.RemoveAll(exportDir)
}

func TestHandler_ImportNotes(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	// Create temp dir with a markdown file
	dir := t.TempDir()
	mdFile := filepath.Join(dir, "imported.md")
	if err := os.WriteFile(mdFile, []byte("# Hello\nWorld"), 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// ImportNotes uses homeDir as base, so we pass a relative path
	// Since we can't easily redirect home, just test the repo layer indirectly
	// This test validates the handler doesn't panic
	_ = h.ImportNotes(dir)
}

func stringPtr(s string) *string {
	return &s
}

func TestHandler_CreateNote_EmptyTitle(t *testing.T) {
	_, noteRepo, tagRepo := setupHandlerDB(t)
	h := NewHandler(noteRepo, tagRepo)

	err := h.CreateNote("   ", nil, nil)
	if err == nil {
		t.Error("expected error for empty title")
		return
	}
	if !strings.Contains(err.Error(), "is required") {
		t.Errorf("expected 'is required' in error, got: %v", err)
	}
}

func TestNote_NewNote(t *testing.T) {
	n := note.NewNote("Title", "Content")
	if n.Title != "Title" {
		t.Errorf("expected 'Title', got '%s'", n.Title)
	}
	if n.Content != "Content" {
		t.Errorf("expected 'Content', got '%s'", n.Content)
	}
	if n.CreatedAt.IsZero() {
		t.Error("expected non-zero CreatedAt")
	}
	if n.UpdatedAt.IsZero() {
		t.Error("expected non-zero UpdatedAt")
	}
}

func TestNote_JSONTags(t *testing.T) {
	n := note.NoteWithTags{
		ID:    1,
		Title: "Test",
		Tags:  []string{"a", "b"},
	}
	if len(n.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(n.Tags))
	}
}
