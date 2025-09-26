package validation

import (
	"fmt"
	"strings"
)

type Validator struct {
	Field   string
	Message string
}

func (e *Validator) Error() string {
	return fmt.Sprintf("%s %s", e.Field, e.Message)
}

func NewValidator() *Validator {
	return &Validator{}
}

func (v *Validator) ValidateNote(title string) error {
	if strings.TrimSpace(title) == "" {
		return &Validator{
			Field:   "title",
			Message: "is required",
		}
	}

	return nil
}
