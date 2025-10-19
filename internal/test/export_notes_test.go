package test

import (
	"testing"

	"github.com/snip/internal/note"
)

func TestExportNotes(t *testing.T) {
	tests := []struct {
		name        string
		since       string
		format      string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:   "successful export JSON format",
			since:  "",
			format: "json",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:   "successful export Markdown format",
			since:  "",
			format: "markdown",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:   "successful export with since date",
			since:  "2024-01-01",
			format: "json",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:   "invalid format",
			since:  "",
			format: "invalid",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "invalid format",
		},
		{
			name:   "invalid since date format",
			since:  "invalid-date",
			format: "json",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: true,
			errorMsg:    "invalid --since value",
		},
		{
			name:   "no notes to export",
			since:  "",
			format: "json",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = []*note.NoteWithTags{}
			},
			expectError: false,
		},
		{
			name:   "repository error",
			since:  "",
			format: "json",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
			},
			expectError: true,
			errorMsg:    "failed to export notes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.ExportNotes(tt.since, tt.format)

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

func TestExportNotes_EdgeCases(t *testing.T) {
	t.Run("export with future date", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.ExportNotes("2030-01-01", "json")

		if err != nil {
			t.Errorf("Expected no error for future date, got: %v", err)
		}
	})

	t.Run("export with very old date", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.ExportNotes("1900-01-01", "json")

		if err != nil {
			t.Errorf("Expected no error for old date, got: %v", err)
		}
	})

	t.Run("export with special characters in format", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.ExportNotes("", "json@#$")

		if err == nil {
			t.Errorf("Expected error for invalid format with special characters, got none")
		}
	})
}

func BenchmarkExportNotes(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.ExportNotes("", "json")
		if err != nil {
			b.Fatalf("ExportNotes failed: %v", err)
		}
	}
}
