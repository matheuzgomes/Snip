package cmd

import (
	"fmt"

	"github.com/matheuzgomes/Snip/internal/database"
	"github.com/matheuzgomes/Snip/internal/handler"
	"github.com/matheuzgomes/Snip/internal/repository"
)

var (
	globalNoteRepo repository.NoteRepository
	globalTagRepo  repository.TagRepository
)

func getRepository() (repository.NoteRepository, repository.TagRepository, error) {
	if globalNoteRepo != nil && globalTagRepo != nil {
		return globalNoteRepo, globalTagRepo, nil
	}

	db, err := database.Connect()
	if err != nil {
		return nil, nil, err
	}

	globalNoteRepo, err = repository.NewNoteRepository(db)
	if err != nil {
		return nil, nil, err
	}

	globalTagRepo, err = repository.NewTagRepository(db)
	if err != nil {
		return nil, nil, err
	}

	return globalNoteRepo, globalTagRepo, nil
}

func setupHandler() (handler.Handler, error) {
	noteRepo, tagRepo, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewHandler(noteRepo, tagRepo)

	return h, nil
}

func executeWithHandler(fn func(handler.Handler) error) error {
	h, err := setupHandler()
	if err != nil {
		return fmt.Errorf("failed to setup handler: %w", err)
	}

	return fn(h)
}
