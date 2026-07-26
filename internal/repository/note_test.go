package repository

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/matheuzgomes/Snip/internal/note"
	"github.com/matheuzgomes/Snip/internal/tag"
)

func setupNoteDB(t *testing.T) *sql.DB {
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

		CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
		CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);

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

	return db
}

func TestNoteRepo_CreateAndGetByID(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	n := note.NewNote("Test Note", "Hello World")
	if err := repo.Create(n); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if n.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	got, err := repo.GetByID(n.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if got.Title != "Test Note" {
		t.Errorf("expected title 'Test Note', got '%s'", got.Title)
	}
	if got.Content != "Hello World" {
		t.Errorf("expected content 'Hello World', got '%s'", got.Content)
	}

	_, err = repo.GetByID(9999)
	if err == nil {
		t.Error("expected error for nonexistent note")
	}
}

func TestNoteRepo_GetAll(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	for i, title := range []string{"A", "B", "C"} {
		n := note.NewNote(title, "content")
		time.Sleep(1 * time.Millisecond) // ensure distinct timestamps
		if err := repo.Create(n); err != nil {
			t.Fatalf("Create(%q) failed: %v", title, err)
		}
		// quick sanity
		if n.ID != i+1 {
			t.Errorf("expected ID %d, got %d", i+1, n.ID)
		}
	}

	// Desc (default)
	desc, err := repo.GetAll(false, 0)
	if err != nil {
		t.Fatalf("GetAll desc failed: %v", err)
	}
	if len(desc) != 3 {
		t.Fatalf("expected 3 notes, got %d", len(desc))
	}
	if desc[0].Title != "C" {
		t.Errorf("expected 'C' first in desc, got '%s'", desc[0].Title)
	}

	// Asc
	asc, err := repo.GetAll(true, 0)
	if err != nil {
		t.Fatalf("GetAll asc failed: %v", err)
	}
	if asc[0].Title != "A" {
		t.Errorf("expected 'A' first in asc, got '%s'", asc[0].Title)
	}
}

func TestNoteRepo_Update(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	n := note.NewNote("Original", "old content")
	if err := repo.Create(n); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	oldUpdated := n.UpdatedAt

	if err := repo.Update(n.ID, "new content", "New Title"); err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	got, err := repo.GetByID(n.ID)
	if err != nil {
		t.Fatalf("GetByID after update failed: %v", err)
	}
	if got.Title != "New Title" {
		t.Errorf("expected title 'New Title', got '%s'", got.Title)
	}
	if got.Content != "new content" {
		t.Errorf("expected content 'new content', got '%s'", got.Content)
	}
	if !got.UpdatedAt.After(oldUpdated) {
		t.Error("expected updated_at to be after original")
	}
}

func TestNoteRepo_Patch(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	n := note.NewNote("Old Title", "content")
	if err := repo.Create(n); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	oldUpdated := n.UpdatedAt

	if err := repo.Patch(n.ID, "Patched Title"); err != nil {
		t.Fatalf("Patch failed: %v", err)
	}

	got, err := repo.GetByID(n.ID)
	if err != nil {
		t.Fatalf("GetByID after patch failed: %v", err)
	}
	if got.Title != "Patched Title" {
		t.Errorf("expected 'Patched Title', got '%s'", got.Title)
	}
	if got.Content != "content" {
		t.Error("content should not change on patch")
	}
	if !got.UpdatedAt.After(oldUpdated) {
		t.Error("expected updated_at to be updated on patch")
	}
}

func TestNoteRepo_Delete(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	n := note.NewNote("To Delete", "content")
	if err := repo.Create(n); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(n.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	if err := repo.CheckByID(n.ID); err == nil {
		t.Error("expected error after delete")
	}
}

func TestNoteRepo_Search(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	if err := repo.Create(note.NewNote("Meeting Notes", "Discuss project roadmap")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if err := repo.Create(note.NewNote("Shopping List", "Milk and eggs")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	results, err := repo.Search("meeting")
	if err != nil {
		t.Fatalf("Search failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	if results[0].Title != "Meeting Notes" {
		t.Errorf("expected 'Meeting Notes', got '%s'", results[0].Title)
	}

	results, err = repo.Search("roadmap")
	if err != nil {
		t.Fatalf("Search 'roadmap' failed: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected 1 result for 'roadmap', got %d", len(results))
	}

	results, err = repo.Search("nonexistent")
	if err != nil {
		t.Fatalf("Search 'nonexistent' failed: %v", err)
	}
	if len(results) != 0 {
		t.Errorf("expected 0 results, got %d", len(results))
	}

	// FTS4 special chars should not error
	results, err = repo.Search("C++")
	if err != nil {
		t.Fatalf("Search with special chars failed: %v", err)
	}

	results, err = repo.Search("test*query")
	if err != nil {
		t.Fatalf("Search with wildcard chars failed: %v", err)
	}
}

func TestNoteRepo_GetRecent(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	for i := 0; i < 5; i++ {
		n := note.NewNote("Note", "content")
		time.Sleep(1 * time.Millisecond)
		if err := repo.Create(n); err != nil {
			t.Fatalf("Create failed: %v", err)
		}
	}

	recent, err := repo.GetRecent(3)
	if err != nil {
		t.Fatalf("GetRecent failed: %v", err)
	}
	if len(recent) != 3 {
		t.Errorf("expected 3 recent notes, got %d", len(recent))
	}
}

func TestNoteRepo_Tags(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	noteRepo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}
	tagRepo, err := NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	n := note.NewNote("Tagged Note", "content")
	if err := noteRepo.Create(n); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	work := tag.NewTag("work")
	if err := tagRepo.Create(work); err != nil {
		t.Fatalf("Create tag failed: %v", err)
	}

	if err := noteRepo.AddTagToNote(n.ID, work.ID); err != nil {
		t.Fatalf("AddTagToNote failed: %v", err)
	}

	// Duplicate should be ignored (INSERT OR IGNORE)
	if err := noteRepo.AddTagToNote(n.ID, work.ID); err != nil {
		t.Fatalf("AddTagToNote duplicate failed: %v", err)
	}

	tags, err := noteRepo.GetTagsByNote(n.ID)
	if err != nil {
		t.Fatalf("GetTagsByNote failed: %v", err)
	}
	if len(tags) != 1 {
		t.Fatalf("expected 1 tag, got %d", len(tags))
	}
	if tags[0].Name != "work" {
		t.Errorf("expected tag 'work', got '%s'", tags[0].Name)
	}

	// GetByID should include tags
	got, err := noteRepo.GetByID(n.ID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if len(got.Tags) != 1 || got.Tags[0] != "work" {
		t.Errorf("expected tags ['work'], got %v", got.Tags)
	}

	// Remove tags
	if err := noteRepo.RemoveTagFromNote(n.ID); err != nil {
		t.Fatalf("RemoveTagFromNote failed: %v", err)
	}
	tags, err = noteRepo.GetTagsByNote(n.ID)
	if err != nil {
		t.Fatalf("GetTagsByNote after remove failed: %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("expected 0 tags after remove, got %d", len(tags))
	}

	// Cascade: delete note should clean up notes_tags
	if err := noteRepo.Delete(n.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}
}

func TestNoteRepo_ExportJSON(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	if err := repo.Create(note.NewNote("Export Test", "some content")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	dir := t.TempDir()
	if err := repo.ExportNotes(dir, nil, "json"); err != nil {
		t.Fatalf("ExportNotes json failed: %v", err)
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(files) == 0 {
		t.Error("expected at least 1 exported file")
	}
}

func TestNoteRepo_ExportMarkdown(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	if err := repo.Create(note.NewNote("MD Export", "# heading\ncontent")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	dir := t.TempDir()
	if err := repo.ExportNotes(dir, nil, "markdown"); err != nil {
		t.Fatalf("ExportNotes markdown failed: %v", err)
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir failed: %v", err)
	}
	if len(entries) == 0 {
		t.Error("expected at least 1 exported file")
	}

	// Verify file contents
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".md" {
			t.Errorf("expected .md extension, got '%s'", e.Name())
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			t.Errorf("ReadFile failed: %v", err)
			continue
		}
		if len(data) == 0 {
			t.Error("exported file is empty")
		}
	}
}

func TestNoteRepo_ExportInvalidFormat(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	if err := repo.Create(note.NewNote("Test", "content")); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	dir := t.TempDir()
	if err := repo.ExportNotes(dir, nil, "xml"); err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestNoteRepo_ExportSince(t *testing.T) {
	db := setupNoteDB(t)
	defer db.Close()

	repo, err := NewNoteRepository(db)
	if err != nil {
		t.Fatalf("NewNoteRepository failed: %v", err)
	}

	old := note.NewNote("Old Note", "old")
	old.CreatedAt = time.Now().Add(-48 * time.Hour)
	old.UpdatedAt = old.CreatedAt
	if err := repo.Create(old); err != nil {
		t.Fatalf("Create old note failed: %v", err)
	}

	recent := note.NewNote("Recent Note", "recent")
	if err := repo.Create(recent); err != nil {
		t.Fatalf("Create recent note failed: %v", err)
	}

	since := time.Now().Add(-24 * time.Hour)
	dir := t.TempDir()
	if err := repo.ExportNotes(dir, &since, "json"); err != nil {
		t.Fatalf("ExportNotes since failed: %v", err)
	}

	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Errorf("expected 1 file (only recent), got %d", len(files))
	}
}
