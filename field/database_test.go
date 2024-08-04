package field_test

import (
	"bytes"
	"encoding/binary"
	"io"
	"testing"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
)

type MockFile struct {
	content []byte
	offset  int64
}

func NewMockFile(content []byte) *MockFile {
	return &MockFile{content: content, offset: 0}
}

// Read reads up to len(p) bytes into p from the file
func (file *MockFile) Read(p []byte) (int, error) {
	if file.offset >= int64(len(file.content)) {
		return 0, io.EOF
	}
	n := copy(p, file.content[file.offset:])
	file.offset += int64(n)
	return n, nil
}

// ReadAt reads len(p) bytes into p from the file at offset off
func (file *MockFile) ReadAt(p []byte, off int64) (int, error) {
	if off >= int64(len(file.content)) {
		return 0, io.EOF
	}
	n := copy(p, file.content[off:])
	if n < len(p) {
		return n, io.EOF
	}
	return n, nil
}

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

func TestDatabaseEmpty(t *testing.T) {
	buffer, err := getMessagesBuffer([]field.Message{})
	if err != nil {
		t.Fatal(err)
	}

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ReadOldestMessage(0)
	if err == nil {
		t.Fatal("Should fail because the file is empty!")
	}
}

func TestDatabaseSingle(t *testing.T) {
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

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	message, err := db.ReadNewestMessage(0)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, message, msg)
}

func TestDatabaseLatest(t *testing.T) {
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

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()))
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
