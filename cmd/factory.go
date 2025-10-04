package cmd

import (
	"fmt"
	"sync"

	"github.com/snip/internal/database"
	"github.com/snip/internal/handler"
	"github.com/snip/internal/repository"
)

var (
	globalNoteRepo repository.NoteRepository
	globalTagRepo repository.TagRepository
	repoOnce   sync.Once
)

func getRepository() (repository.NoteRepository, repository.TagRepository, error) {
	var err error
	repoOnce.Do(func() {
		db, connectErr := database.Connect()
		if connectErr != nil {
			err = connectErr
			return
		}
		globalNoteRepo, err = repository.NewNoteRepository(db)
		globalTagRepo, err = repository.NewTagRepository(db)
	})
	return globalNoteRepo, globalTagRepo, err
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
