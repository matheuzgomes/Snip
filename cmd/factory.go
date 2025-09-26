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

func setupHandler() (handler.Handler, func() error, error) {
	repo, err := getRepository()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewHandler(repo)

	return h, func() error { return nil }, nil
}

func executeWithHandler(fn func(handler.Handler) error) error {
	h, _, err := setupHandler()
	if err != nil {
		return fmt.Errorf("failed to setup handler: %w", err)
	}

	return fn(h)
}
