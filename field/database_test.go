package field_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
)

func getChunkBuffer(chunks []string, fn func(string) []byte) bytes.Buffer {
	var buffer bytes.Buffer

	byteHeader := make([]byte, 8)
	binary.BigEndian.PutUint64(byteHeader, uint64(len(chunks)))
	buffer.Write(byteHeader)

	for _, message := range chunks {
		buffer.Write(fn(message))
	}

	return buffer
}

func TestDatabaseEmpty(t *testing.T) {
	buffer := getChunkBuffer([]string{}, func(str string) []byte {
		var data [1024]byte
		copy(data[:], []byte(str))
		return data[:]
	})

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()), 1024)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.ReadChunk(0)
	if err == nil {
		t.Fatal("Should fail because the file is empty!")
	}
}

func TestDatabaseSingle(t *testing.T) {
	msg := "Hello world!\n# This is markdown! 🎬\nWith UTF-8 ENCODING!"
	buffer := getChunkBuffer([]string{msg}, func(str string) []byte {
		var data [1024]byte
		copy(data[:], []byte(str))
		return data[:]
	})

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()), 1024)
	if err != nil {
		t.Fatal(err)
	}

	chunk, err := db.ReadChunk(0)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, chunk, msg)
}

func TestDatabaseOverflow(t *testing.T) {
	msg := "Hello world!\n# This is markdown! 🎬\nWith UTF-8 ENCODING!"
	buffer := getChunkBuffer([]string{msg}, func(str string) []byte {
		var data [16]byte
		copy(data[:], []byte(str))
		return data[:]
	})

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()), 16)
	if err != nil {
		t.Fatal(err)
	}

	chunk, err := db.ReadChunk(0)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, chunk, msg[:16])
}

func TestDatabaseMultiple(t *testing.T) {
	msg := "Hello world!\n# This is markdown! 🎬\nWith UTF-8 ENCODING!"
	buffer := getChunkBuffer([]string{msg, msg, msg}, func(str string) []byte {
		var data [1024]byte
		copy(data[:], []byte(str))
		return data[:]
	})

	db, err := field.NewDatabase(NewMockFile(buffer.Bytes()), 1024)
	if err != nil {
		t.Fatal(err)
	}

	chunk, err := db.ReadChunk(1)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, chunk, msg)
}

func TestDatabaseWrite(t *testing.T) {
	msg := "# Hello World!\nThis is an example."

	buffer := getChunkBuffer([]string{msg}, func(str string) []byte {
		var data [1024]byte
		copy(data[:], []byte(str))
		return data[:]
	})

	file := NewMockFile(buffer.Bytes())
	file.WriteAt(buffer.Bytes(), 0)

	db, err := field.NewDatabase(file, 1024)
	if err != nil {
		t.Fatal(err)
	}

	err = db.WriteChunk(msg)
	if err != nil {
		t.Fatal(err)
	}

	var data = make([]byte, 1024)
	file.Read(data)

	str, err := db.ReadChunk(0)
	if err != nil {
		t.Fatal(err)
	}

	parser.Assert(t, str, msg)
}
