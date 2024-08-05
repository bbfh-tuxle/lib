package fields_test

import (
	"testing"

	"github.com/bbfh-tuxle/lib/tuxle/fields"
	"gotest.tools/assert"
)

func assertRandomStringLength(t *testing.T, length int) {
	assert.Equal(t, len(fields.RandomString(length)), length)
}

func TestEncodeString(t *testing.T) {
	assertRandomStringLength(t, 0)
	assertRandomStringLength(t, 16)
	assertRandomStringLength(t, 128)
	assertRandomStringLength(t, 150)

	str := fields.RandomString(256)
	if str[128:] == str[:128] {
		t.Fatal("The string is repeative and not random.")
	}
}

func TestPasswordMatch(t *testing.T) {
	password := fields.GenPassword("example")
	if !password.MatchesWith("example") {
		t.Fatal("Passwords must match")
	}
	if password.MatchesWith("example2") {
		t.Fatal("Passwords must not match")
	}
	if password.MatchesWith("elpmaxe") {
		t.Fatal("Passwords must not match")
	}
}

func TestPasswordParse(t *testing.T) {
	password := fields.GenPassword("example")

	got, err := fields.ParsePassword(password.String())
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, password.String(), got.String())
}
