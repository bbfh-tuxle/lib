package channels_test

import (
	"bytes"
	"testing"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/channels"
	"gotest.tools/assert"
)

func makeChunksBuffer256(chunks []string) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	stream.WriteUint64(&buffer, uint64(len(chunks)))

	for _, chunk := range chunks {
		stream.WriteChars256(&buffer, chunk)
	}

	return buffer, nil
}

func makeChunksBuffer16(chunks []string) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	stream.WriteUint64(&buffer, uint64(len(chunks)))

	for _, chunk := range chunks {
		stream.WriteChars16(&buffer, chunk)
	}

	return buffer, nil
}

func TestDatabaseEmpty(t *testing.T) {
	buffer, err := makeChunksBuffer256([]string{})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewDatabase(NewMockFile(buffer.Bytes()), 256)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ReadChunk(0)
	if err == nil {
		t.Fatal("Should fail because the file is empty!")
	}
}

func TestDatabaseSingle(t *testing.T) {
	msg1 := "Hello World!"
	buffer, err := makeChunksBuffer256([]string{msg1})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewDatabase(NewMockFile(buffer.Bytes()), 256)
	if err != nil {
		t.Fatal(err)
	}

	chunk, err := db.ReadChunk(0)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, chunk, msg1)
}

func TestDatabaseMultiple(t *testing.T) {
	msg1 := "Hello World!"
	buffer, err := makeChunksBuffer256([]string{
		"This is a very long message maybe that doesn't fit completely! 🧩",
		msg1,
		"Some message after the chunk.",
	})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewDatabase(NewMockFile(buffer.Bytes()), 256)
	if err != nil {
		t.Fatal(err)
	}

	chunk, err := db.ReadChunk(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, chunk, msg1)
}

func TestDatabaseOverwrite(t *testing.T) {
	msg1 := "Hello World!"
	buffer, err := makeChunksBuffer256([]string{
		"This is a very long message that maybe doesn't fit completely! 🧩",
		"Testing...",
		"Some message after the chunk.",
	})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewDatabase(NewMockFile(buffer.Bytes()), 256)
	if err != nil {
		t.Fatal(err)
	}

	db.OverwriteChunk(msg1, 1)
	chunk, err := db.ReadChunk(1)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, chunk, msg1)
}

func TestDatabaseAppend(t *testing.T) {
	msg1 := "Hello World!"
	buffer, err := makeChunksBuffer16([]string{
		"This is a very long message that maybe doesn't fit completely! 🧩",
		"Testing...",
	})
	if err != nil {
		t.Fatal(err)
	}

	db, err := channels.NewDatabase(NewMockFile(buffer.Bytes()), 16)
	if err != nil {
		t.Fatal(err)
	}

	db.AppendChunk(msg1)
	chunk, err := db.ReadChunk(2)
	if err != nil {
		t.Fatal(err)
	}

	assert.DeepEqual(t, chunk, msg1)
}
