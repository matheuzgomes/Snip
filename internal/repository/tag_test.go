package repository

import (
	"database/sql"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/matheuzgomes/Snip/internal/tag"
)

func setupTagDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}

	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("failed to enable foreign keys: %v", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		)
	`)
	if err != nil {
		t.Fatalf("failed to create tags table: %v", err)
	}

	return db
}

func TestTagRepo_Create(t *testing.T) {
	db := setupTagDB(t)
	defer db.Close()

	repo, err := NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	tg := tag.NewTag("work")
	if err := repo.Create(tg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if tg.ID == 0 {
		t.Error("expected non-zero ID after create")
	}

	// Duplicate should fail due to UNIQUE constraint
	tg2 := tag.NewTag("work")
	if err := repo.Create(tg2); err == nil {
		t.Error("expected error on duplicate tag name")
	}
}

func TestTagRepo_GetByName(t *testing.T) {
	db := setupTagDB(t)
	defer db.Close()

	repo, err := NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	tg := tag.NewTag("personal")
	if err := repo.Create(tg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	found, err := repo.GetByName("personal")
	if err != nil {
		t.Fatalf("GetByName failed: %v", err)
	}
	if found.Name != "personal" {
		t.Errorf("expected name 'personal', got '%s'", found.Name)
	}

	_, err = repo.GetByName("nonexistent")
	if err != ErrTagNotFound {
		t.Errorf("expected ErrTagNotFound, got %v", err)
	}
}

func TestTagRepo_GetOrCreate(t *testing.T) {
	db := setupTagDB(t)
	defer db.Close()

	repo, err := NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	// Create new
	tg1, err := repo.GetOrCreate("meeting")
	if err != nil {
		t.Fatalf("GetOrCreate failed: %v", err)
	}
	if tg1.Name != "meeting" {
		t.Errorf("expected 'meeting', got '%s'", tg1.Name)
	}

	// Get existing
	tg2, err := repo.GetOrCreate("meeting")
	if err != nil {
		t.Fatalf("GetOrCreate (existing) failed: %v", err)
	}
	if tg2.ID != tg1.ID {
		t.Errorf("expected same ID %d, got %d", tg1.ID, tg2.ID)
	}
}

func TestTagRepo_GetAll(t *testing.T) {
	db := setupTagDB(t)
	defer db.Close()

	repo, err := NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	for _, name := range []string{"zulu", "alpha", "bravo"} {
		if err := repo.Create(tag.NewTag(name)); err != nil {
			t.Fatalf("Create(%q) failed: %v", name, err)
		}
	}

	tags, err := repo.GetAll()
	if err != nil {
		t.Fatalf("GetAll failed: %v", err)
	}

	if len(tags) != 3 {
		t.Errorf("expected 3 tags, got %d", len(tags))
	}

	// Should be ordered by name
	if tags[0].Name != "alpha" {
		t.Errorf("expected 'alpha' first, got '%s'", tags[0].Name)
	}
}

func TestTagRepo_Delete(t *testing.T) {
	db := setupTagDB(t)
	defer db.Close()

	repo, err := NewTagRepository(db)
	if err != nil {
		t.Fatalf("NewTagRepository failed: %v", err)
	}

	tg := tag.NewTag("temp")
	if err := repo.Create(tg); err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(tg.ID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	_, err = repo.GetByName("temp")
	if err != ErrTagNotFound {
		t.Errorf("expected ErrTagNotFound after delete, got %v", err)
	}
}
