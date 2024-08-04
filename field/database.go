package field

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// The size of the Database header size in bytes
const DB_HEADER_SIZE int64 = 8

// A file that database will read from.
//
// An *os.File will always satisfy this interface.
type File interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
}

// A single `.db` file with helper functions for reading/writing.
type Database struct {
	Size uint64 // The amount of messages in the database
	file File
}

func NewDatabase(file File) (Database, error) {
	var headerBytes = make([]byte, DB_HEADER_SIZE)
	_, err := file.Read(headerBytes)
	if err != nil {
		return Database{}, err
	}

	return Database{
		Size: binary.BigEndian.Uint64(headerBytes),
		file: file,
	}, nil
}

// Reads message at a certain index starting from the OLDEST message.
func (db Database) ReadOldestMessage(index int64) (Message, error) {
	if db.Size == 0 {
		return Message{}, errors.New("Trying to read an empty database!")
	}

	var data = make([]byte, MESSAGE_SIZE)
	_, err := db.file.ReadAt(data, DB_HEADER_SIZE+MESSAGE_SIZE*index)
	if err != nil {
		return Message{}, err
	}

	return ReadMessage(bytes.NewReader(data))
}

// Reads message at a certain index starting from the NEWEST message.
//
// Returns an error if index is out of bounds.
func (db Database) ReadNewestMessage(index int64) (Message, error) {
	if db.Size == 0 {
		return Message{}, errors.New("Trying to read an empty database!")
	}

	if index > int64(db.Size) {
		return Message{}, fmt.Errorf("Index (%d) is out of bounds (%d)!", index, db.Size)
	}

	var data = make([]byte, MESSAGE_SIZE)
	_, err := db.file.ReadAt(data, DB_HEADER_SIZE+MESSAGE_SIZE*(int64(db.Size)-index-1))
	if err != nil {
		return Message{}, err
	}

	return ReadMessage(bytes.NewReader(data))
}
