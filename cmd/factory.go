package cmd

import (
	"fmt"
	"sync"

	"github.com/snip/internal/handler"
	"github.com/snip/internal/note"
)

var (
	globalRepo note.Repository
	repoOnce   sync.Once
)

func getRepository() (note.Repository, error) {
	var err error
	repoOnce.Do(func() {
		globalRepo, err = note.NewRepository()
	})
	return globalRepo, err
}

func setupHandler() (handler.Handler, error) {
	repo, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewHandler(repo)

	return h, nil
}

func executeWithHandler(fn func(handler.Handler) error) error {
	h, err := setupHandler()
	if err != nil {
		return fmt.Errorf("failed to setup handler: %w", err)
	}

	return fn(h)
}
