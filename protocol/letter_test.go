package protocol_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"strings"
	"testing"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
	"github.com/bbfh-tuxle/lib/protocol"
)

func TestLetterType(t *testing.T) {
	data, err := protocol.ReadLetterType(bufio.NewReader(strings.NewReader("DONE message\r")))
	if err != nil {
		t.Fatal(err)
	}
	parser.Assert(t, data, "DONE")

	data, err = protocol.ReadLetterType(bufio.NewReader(strings.NewReader("INVALID")))
	if err != io.EOF {
		t.Fatal(err)
	}
}

func TestLetterNoParams(t *testing.T) {
	var number uint32 = 0
	// Convert the uint32 value to a byte slice (little-endian)
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, number)

	// Combine all parts
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	buffer.Write(byteSlice)
	buffer.WriteString("\nBody\r")

	data, err := protocol.ReadLetter(bufio.NewReader(strings.NewReader(buffer.String())))
	if err != nil {
		t.Fatal(err)
	}
	parser.Assert(t, data, protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: field.NewParameters(),
		Body:       "Body",
	})
}

func TestLetterHeaderOnly(t *testing.T) {
	var number uint32 = 0
	// Convert the uint32 value to a byte slice (little-endian)
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, number)

	// Combine all parts
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	buffer.Write(byteSlice)
	buffer.WriteString("\n\r")

	data, err := protocol.ReadLetter(bufio.NewReader(strings.NewReader(buffer.String())))
	if err != nil {
		t.Fatal(err)
	}
	parser.Assert(t, data, protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: field.NewParameters(),
		Body:       "",
	})
}

func TestLetterWrongType(t *testing.T) {
	_, err := protocol.ReadLetter(bufio.NewReader(strings.NewReader("TEST message\nINVALID\n\r")))
	if err == nil {
		t.Fatal("Expected to fail")
	}
}

func TestLetterParameters(t *testing.T) {
	var number uint32 = 2
	// Convert the uint32 value to a byte slice (little-endian)
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, number)

	// Combine all parts
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	buffer.Write(byteSlice)
	buffer.WriteString("\nKeyA=ValueA\nKeyB=ValueB\n\r")

	data, err := protocol.ReadLetter(bufio.NewReader(strings.NewReader(buffer.String())))
	if err != nil {
		t.Fatal(err)
	}
	params := field.NewParameters()
	params["KeyA"] = "ValueA"
	params["KeyB"] = "ValueB"
	parser.Assert(t, data, protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: params,
		Body:       "",
	})
}

func TestLetterComplete(t *testing.T) {
	body := "This is an example body.\nIt should work!\n\nEven with multiple spaces!"
	var number uint32 = 1

	// Convert the uint32 value to a byte slice (little-endian)
	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, number)

	// Combine all parts
	var buffer bytes.Buffer
	buffer.WriteString("TEST message\n")
	buffer.Write(byteSlice)
	buffer.WriteString("\nKeyA=ValueA\n")
	buffer.WriteString(body)
	buffer.WriteString("\r")

	data, err := protocol.ReadLetter(bufio.NewReader(strings.NewReader(buffer.String())))
	if err != nil {
		t.Fatal(err)
	}
	params := field.NewParameters()
	params["KeyA"] = "ValueA"
	parser.Assert(t, data, protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: params,
		Body:       body,
	})
}

func TestLetterWriteComplete(t *testing.T) {
	body := "This is an example body.\nIt should work!\n\nEven with multiple spaces!"

	params := field.NewParameters()
	params["KeyA"] = "ValueA"

	letter := protocol.Letter{
		Type:       "TEST",
		Endpoint:   "message",
		Parameters: params,
		Body:       body,
	}

	var buffer bytes.Buffer
	letter.Write(&buffer)

	got, err := protocol.ReadLetter(bufio.NewReader(strings.NewReader(buffer.String())))
	if err != nil {
		t.Fatal(err)
	}
	parser.Assert(t, got, letter)
}
