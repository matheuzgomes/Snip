package handler

import (
	"github.com/snip/internal/note"
	"github.com/snip/internal/validation"
)

type Handler interface {
	CreateNote(title string, message *string) error
	ListNotes(isAsc, verbose bool) error
	GetNote(idStr string, verbose bool) error
	FindNotes(term string) error
	UpdateNote(idStr string, title string) error
	DeleteNote(idStr string) error
}

type handler struct {
	noteRepo      note.Repository
	validator     *validation.Validator
	editorHandler *EditorHandler
	dateFormat    string
}

func NewHandler(repo note.Repository) Handler {
	return &handler{
		noteRepo:      repo,
		validator:     validation.NewValidator(),
		dateFormat:    "2006-01-02 15:04:05",
		editorHandler: NewEditorHandler(),
	}
}
