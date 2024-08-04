package field_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
)

func getMessagesBuffer(messages []field.Message) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	byteHeader := make([]byte, 8)
	binary.BigEndian.PutUint64(byteHeader, uint64(len(messages)))
	buffer.Write(byteHeader)

	for _, message := range messages {
		err := message.Write(&buffer)
		if err != nil {
			return buffer, err
		}
	}

	return buffer, nil
}

func TestMessageListEmpty(t *testing.T) {
	buffer, err := getMessagesBuffer([]field.Message{})
	if err != nil {
		t.Fatal(err)
	}

	db, err := field.NewMessageList(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ReadOldestMessage(0)
	if err == nil {
		t.Fatal("Should fail because the file is empty!")
	}
}

func TestMessageListSingle(t *testing.T) {
	msg := field.Message{
		Timestamp: 1722765647,
		ChunkId:   2,
		ChunkLine: 6,
		UserId:    "1234567890abcdef",
	}
	buffer, err := getMessagesBuffer([]field.Message{
		msg,
	})
	if err != nil {
		t.Fatal(err)
	}

	db, err := field.NewMessageList(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	message, err := db.ReadNewestMessage(0)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, message, msg)
}

func TestMessageListLatest(t *testing.T) {
	latest := field.Message{
		Timestamp: 1722765647,
		ChunkId:   2,
		ChunkLine: 6,
		UserId:    "1234567890abcdef",
	}
	buffer, err := getMessagesBuffer([]field.Message{
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

	db, err := field.NewMessageList(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	message, err := db.ReadNewestMessage(0)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, message, latest)

	oldest, err := db.ReadOldestMessage(1)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, message, oldest)

	_, err = db.ReadOldestMessage(2)
	if err != io.EOF {
		t.Fatal(err)
	}
}
