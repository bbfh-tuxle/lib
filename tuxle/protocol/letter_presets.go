package protocol

import (
	"fmt"

	"github.com/bbfh-tuxle/lib/tuxle/fields"
)

// Get a new error letter preset.
func NewErrorLetter(name string, body string, args ...any) *Letter {
	return &Letter{
		Type:       "ERROR",
		Endpoint:   name,
		Parameters: fields.Parameters{},
		Body:       fmt.Sprintf(body, args...),
	}
}

// Get a new success (okay) letter preset.
func NewSuccessLetter(name string) *Letter {
	return &Letter{
		Type:       "OKAY",
		Endpoint:   name,
		Parameters: fields.Parameters{},
		Body:       "",
	}
}
