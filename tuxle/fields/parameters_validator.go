package fields

import (
	"fmt"
)

// Validates that the value of a key is of a valid value.
type ValidateFunc func(value string) bool

// Always returns 'true', which essentially just ensures that the field exists.
var Exists ValidateFunc = func(_ string) bool { return true }

// Validate fields in Parameters using custom validation function.
//
// []errors is an empty slice if all validations pass.
// Use fields.Exists to simply validate that the field exists.
func (params Parameters) Validate(functions map[string]ValidateFunc) []error {
	var items []error

	for key, fn := range functions {
		value, ok := params[key]
		if !ok {
			items = append(items, fmt.Errorf("Key %q isn't found in parameters.", key))
		} else if !fn(value) {
			items = append(items, fmt.Errorf("Value of %s (%q) has invalid format.", key, value))
		}
	}

	return items
}
