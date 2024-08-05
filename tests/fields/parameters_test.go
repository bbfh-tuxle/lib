package fields_test

import (
	"testing"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/fields"
	"gotest.tools/assert"
)

func TestParametersParse(t *testing.T) {
	params, err := fields.ReadAllParameters(stream.NewReader("Key=Value\n"))
	if err != nil {
		t.Fatal(err)
	}

	expected := fields.Parameters{}
	expected["Key"] = "Value"

	assert.DeepEqual(t, params, expected)
}

func TestParametersParseMultiple(t *testing.T) {
	params, err := fields.ReadAllParameters(stream.NewReader("Key=Value\nKeyN=ValueN\n"))
	if err != nil {
		t.Fatal(err)
	}

	expected := fields.Parameters{}
	expected["Key"] = "Value"
	expected["KeyN"] = "ValueN"

	assert.DeepEqual(t, params, expected)
}

func TestParametersParseFixed(t *testing.T) {
	params, err := fields.ReadParameters(stream.NewReader("Key=Value\nKeyN=ValueN\n"), 1)
	if err != nil {
		t.Fatal(err)
	}

	expected := fields.Parameters{}
	expected["Key"] = "Value"

	assert.DeepEqual(t, params, expected)
}
