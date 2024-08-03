package parser

import (
	"reflect"
	"testing"
)

// Used for testing. Asserts that the values are equals, and Fatals otherwise.
func Assert[T any](t *testing.T, got T, expected T) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected `%v`, Got `%v`", expected, got)
	}
}
