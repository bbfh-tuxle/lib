package field_test

import (
	"bufio"
	"strings"
	"testing"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
)

func TestParametersParse(t *testing.T) {
	params, err := field.ReadParameters(bufio.NewReader(strings.NewReader("Key=Value\n")), -1)
	if err != nil {
		t.Fatal(err)
	}

	expected := field.NewParameters()
	expected["Key"] = "Value"

	parser.Assert(t, params, expected)
}

func TestParametersParseMultiple(t *testing.T) {
	params, err := field.ReadParameters(
		bufio.NewReader(strings.NewReader("Key=Value\nKeyN=ValueN\n")),
		-1,
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := field.NewParameters()
	expected["Key"] = "Value"
	expected["KeyN"] = "ValueN"

	parser.Assert(t, params, expected)
}

func TestParametersParseFixed(t *testing.T) {
	params, err := field.ReadParameters(
		bufio.NewReader(strings.NewReader("Key=Value\nKeyN=ValueN\n")),
		1,
	)
	if err != nil {
		t.Fatal(err)
	}

	expected := field.NewParameters()
	expected["Key"] = "Value"

	parser.Assert(t, params, expected)
}
