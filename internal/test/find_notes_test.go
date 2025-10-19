package test

import (
	"testing"

	"github.com/snip/internal/note"
)

func TestFindNotes(t *testing.T) {
	tests := []struct {
		name        string
		term        string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name: "successful search by title",
			term: "First",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name: "successful search by content",
			term: "content",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name: "successful search case insensitive",
			term: "FIRST",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name: "no results found",
			term: "nonexistent",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name: "empty search term",
			term: "",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name: "whitespace only search term",
			term: "   ",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name: "repository error",
			term: "test",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "failed to search notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.FindNotes(tt.term)

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

func TestFindNotes_EdgeCases(t *testing.T) {
	t.Run("search with special characters", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.FindNotes("@#$%")

		if err != nil {
			t.Errorf("Expected no error for special characters, got: %v", err)
		}
	})

	t.Run("search with very long term", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		longTerm := "This is a very long search term that might cause issues in some systems but should still be valid for our search functionality"
		err := h.FindNotes(longTerm)

		if err != nil {
			t.Errorf("Expected no error for long search term, got: %v", err)
		}
	})

	t.Run("search in empty database", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = []*note.NoteWithTags{}

		err := h.FindNotes("test")

		if err != nil {
			t.Errorf("Expected no error for empty database, got: %v", err)
		}
	})
}

func BenchmarkFindNotes(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.FindNotes("test")
		if err != nil {
			b.Fatalf("FindNotes failed: %v", err)
		}
	}
}
