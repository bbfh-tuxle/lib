package channels_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/channels"
	"gotest.tools/assert"
)

func makeEntriesBuffer(messages []*channels.Entry) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	stream.WriteUint64(&buffer, uint64(len(messages)))

	for _, message := range messages {
		err := message.Write(&buffer)
		if err != nil {
			return buffer, err
		}
	}

	return buffer, nil
}

func TestMessageListEmpty(t *testing.T) {
	buffer, err := makeEntriesBuffer([]*channels.Entry{})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewListFile(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ReadOldestEntry(0)
	if err == nil {
		t.Fatal("Should fail because the file is empty!")
	}
}

func TestMessageListSingle(t *testing.T) {
	msg := &channels.Entry{
		Timestamp: 1722765647,
		ChunkId:   2,
		ChunkLine: 6,
		UserId:    "1234567890abcdef",
	}
	buffer, err := makeEntriesBuffer([]*channels.Entry{
		msg,
	})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewListFile(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	message, err := db.ReadNewestEntry(0)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, message, msg)
}

func TestMessageListLatest(t *testing.T) {
	latest := &channels.Entry{
		Timestamp: 1722765647,
		ChunkId:   2,
		ChunkLine: 6,
		UserId:    "1234567890abcdef",
	}
	buffer, err := makeEntriesBuffer([]*channels.Entry{
		{
			Timestamp: 1722765698,
			ChunkId:   1,
			ChunkLine: 5681285,
			UserId:    "a",
		},
		latest,
	})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewListFile(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	message, err := db.ReadNewestEntry(0)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, message, latest)

	oldest, err := db.ReadOldestEntry(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, message, oldest)

	_, err = db.ReadOldestEntry(2)
	if err != io.EOF {
		t.Fatal(err)
	}
}
