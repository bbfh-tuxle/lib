package fields_test

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/bbfh-tuxle/lib/tuxle/fields"
	"gotest.tools/assert"
)

func TestSession(t *testing.T) {
	session := fields.SessionKey{
		UserId:  "root",
		Token:   "4936c321070dbb1d067506b692e993e7bdf5ca70512e1237e20afd795f6b0e5f",
		Address: "127.0.0.1:16400",
	}

	var buffer bytes.Buffer
	session.Write(&buffer)
	got, err := fields.ReadSessionKey(bufio.NewReader(bytes.NewReader(buffer.Bytes())))
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, got, session)
}

func TestNoAddress(t *testing.T) {
	session := fields.SessionKey{
		UserId:  "root",
		Token:   "4936c321070dbb1d067506b692e993e7bdf5ca70512e1237e20afd795f6b0e5f",
		Address: "",
	}

	var buffer bytes.Buffer
	session.Write(&buffer)
	_, err := fields.ReadSessionKey(bufio.NewReader(bytes.NewReader(buffer.Bytes())))
	if err == nil {
		t.Fatal("Must fail!")
	}
}
