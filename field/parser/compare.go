package parser

import (
	"reflect"
	"testing"
)

func Assert[T any](t *testing.T, got T, expected T) {
	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Expected `%v`, Got `%v`", expected, got)
	}
}
