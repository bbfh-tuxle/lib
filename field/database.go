package field

import (
	"bytes"
	"encoding/binary"
	"errors"
)

// A single `.db` file with helper functions for reading/writing.
type Database struct {
	Size      uint64 // The amount of messages in the message list
	ChunkSize int64
	file      File
}

func NewDatabase(file File, chunkSize int64) (*Database, error) {
	var headerBytes = make([]byte, LIST_HEADER_SIZE)
	_, err := file.Read(headerBytes)
	if err != nil {
		return nil, err
	}

	return &Database{
		Size:      binary.BigEndian.Uint64(headerBytes),
		ChunkSize: chunkSize,
		file:      file,
	}, nil
}

// Reads the text chunk at specified index with all null-bytes trimmed.
func (db *Database) ReadChunk(index int64) (string, error) {
	if db.Size == 0 {
		return "", errors.New("Trying to read an empty message list!")
	}

	var data = make([]byte, db.ChunkSize)
	_, err := db.file.ReadAt(data, LIST_HEADER_SIZE+db.ChunkSize*index)
	if err != nil {
		return "", err
	}

	return string(bytes.Trim(data, "\x00")), nil
}

// Reads the text chunk at specified index with all null-bytes trimmed.
func (db *Database) WriteChunk(chunk string) error {
	db.Size += 1

	var data = make([]byte, db.ChunkSize)
	copy(data[:], []byte(chunk))

	byteSize := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSize, db.Size)

	_, err := db.file.WriteAt(byteSize, 0)
	if err != nil {
		return err
	}

	_, err = db.file.WriteAt(data, LIST_HEADER_SIZE+db.ChunkSize*int64(db.Size))
	return err
}
