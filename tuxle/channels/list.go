package channels

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/bbfh-tuxle/lib/internal/stream"
)

var LIST_CONTENT_OFFSET int64 = 8

type ListFile struct {
	Size int64 // The number of entries in the list
	file File
}

func NewListFile(file File) (*ListFile, error) {
	size, err := stream.ReadInt64(file)
	if err != nil {
		var buffer bytes.Buffer
		stream.WriteInt64(&buffer, 0)

		_, err = file.WriteAt(buffer.Bytes(), 0)
		if err != nil {
			return nil, err
		}
	}

	return &ListFile{
		Size: size,
		file: file,
	}, nil
}

// Reads entry at a certain index starting from the OLDEST entry.
//
// Entry is nil if an error occured.
func (list *ListFile) ReadOldestEntry(index int64) (*Entry, error) {
	if list.Size == 0 {
		return nil, errors.New("Trying to read an empty list!")
	}

	var data = make([]byte, ENTRY_SIZE)
	_, err := list.file.ReadAt(data, LIST_CONTENT_OFFSET+ENTRY_SIZE*index)
	if err != nil {
		return nil, err
	}

	return ReadEntry(bytes.NewReader(data))
}

// Reads entry at a certain index starting from the NEWEST entry.
//
// Returns an error if index is out of bounds.
//
// Entry is nil if an error occured.
func (list *ListFile) ReadNewestEntry(index int64) (*Entry, error) {
	if list.Size == 0 {
		return nil, errors.New("Trying to read an empty list!")
	}

	if index > list.Size {
		return nil, fmt.Errorf("Index (%d) is out of bounds (%d)!", index, list.Size)
	}

	var data = make([]byte, ENTRY_SIZE)
	_, err := list.file.ReadAt(data, LIST_CONTENT_OFFSET+ENTRY_SIZE*(list.Size-index-1))
	if err != nil {
		return nil, err
	}

	return ReadEntry(bytes.NewReader(data))
}

// Appends an entry to the end of the list file.
func (list *ListFile) AppendEntry(entry *Entry) error {
	list.Size += 1

	byteSize := make([]byte, 8)
	binary.BigEndian.PutUint64(byteSize, uint64(list.Size))

	_, err := list.file.WriteAt(byteSize, 0)
	if err != nil {
		return err
	}

	var buffer bytes.Buffer
	entry.Write(&buffer)

	_, err = list.file.WriteAt(buffer.Bytes(), LIST_CONTENT_OFFSET+ENTRY_SIZE*(list.Size-1))
	return err
}

// Writes an entry to certain position in the list file.
func (list *ListFile) OverwriteEntry(entry *Entry, index int64) error {
	if index > list.Size {
		return errors.New("Index is out of bounds!")
	}

	var buffer bytes.Buffer
	entry.Write(&buffer)

	_, err := list.file.WriteAt(buffer.Bytes(), LIST_CONTENT_OFFSET+ENTRY_SIZE*index)
	return err
}
