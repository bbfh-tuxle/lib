package channels

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/bbfh-tuxle/lib/internal/stream"
)

// A single `.db` file with helper functions for reading/writing.
type Database struct {
	Size      int64
	ChunkSize int64
	file      File
}

func NewDatabase(file File, chunkSize int64) (*Database, error) {
	size, err := stream.ReadInt64(file)
	if err != nil {
		var buffer bytes.Buffer
		stream.WriteInt64(&buffer, 0)

		_, err = file.WriteAt(buffer.Bytes(), 0)
		if err != nil {
			return nil, err
		}
	}

	return &Database{
		Size:      size,
		ChunkSize: chunkSize,
		file:      file,
	}, nil
}

// Reads chunk at a specific index.
//
// Returns io.EOF when reading out of bounds data
func (db *Database) ReadChunk(index int64) (string, error) {
	if db.Size == 0 {
		return "", errors.New("Trying to read an empty database!")
	}

	var data = make([]byte, db.ChunkSize)
	_, err := db.file.ReadAt(data, LIST_CONTENT_OFFSET+db.ChunkSize*index)
	if err != nil {
		return "", err
	}

	return string(bytes.Trim(data, "\x00")), nil
}

func (db *Database) IncrementSize() error {
	db.Size += 1

	byteSize := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSize, uint64(db.Size))

	_, err := db.file.WriteAt(byteSize, 0)
	return err
}

// Appends a chunk to the end of the database.
func (db *Database) AppendChunk(chunk string) error {
	var data = make([]byte, db.ChunkSize)
	copy(data[:], []byte(chunk))

	err := db.IncrementSize()
	if err != nil {
		return err
	}

	_, err = db.file.WriteAt(data, LIST_CONTENT_OFFSET+db.ChunkSize*(db.Size-1))
	return err
}

// Writes a chunk to certain position in the database.
func (db *Database) OverwriteChunk(chunk string, index int64) error {
	if index > db.Size {
		return errors.New("Index is out of bounds!")
	}

	var data = make([]byte, db.ChunkSize)
	copy(data[:], []byte(chunk))

	_, err := db.file.WriteAt(data, LIST_CONTENT_OFFSET+db.ChunkSize*index)
	return err
}

// Appends a chunk to the end of the database.
//
// Limits the actual size of the file to several elements, replacing the oldest ones when needed.
func (db *Database) LimitedAppendChunk(chunk string, cycle int64) error {
	if db.Size < cycle {
		return db.AppendChunk(chunk)
	}

	err := db.IncrementSize()
	if err != nil {
		return err
	}

	return db.OverwriteChunk(chunk, db.Size%cycle)
}

// Returns whether database contains the chunk.
func (db *Database) ContainsChunk(chunk string) bool {
	for i := range db.Size {
		line, err := db.ReadChunk(i)
		if err != nil {
			break
		}

		if line == chunk {
			return true
		}
	}

	return false
}
