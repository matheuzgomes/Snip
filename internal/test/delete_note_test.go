package test

import (
	"testing"

	"github.com/snip/internal/note"
)

func TestDeleteNote(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful delete note",
			idStr: "1",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:        "invalid id format",
			idStr:       "invalid",
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg:    "invalid note ID",
		},
		{
			name:  "note not found",
			idStr: "999",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "note not found",
		},
		{
			name:  "repository error",
			idStr: "1",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "this note does not exist",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.DeleteNote(tt.idStr)

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

func TestDeleteNote_EdgeCases(t *testing.T) {
	t.Run("negative id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.DeleteNote("-1")

		if err == nil {
			t.Errorf("Expected error for negative ID, got none")
		}
	})

	t.Run("zero id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.DeleteNote("0")

		if err == nil {
			t.Errorf("Expected error for zero ID, got none")
		}
	})

	t.Run("delete from empty list", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = []*note.NoteWithTags{}

		err := h.DeleteNote("1")

		if err == nil {
			t.Errorf("Expected error for deleting from empty list, got none")
		}
	})
}

func BenchmarkDeleteNote(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.DeleteNote("1")
		if err != nil {
			b.Fatalf("DeleteNote failed: %v", err)
		}
		// Reset the mock for next iteration
		mockNoteRepo.notesWithTags = createTestNotes()
	}
}
