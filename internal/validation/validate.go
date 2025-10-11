package validation

import (
	"fmt"
	"strconv"
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


func (v *Validator) CheckString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (v *Validator) CheckInt(s string) int {
	num, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return num
}