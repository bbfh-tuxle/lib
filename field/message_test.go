package field_test

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
)

func TestParseMessage(t *testing.T) {
	message := field.Message{
		Timestamp: 1722765647,
		ChunkId:   1,
		ChunkLine: 0,
		UserId:    "root",
	}

	var buffer bytes.Buffer
	message.Write(&buffer)

	got, err := field.ReadMessage(bufio.NewReader(strings.NewReader(buffer.String())))
	if err != nil {
		t.Fatal(err)
	}
	parser.Assert(t, got, message)
}

func TestParseMultipleMessages(t *testing.T) {
	messages := []field.Message{
		{
			Timestamp: 1722765647,
			ChunkId:   2,
			ChunkLine: 6,
			UserId:    "1234567890abcdef",
		},
		{
			Timestamp: 1722765698,
			ChunkId:   1,
			ChunkLine: 5681285,
			UserId:    "a",
		},
	}

	// Separate into separate for loops to simulate messages being stored in a sequence
	var buffer bytes.Buffer
	for _, message := range messages {
		message.Write(&buffer)
	}

	var reader = bufio.NewReader(strings.NewReader(buffer.String()))
	for _, message := range messages {
		got, err := field.ReadMessage(reader)
		if err != nil {
			t.Fatal(err)
		}
		parser.Assert(t, got, message)
	}
}
