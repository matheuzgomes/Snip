package test

import (
	"testing"
)

func TestListNotes(t *testing.T) {
	tests := []struct {
		name        string
		isAsc       bool
		verbose     bool
		tag         *string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful list notes ascending",
			isAsc:   true,
			verbose: false,
			tag:     nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:    "successful list notes descending",
			isAsc:   false,
			verbose: true,
			tag:     nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:    "successful list notes with tag filter",
			isAsc:   true,
			verbose: false,
			tag:     stringPtr("work"),
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
				noteRepo.notesWithTags = createTestNotesWithTag("work")
			},
			expectError: false,
		},
		{
			name:    "repository error",
			isAsc:   true,
			verbose: false,
			tag:     nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "failed to fetch notes",
		},
		{
			name:    "tags repository error",
			isAsc:   true,
			verbose: false,
			tag:     stringPtr("test-tag"),
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				tagRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "no note found for this tag",
		},
		{
			name:    "empty notes list",
			isAsc:   true,
			verbose: false,
			tag:     nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
			},
			expectError: false,
			errorMsg:    "No notes found.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.ListNotes(tt.isAsc, tt.verbose, tt.tag)

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

func TestListNotes_EdgeCases(t *testing.T) {
	t.Run("empty notes list", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil

		err := h.ListNotes(true, false, nil)

		if err != nil {
			t.Errorf("Expected no error for empty list, got: %v", err)
		}
	})

	t.Run("invalid tag filter", func(t *testing.T) {
		h, mockNoteRepo, mockTagRepo := createTestHandler()
		mockNoteRepo.err = nil
		mockTagRepo.err = ErrNoteNotFound

		err := h.ListNotes(true, false, stringPtr("nonexistent-tag"))

		if err == nil {
			t.Errorf("Expected error for invalid tag, got none")
		}
	})
}

func BenchmarkListNotes(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.ListNotes(true, false, nil)
		if err != nil {
			b.Fatalf("ListNotes failed: %v", err)
		}
	}
}
