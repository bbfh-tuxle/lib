package protocol_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/fields"
	"github.com/bbfh-tuxle/lib/tuxle/protocol"
	"gotest.tools/assert"
)

func TestLetterType(t *testing.T) {
	data, err := protocol.ReadLetterType(stream.NewReader("DONE message\r"))
	if err != nil {
		t.Fatal(err)
	}
	assert.DeepEqual(t, data, "DONE")

	data, err = protocol.ReadLetterType(stream.NewReader("INVALID"))
	if err != io.EOF {
		t.Fatal(err)
	}
}

func TestLetterNoParams(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	stream.WriteUint32(&buffer, 0)
	buffer.WriteString("Body\r")

	data, err := protocol.ReadLetter(stream.NewReader(buffer.String()))
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, data, &protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: fields.Parameters{},
		Body:       "Body",
	})
}

func TestLetterHeaderOnly(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	stream.WriteUint32(&buffer, 0)
	buffer.WriteString("\r")

	data, err := protocol.ReadLetter(stream.NewReader(buffer.String()))
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, data, &protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: fields.Parameters{},
		Body:       "",
	})
}

func TestLetterWrongType(t *testing.T) {
	_, err := protocol.ReadLetter(stream.NewReader("TEST message\nINVALID\r"))
	if err == nil {
		t.Fatal("Expected to fail")
	}
}

func TestLetterParameters(t *testing.T) {
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	stream.WriteUint32(&buffer, 2)
	buffer.WriteString("KeyA=ValueA\nKeyB=ValueB\n\r")

	data, err := protocol.ReadLetter(stream.NewReader(buffer.String()))
	if err != nil {
		t.Fatal(err)
	}
	params := fields.Parameters{}
	params["KeyA"] = "ValueA"
	params["KeyB"] = "ValueB"
	assert.DeepEqual(t, data, &protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: params,
		Body:       "",
	})
}

func TestLetterComplete(t *testing.T) {
	body := "This is an example body.\nIt should work!\n\nEven with multiple spaces!"

	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	stream.WriteUint32(&buffer, 1)
	buffer.WriteString("KeyA=ValueA\n")
	buffer.WriteString(body)
	buffer.WriteString("\r")

	data, err := protocol.ReadLetter(stream.NewReader(buffer.String()))
	if err != nil {
		t.Fatal(err)
	}
	params := fields.Parameters{}
	params["KeyA"] = "ValueA"
	assert.DeepEqual(t, data, &protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: params,
		Body:       body,
	})
}

func TestLetterWriteComplete(t *testing.T) {
	body := "This is an example body.\nIt should work!\n\nEven with multiple spaces!"

	params := fields.Parameters{}
	params["KeyA"] = "ValueA"

	letter := protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: params,
		Body:       body,
	}

	var buffer bytes.Buffer
	letter.Write(&buffer)

	got, err := protocol.ReadLetter(stream.NewReader(buffer.String()))
	if err != nil {
		t.Fatal(err)
	}
	assert.DeepEqual(t, got, &letter)
}
