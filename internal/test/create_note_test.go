package test

import (
	"errors"
	"testing"
)

func TestCreateNote(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		message     *string
		tag         *string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful note creation with message",
			title: "Test Note",
			message: stringPtr("Test content"),
			tag: nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
			},
			expectError: false,
		},
		{
			name:  "successful note creation with tag",
			title: "Test Note",
			message: stringPtr("Test content"),
			tag: stringPtr("test-tag"),
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
			},
			expectError: false,
		},
		{
			name:        "empty title validation error",
			title:       "",
			message:     nil,
			tag:         nil,
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg:    "title is required",
		},
		{
			name:        "whitespace only title validation error",
			title:       "   ",
			message:     nil,
			tag:         nil,
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg: "title is required",
		},
		{
			name:  "repository create error",
			title: "Test Note",
			message: stringPtr("Test content"),
			tag: nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = errors.New("database connection failed")
			},
			expectError: true,
			errorMsg: "failed to create note",
		},
		{
			name:  "tag repository associate tag with note error",
			title: "Test Note",
			message: stringPtr("Test Tag"),
			tag: stringPtr("test-tag"),
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				tagRepo.err = errors.New("failed to associate tags with note:")
			},
			expectError: true,
			errorMsg: "failed to associate tags with note:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.CreateNote(tt.title, tt.message, tt.tag)

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
					return
				}

				if len(mockNoteRepo.notes) == 0 {
					t.Errorf("Expected note to be created but none found")
					return
				}

				createdNote := mockNoteRepo.notes[0]
				if createdNote.Title != tt.title {
					t.Errorf("Expected title '%s', got '%s'", tt.title, createdNote.Title)
				}
			}
		})
	}
}

func TestCreateNote_EdgeCases(t *testing.T) {
	t.Run("very long title", func(t *testing.T) {
		h, mockNoteRepo, mockTagRepo := createTestHandler()
		mockNoteRepo.err = nil
		mockTagRepo.err = nil

		longTitle := "This is a very long title that might cause issues in some systems but should still be valid for our note creation"
		message := "Test content"

		err := h.CreateNote(longTitle, &message, nil)

		if err != nil {
			t.Errorf("Expected no error for long title, got: %v", err)
		}

		if len(mockNoteRepo.notes) == 0 {
			t.Errorf("Expected note to be created")
		}
	})

	t.Run("special characters in title", func(t *testing.T) {
		h, mockNoteRepo, mockTagRepo := createTestHandler()
		mockNoteRepo.err = nil
		mockTagRepo.err = nil

		specialTitle := "Note with special chars: @#$%^&*()_+-=[]{}|;':\",./<>?"
		message := "Test content"

		err := h.CreateNote(specialTitle, &message, nil)

		if err != nil {
			t.Errorf("Expected no error for special characters, got: %v", err)
		}
	})

}

func BenchmarkCreateNote(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil

	title := "Benchmark Note"
	message := "Benchmark content"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.CreateNote(title, &message, nil)
		if err != nil {
			b.Fatalf("CreateNote failed: %v", err)
		}
		mockNoteRepo.notes = nil
	}
}
