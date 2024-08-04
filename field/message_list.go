package field

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

// The size of the MessageList header size in bytes
const LIST_HEADER_SIZE int64 = 8

// A file that message list will read from.
//
// An *os.File will always satisfy this interface.
type File interface {
	Read([]byte) (int, error)
	ReadAt([]byte, int64) (int, error)
}

// A single `.list` file with helper functions for reading/writing.
type MessageList struct {
	Size uint64 // The amount of messages in the message list
	file File
}

func NewMessageList(file File) (MessageList, error) {
	var headerBytes = make([]byte, LIST_HEADER_SIZE)
	_, err := file.Read(headerBytes)
	if err != nil {
		return MessageList{}, err
	}

	return MessageList{
		Size: binary.BigEndian.Uint64(headerBytes),
		file: file,
	}, nil
}

// Reads message at a certain index starting from the OLDEST message.
func (list MessageList) ReadOldestMessage(index int64) (Message, error) {
	if list.Size == 0 {
		return Message{}, errors.New("Trying to read an empty message list!")
	}

	var data = make([]byte, MESSAGE_SIZE)
	_, err := list.file.ReadAt(data, LIST_HEADER_SIZE+MESSAGE_SIZE*index)
	if err != nil {
		return Message{}, err
	}

	return ReadMessage(bytes.NewReader(data))
}

// Reads message at a certain index starting from the NEWEST message.
//
// Returns an error if index is out of bounds.
func (list MessageList) ReadNewestMessage(index int64) (Message, error) {
	if list.Size == 0 {
		return Message{}, errors.New("Trying to read an empty message list!")
	}

	if index > int64(list.Size) {
		return Message{}, fmt.Errorf("Index (%d) is out of bounds (%d)!", index, list.Size)
	}

	var data = make([]byte, MESSAGE_SIZE)
	_, err := list.file.ReadAt(data, LIST_HEADER_SIZE+MESSAGE_SIZE*(int64(list.Size)-index-1))
	if err != nil {
		return Message{}, err
	}

	return ReadMessage(bytes.NewReader(data))
}
