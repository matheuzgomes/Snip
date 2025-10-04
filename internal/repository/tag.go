package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/snip/internal/tag"
)

var ErrTagNotFound = errors.New("tag not found")

type TagRepository interface {
	Create(tag *tag.Tag) error
	GetByName(name string) (*tag.Tag, error)
	GetAll() ([]*tag.Tag, error)
	Delete(id int) error
	GetOrCreate(name string) (*tag.Tag, error)
	Close() error
}

type tagRepository struct {
	db *sql.DB
}

func NewTagRepository(db *sql.DB) (TagRepository, error) {
	return &tagRepository{db: db}, nil
}

func (r *tagRepository) Close() error {
	return r.db.Close()
}

func (r *tagRepository) Create(tag *tag.Tag) error {
	query := `INSERT INTO tags (name) VALUES (?)`

	result, err := r.db.Exec(query, tag.Name)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tag.ID = int(id)
	return nil
}

func (r *tagRepository) GetByName(name string) (*tag.Tag, error) {
	query := `SELECT id, name FROM tags WHERE name = ?`

	tag := &tag.Tag{}
	err := r.db.QueryRow(query, name).Scan(&tag.ID, &tag.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTagNotFound
		}
		return nil, err
	}

	return tag, nil
}

func (r *tagRepository) GetAll() ([]*tag.Tag, error) {
	query := `SELECT id, name FROM tags ORDER BY name`

	rows, err := r.db.Query(query)
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

func (r *tagRepository) Delete(id int) error {
	query := `DELETE FROM tags WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *tagRepository) GetOrCreate(name string) (*tag.Tag, error) {
	retrievedTag, err := r.GetByName(name)
	if err == nil {
		return retrievedTag, nil
	}

	if errors.Is(err, ErrTagNotFound) {
		newTag := tag.NewTag(name)
		if err := r.Create(newTag); err != nil {
			return nil, fmt.Errorf("failed to create tag: %w", err)
		}
		return newTag, nil
	}

	return nil, err
}
