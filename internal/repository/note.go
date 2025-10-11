package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/snip/internal/note"
	"github.com/snip/internal/tag"
)

type NoteRepository interface {
	Create(note *note.Note) error
	GetByID(id int) (*note.NoteWithTags, error)
	GetAll(isAsc bool, tagID int) ([]*note.NoteWithTags, error)
	Update(id int, content string, title string) error
	Delete(id int) error
	Search(term string) ([]*note.Note, error)
	CheckByID(id int) error
	Patch(id int, title string) error
	GetRecent(limit int) ([]*note.NoteWithTags, error)

	// Tag operations
	AddTagToNote(noteID, tagID int) error
	RemoveTagFromNote(noteID int) error
	GetTagsByNote(noteID int) ([]*tag.Tag, error)

	Close() error
}

type repository struct {
	db *sql.DB
}

func NewNoteRepository(db *sql.DB) (NoteRepository, error) {
	return &repository{db: db}, nil
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) Create(note *note.Note) error {
	query := `
		INSERT INTO notes (title, content, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`

	result, err := r.db.Exec(query, note.Title, note.Content, note.CreatedAt, note.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	note.ID = int(id)
	return nil
}

func (r *repository) GetByID(id int) (*note.NoteWithTags, error) {
	query := `
		SELECT n.id, n.title, n.content, n.created_at, n.updated_at, GROUP_CONCAT(t.name) AS tags
		FROM notes n
		LEFT JOIN notes_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		WHERE n.id = ?
	`

	note := &note.NoteWithTags{}

	err := r.db.QueryRow(query, id).Scan(
		&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.Tags,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	return note, nil
}

func (r *repository) CheckByID(id int) error {
	query := `SELECT id FROM notes WHERE id = ?`

	if err := r.db.QueryRow(query, id).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return errors.New("not found")
		}
		return err
	}

	return nil
}

func (r *repository) GetAll(isAsc bool, tagID int) ([]*note.NoteWithTags, error) {

	orderBy := "DESC"

	if isAsc {
		orderBy = "ASC"
	}

	args := []any{}

	query := `
		SELECT n.id, n.title, n.content, n.created_at, n.updated_at, GROUP_CONCAT(t.name) AS tags
		FROM notes n
		LEFT JOIN notes_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		`


	if tagID != 0 {
		query += `WHERE nt.tag_id = ?`
		args = append(args, tagID)
	}

	query += ` GROUP BY n.id`

	query += ` ORDER BY n.created_at ` + orderBy

	db, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var notes []*note.NoteWithTags
	for db.Next() {
		note := &note.NoteWithTags{}
		err := db.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.Tags)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *repository) Update(id int, content string, title string) error {
	var query string
	var args []any

	query = `
        UPDATE notes
        SET content = ?, updated_at = ?
    `
	args = []any{content, time.Now()}

	if title != "" {
		query += `, title = ?`
		args = append(args, title)
	}

	query += ` WHERE id = ?`

	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	return err
}

func (r *repository) Delete(id int) error {
	query := `DELETE FROM notes WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *repository) Search(term string) ([]*note.Note, error) {
	query := `
		SELECT n.id, n.title, n.content
		FROM notes_fts n
		WHERE notes_fts MATCH ?
	`

	db, err := r.db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var notes []*note.Note
	for db.Next() {
		note := &note.Note{}
		err := db.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}

func (r *repository) AddTagToNote(noteID, tagID int) error {
	query := `INSERT OR IGNORE INTO notes_tags (note_id, tag_id) VALUES (?, ?)`
	_, err := r.db.Exec(query, noteID, tagID)
	return err
}

func (r *repository) RemoveTagFromNote(noteID int) error {
	query := `DELETE FROM notes_tags WHERE note_id = ?`
	_, err := r.db.Exec(query, noteID)
	return err
}

func (r *repository) GetTagsByNote(noteID int) ([]*tag.Tag, error) {
	query := `
		SELECT t.id, t.name
		FROM tags t
		INNER JOIN notes_tags nt ON t.id = nt.tag_id
		WHERE nt.note_id = ?
		ORDER BY t.name
	`

	rows, err := r.db.Query(query, noteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*tag.Tag
	for rows.Next() {
		tag := &tag.Tag{}
		err := rows.Scan(&tag.ID, &tag.Name)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (r *repository) Patch(id int, title string) error {
	query := `UPDATE notes SET title = ? WHERE id = ?`
	_, err := r.db.Exec(query, title, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetRecent(limit int) ([]*note.NoteWithTags, error) {
	query := `
		SELECT n.id, n.title, n.content, n.created_at, n.updated_at, GROUP_CONCAT(t.name) AS tags
		FROM notes n
		LEFT JOIN notes_tags nt ON n.id = nt.note_id
		LEFT JOIN tags t ON nt.tag_id = t.id
		GROUP BY n.id
		ORDER BY n.updated_at DESC
		LIMIT ?
	`

	notes := []*note.NoteWithTags{}

	db, err := r.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	for db.Next() {
		note := &note.NoteWithTags{}
		err := db.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt, &note.Tags)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}