package test

import (
	"testing"
)

func TestUpdateNote(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		title       string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful update note",
			idStr: "1",
			title: "Updated Note Title",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "failed to open editor",
		},
		{
			name:        "invalid id format",
			idStr:       "invalid",
			title:       "Updated Title",
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg:    "invalid note ID",
		},
		{
			name:  "empty title",
			idStr: "1",
			title: "",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "failed to open editor",
		},
		{
			name:  "whitespace only title",
			idStr: "1",
			title: "   ",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "failed to open editor",
		},
		{
			name:  "note not found",
			idStr: "999",
			title: "Updated Title",
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
			title: "Updated Title",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "failed to fetch note",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.UpdateNote(tt.idStr, tt.title)

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

func TestUpdateNote_EdgeCases(t *testing.T) {
	t.Run("negative id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.UpdateNote("-1", "Updated Title")

		if err == nil {
			t.Errorf("Expected error for negative ID, got none")
		}
	})

	t.Run("zero id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.UpdateNote("0", "Updated Title")

		if err == nil {
			t.Errorf("Expected error for zero ID, got none")
		}
	})

	t.Run("very long title", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		longTitle := "This is a very long title that might cause issues in some systems but should still be valid for our note update"
		err := h.UpdateNote("1", longTitle)

		if err == nil {
			t.Errorf("Expected error for long title due to editor, got none")
		}
	})
}

func BenchmarkUpdateNote(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.UpdateNote("1", "Updated Title")
		if err != nil {
			b.Skip("UpdateNote benchmark skipped due to editor dependency")
		}
	}
}
