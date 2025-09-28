package note

import (
	"database/sql"
	"errors"
	"time"

	"github.com/snip/internal/database"
)

type Repository interface {
	Create(note *Note) error
	GetByID(id int) (*Note, error)
	GetAll(isAsc bool) ([]*Note, error)
	Update(id int, content string, title string) error
	Delete(id int) error
	Search(term string) ([]*Note, error)
	CheckByID(id int) error
	Close() error
}

type repository struct {
	db *sql.DB
}

func NewRepository() (Repository, error) {
	db, err := database.Connect()
	if err != nil {
		return nil, err
	}

	return &repository{db: db}, nil
}

func (r *repository) Close() error {
	return r.db.Close()
}

func (r *repository) Create(note *Note) error {
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

func (r *repository) GetByID(id int) (*Note, error) {
	query := `
		SELECT id, title, content, created_at, updated_at
		FROM notes WHERE id = ?
	`

	note := &Note{}

	err := r.db.QueryRow(query, id).Scan(
		&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt,
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

func (r *repository) GetAll(isAsc bool) ([]*Note, error) {

	orderBy := "DESC"

	if isAsc {
		orderBy = "ASC"
	}

	query := `
		SELECT id, title, content, created_at, updated_at
		FROM notes ORDER BY created_at ` + orderBy

	db, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var notes []*Note
	for db.Next() {
		note := &Note{}
		err := db.Scan(&note.ID, &note.Title, &note.Content, &note.CreatedAt, &note.UpdatedAt)
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

func (r *repository) Search(term string) ([]*Note, error) {
	query := `
		SELECT id, title, content
		FROM notes_fts 
		WHERE notes_fts MATCH ?
	`

	db, err := r.db.Query(query, term)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var notes []*Note
	for db.Next() {
		note := &Note{}
		err := db.Scan(&note.ID, &note.Title, &note.Content)
		if err != nil {
			return nil, err
		}
		notes = append(notes, note)
	}

	return notes, nil
}
