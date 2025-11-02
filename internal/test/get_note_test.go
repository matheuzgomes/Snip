package test

import (
	"testing"
)

func TestGetNote(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		verbose     bool
		render      bool
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:    "successful get note with verbose",
			idStr:   "1",
			verbose: true,
			render: false,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:    "successful get note without verbose",
			idStr:   "2",
			verbose: false,
			render: false,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:        "invalid id format",
			idStr:       "invalid",
			verbose:     false,
			render: false,
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg:    "invalid note ID",
		},
		{
			name:    "note not found",
			idStr:   "999",
			verbose: false,
			render: false,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "note not found",
		},
		{
			name:    "repository error",
			idStr:   "1",
			verbose: false,
			render: false,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "failed to fetch note",
		},
		{
			name:    "successful get note with render",
			idStr:   "1",
			verbose: false,
			render: true,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.GetNote(tt.idStr, tt.verbose, tt.render)

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

func TestGetNote_EdgeCases(t *testing.T) {
	t.Run("negative id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.GetNote("-1", false, false)

		if err == nil {
			t.Errorf("Expected error for negative ID, got none")
		}
	})

	t.Run("zero id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.GetNote("0", false, false)

		if err == nil {
			t.Errorf("Expected error for zero ID, got none")
		}
	})
}

func BenchmarkGetNote(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.GetNote("1", false, false)
		if err != nil {
			b.Fatalf("GetNote failed: %v", err)
		}
	}
}
