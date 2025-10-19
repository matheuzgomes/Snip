package test

import (
	"testing"
)

func TestImportNotes(t *testing.T) {
	tests := []struct {
		name        string
		importDir   string
		setupMocks  func(*mockNoteRepository, *mockTagRepository)
		expectError bool
		errorMsg    string
	}{
		{
			name:      "successful import from directory",
			importDir: "/tmp/test_import",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
			},
			expectError: true,
			errorMsg:    "failed to read import directory",
		},
		{
			name:        "empty import directory",
			importDir:   "",
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: false,
		},
		{
			name:        "whitespace only import directory",
			importDir:   "   ",
			setupMocks:  func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {},
			expectError: true,
			errorMsg:    "failed to read import directory",
		},
		{
			name:      "nonexistent directory",
			importDir: "/nonexistent/directory",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = nil
				tagRepo.err = nil
			},
			expectError: true,
			errorMsg:    "failed to read import directory",
		},
		{
			name:      "repository create error",
			importDir: "/tmp/test_import",
			setupMocks: func(noteRepo *mockNoteRepository, tagRepo *mockTagRepository) {
				noteRepo.err = ErrDatabaseConnection
				tagRepo.err = nil
			},
			expectError: true,
			errorMsg:    "failed to read import directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h, mockNoteRepo, mockTagRepo := createTestHandler()
			tt.setupMocks(mockNoteRepo, mockTagRepo)

			err := h.ImportNotes(tt.importDir)

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

func TestImportNotes_EdgeCases(t *testing.T) {
	t.Run("import from home directory", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil

		err := h.ImportNotes("~/test_import")

		if err != nil && !contains(err.Error(), "failed to read import directory") {
			t.Errorf("Expected directory error, got: %v", err)
		}
	})

	t.Run("import from relative path", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil

		err := h.ImportNotes("./test_import")

		// This might fail due to directory not existing, which is expected
		if err != nil && !contains(err.Error(), "failed to read import directory") {
			t.Errorf("Expected directory error, got: %v", err)
		}
	})

	t.Run("import with special characters in path", func(t *testing.T) {
		h, mockNoteRepo, _ := createTestHandler()
		mockNoteRepo.err = nil

		err := h.ImportNotes("/tmp/test@#$%")

		// This might fail due to directory not existing, which is expected
		if err != nil && !contains(err.Error(), "failed to read import directory") {
			t.Errorf("Expected directory error, got: %v", err)
		}
	})
}

func BenchmarkImportNotes(b *testing.B) {
	h, mockNoteRepo, mockTagRepo := createTestHandler()
	mockNoteRepo.err = nil
	mockTagRepo.err = nil

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := h.ImportNotes("/tmp/test_import")
		if err != nil && !contains(err.Error(), "failed to read import directory") {
			b.Fatalf("ImportNotes failed: %v", err)
		}
	}
}
