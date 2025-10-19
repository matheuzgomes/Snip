package test

import (
	"testing"
	"time"

	"github.com/snip/internal/note"
)

func TestGetRecentNotes(t *testing.T) {
	tests := []struct {
		name        string
		limit       int
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful get recent notes with limit",
			limit: 2,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "successful get recent notes with zero limit",
			limit: 0,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "successful get recent notes with negative limit",
			limit: -1,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "successful get recent notes with large limit",
			limit: 100,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "no notes found",
			limit: 5,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = []*note.NoteWithTags{}
			},
			expectError: false,
		},
		{
			name:  "repository error",
			limit: 5,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "failed to get recent notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.GetRecentNotes(tt.limit)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestGetRecentNotes_EdgeCases(t *testing.T) {
	t.Run("limit equals number of notes", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		notes := createTestNotes()
		mockNoteRepo.notesWithTags = notes

		err := h.GetRecentNotes(len(notes))

		if err != nil {
			t.Errorf("Expected no error for limit equals notes count, got: %v", err)
		}
	})

	t.Run("limit greater than number of notes", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		notes := createTestNotes()
		mockNoteRepo.notesWithTags = notes

		err := h.GetRecentNotes(len(notes) + 10)

		if err != nil {
			t.Errorf("Expected no error for limit greater than notes count, got: %v", err)
		}
	})

	t.Run("single note", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = []*note.NoteWithTags{
			{
				ID:        1,
				Title:     "Single Note",
				Content:   "Single note content",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
				Tags:      []string{"test"},
			},
		}

		err := h.GetRecentNotes(1)

		if err != nil {
			t.Errorf("Expected no error for single note, got: %v", err)
		}
	})
}

func BenchmarkGetRecentNotes(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.GetRecentNotes(5)
		if err != nil {
			b.Fatalf("GetRecentNotes failed: %v", err)
		}
	}
}
