package channels_test

import (
	"bytes"
	"testing"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/channels"
	"gotest.tools/assert"
)

func TestParseMessage(t *testing.T) {
	message := channels.Entry{
		Timestamp: 1722765647,
		ChunkId:   1,
		ChunkLine: 0,
		UserId:    "root",
	}

	var buffer bytes.Buffer
	message.Write(&buffer)

	got, err := channels.ReadEntry(stream.NewReader(buffer.String()))
	if err != nil {
		t.Fatal(err)
	}
	assert.DeepEqual(t, got, &message)
}

func TestParseMultipleMessages(t *testing.T) {
	messages := []channels.Entry{
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

	var reader = stream.NewReader(buffer.String())
	for _, message := range messages {
		got, err := channels.ReadEntry(reader)
		if err != nil {
			t.Fatal(err)
		}
		assert.DeepEqual(t, got, &message)
	}
}
