package test

import (
	"testing"
)

func TestPatchNote(t *testing.T) {
	tests := []struct {
		name        string
		idStr       string
		title       *string
		tag         *string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:  "successful patch title only",
			idStr: "1",
			title: stringPtr("Patched Title"),
			tag:   nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "successful patch tag only",
			idStr: "1",
			title: nil,
			tag:   stringPtr("new-tag"),
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "successful patch both title and tag",
			idStr: "1",
			title: stringPtr("Patched Title"),
			tag:   stringPtr("new-tag"),
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:        "invalid id format",
			idStr:       "invalid",
			title:       stringPtr("Patched Title"),
			tag:         nil,
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg:    "invalid note ID",
		},
		{
			name:  "empty title",
			idStr: "1",
			title: stringPtr(""),
			tag:   nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "whitespace only title",
			idStr: "1",
			title: stringPtr("   "),
			tag:   nil,
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				noteRepo.notesWithTags = createTestNotes()
			},
			expectError: false,
		},
		{
			name:  "note not found",
			idStr: "999",
			title: stringPtr("Patched Title"),
			tag:   nil,
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
			title: stringPtr("Patched Title"),
			tag:   nil,
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

			err := h.PatchNote(tt.idStr, tt.title, tt.tag)

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

func TestPatchNote_EdgeCases(t *testing.T) {
	t.Run("negative id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.PatchNote("-1", stringPtr("Patched Title"), nil)

		if err == nil {
			t.Errorf("Expected error for negative ID, got none")
		}
	})

	t.Run("zero id", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.PatchNote("0", stringPtr("Patched Title"), nil)

		if err == nil {
			t.Errorf("Expected error for zero ID, got none")
		}
	})

	t.Run("very long title", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		longTitle := "This is a very long title that might cause issues in some systems but should still be valid for our note patch"
		err := h.PatchNote("1", stringPtr(longTitle), nil)

		if err != nil {
			t.Errorf("Expected no error for long title, got: %v", err)
		}
	})

	t.Run("nil title and tag", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil
		mockNoteRepo.notesWithTags = createTestNotes()

		err := h.PatchNote("1", nil, nil)

		if err != nil {
			t.Errorf("Expected no error for nil title and tag, got: %v", err)
		}
	})
}

func BenchmarkPatchNote(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil
	mockNoteRepo.notesWithTags = createTestNotes()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.PatchNote("1", stringPtr("Patched Title"), stringPtr("new-tag"))
		if err != nil {
			b.Fatalf("PatchNote failed: %v", err)
		}
	}
}
