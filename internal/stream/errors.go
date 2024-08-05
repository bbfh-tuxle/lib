package stream

import (
	"errors"
	"strings"
)

func CombineErrors(errs ...error) error {
	var items = make([]string, len(errs))

	for i, err := range errs {
		items[i] = err.Error()
	}

	return errors.New(strings.Join(items, "\n"))
}
