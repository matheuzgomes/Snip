package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func GetDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dbDir := filepath.Join(homeDir, ".snip")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dbDir, "notes.db"), nil
}

func Connect() (*sql.DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := ensureDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureDatabase(db *sql.DB) error {
	query := `
    -- Main Table
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );


    CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS notes_tags (
        note_id INTEGER NOT NULL,
        tag_id INTEGER NOT NULL,
        PRIMARY KEY (note_id, tag_id),
        FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
    );

    -- Index
    CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
    CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);
    
    -- FTS Table
    CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts4(id, title, content);
    
    -- Populate FTS table with existing data (only if empty)
    INSERT OR IGNORE INTO notes_fts(id, title, content) 
    SELECT id, title, content FROM notes 
    WHERE id NOT IN (SELECT id FROM notes_fts);
    
    -- Triggers
    CREATE TRIGGER IF NOT EXISTS notes_fts_ai AFTER INSERT ON notes BEGIN
        INSERT INTO notes_fts(id, title, content) VALUES (new.id, new.title, new.content);
    END;
    
    CREATE TRIGGER IF NOT EXISTS notes_fts_au AFTER UPDATE ON notes BEGIN
        UPDATE notes_fts SET title = new.title, content = new.content WHERE id = old.id;
    END;
    
    CREATE TRIGGER IF NOT EXISTS notes_fts_ad AFTER DELETE ON notes BEGIN
        DELETE FROM notes_fts WHERE id = old.id;
    END;
    `

	_, err := db.Exec(query)
	return err
}
