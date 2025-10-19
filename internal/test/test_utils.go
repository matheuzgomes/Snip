package test

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/snip/internal/handler"
	"github.com/snip/internal/note"
	"github.com/snip/internal/tag"
)

type mockNoteRepository struct {
	notes         []*note.Note
	notesWithTags []*note.NoteWithTags
	err           error
}

func (m *mockNoteRepository) Create(note *note.Note) error {
	if m.err != nil {
		return m.err
	}
	note.ID = len(m.notes) + 1
	m.notes = append(m.notes, note)
	return nil
}

func (m *mockNoteRepository) GetByID(id int) (*note.NoteWithTags, error) {
	if m.err != nil {
		return nil, m.err
	}

	for _, note := range m.notesWithTags {
		if note.ID == id {
			return note, nil
		}
	}
	return nil, ErrNoteNotFound
}

func (m *mockNoteRepository) GetAll(isAsc bool, tagID int) ([]*note.NoteWithTags, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.notesWithTags, nil
}

func (m *mockNoteRepository) Update(id int, content string, title string) error {
	if m.err != nil {
		return m.err
	}

	for _, note := range m.notesWithTags {
		if note.ID == id {
			note.Title = title
			note.Content = content
			note.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrNoteNotFound
}

func (m *mockNoteRepository) Delete(id int) error {
	if m.err != nil {
		return m.err
	}

	for i, note := range m.notesWithTags {
		if note.ID == id {
			m.notesWithTags = append(m.notesWithTags[:i], m.notesWithTags[i+1:]...)
			return nil
		}
	}
	return ErrNoteNotFound
}

func (m *mockNoteRepository) Search(term string) ([]*note.Note, error) {
	if m.err != nil {
		return nil, m.err
	}

	var results []*note.Note
	for _, noteWithTags := range m.notesWithTags {
		if strings.Contains(strings.ToLower(noteWithTags.Title), strings.ToLower(term)) ||
			strings.Contains(strings.ToLower(noteWithTags.Content), strings.ToLower(term)) {
			results = append(results, &note.Note{
				ID:        noteWithTags.ID,
				Title:     noteWithTags.Title,
				Content:   noteWithTags.Content,
				CreatedAt: noteWithTags.CreatedAt,
				UpdatedAt: noteWithTags.UpdatedAt,
			})
		}
	}
	return results, nil
}

func (m *mockNoteRepository) CheckByID(id int) error {
	if m.err != nil {
		return m.err
	}

	for _, note := range m.notesWithTags {
		if note.ID == id {
			return nil
		}
	}
	return ErrNoteNotFound
}

func (m *mockNoteRepository) Patch(id int, title string) error {
	if m.err != nil {
		return m.err
	}

	for _, note := range m.notesWithTags {
		if note.ID == id {
			note.Title = title
			note.UpdatedAt = time.Now()
			return nil
		}
	}
	return ErrNoteNotFound
}

func (m *mockNoteRepository) GetRecent(limit int) ([]*note.NoteWithTags, error) {
	if m.err != nil {
		return nil, m.err
	}

	if limit <= 0 || limit > len(m.notesWithTags) {
		limit = len(m.notesWithTags)
	}

	start := len(m.notesWithTags) - limit
	if start < 0 {
		start = 0
	}

	return m.notesWithTags[start:], nil
}

func (m *mockNoteRepository) ExportNotes(exportDir string, since *time.Time, format string) error {
	if m.err != nil {
		return m.err
	}

	if format != "json" && format != "markdown" {
		return fmt.Errorf("invalid format: %s", format)
	}

	return nil
}

func (m *mockNoteRepository) AddTagToNote(noteID, tagID int) error {
	return nil
}

func (m *mockNoteRepository) RemoveTagFromNote(noteID int) error {
	return nil
}

func (m *mockNoteRepository) GetTagsByNote(noteID int) ([]*tag.Tag, error) {
	return nil, nil
}

func (m *mockNoteRepository) Close() error {
	return nil
}

type mockTagRepository struct {
	err error
}

func (m *mockTagRepository) Create(tag *tag.Tag) error {
	return m.err
}

func (m *mockTagRepository) GetByName(name string) (*tag.Tag, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &tag.Tag{
		ID:   1,
		Name: name,
	}, nil
}

func (m *mockTagRepository) GetAll() ([]*tag.Tag, error) {
	if m.err != nil {
		return nil, m.err
	}
	return []*tag.Tag{
		{ID: 1, Name: "work"},
		{ID: 2, Name: "personal"},
		{ID: 3, Name: "important"},
		{ID: 4, Name: "meeting"},
	}, nil
}

func (m *mockTagRepository) Delete(id int) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockTagRepository) Update(id int, name string) error {
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockTagRepository) GetOrCreate(name string) (*tag.Tag, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &tag.Tag{
		ID:   1,
		Name: name,
	}, nil
}

func (m *mockTagRepository) Close() error {
	return nil
}

func createTestHandler() (handler.Handler, *mockNoteRepository, *mockTagRepository) {
	mockNoteRepo := &mockNoteRepository{}
	mockTagRepo := &mockTagRepository{}

	h := handler.NewHandler(mockNoteRepo, mockTagRepo)
	return h, mockNoteRepo, mockTagRepo
}

func stringPtr(s string) *string {
	return &s
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr ||
		(len(s) > len(substr) &&
			(s[:len(substr)] == substr ||
				s[len(s)-len(substr):] == substr ||
				containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// Common test errors
var (
	ErrDatabaseConnection = errors.New("database connection failed")
	ErrValidationFailed   = errors.New("validation failed")
	ErrNoteNotFound       = errors.New("note not found")
	ErrTagNotFound        = errors.New("no note found for this tag")
)

// Helper functions to create test data
func createTestNotes() []*note.NoteWithTags {
	now := time.Now()
	return []*note.NoteWithTags{
		{
			ID:        1,
			Title:     "First Note",
			Content:   "This is the first note content",
			CreatedAt: now.Add(-2 * time.Hour),
			UpdatedAt: now.Add(-1 * time.Hour),
			Tags:      []string{"work", "important"},
		},
		{
			ID:        2,
			Title:     "Second Note",
			Content:   "This is the second note content",
			CreatedAt: now.Add(-1 * time.Hour),
			UpdatedAt: now.Add(-30 * time.Minute),
			Tags:      []string{"personal"},
		},
		{
			ID:        3,
			Title:     "Third Note",
			Content:   "This is the third note content",
			CreatedAt: now,
			UpdatedAt: now,
			Tags:      []string{"work", "meeting"},
		},
	}
}

func createTestNotesWithTag(tagName string) []*note.NoteWithTags {
	notes := createTestNotes()
	var filtered []*note.NoteWithTags

	for _, note := range notes {
		for _, tag := range note.Tags {
			if tag == tagName {
				filtered = append(filtered, note)
				break
			}
		}
	}

	return filtered
}
