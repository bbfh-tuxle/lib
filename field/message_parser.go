package field

import (
	"bytes"
	"encoding/binary"
	"io"
)

func ReadMessage(reader io.Reader) (Message, error) {
	var timestamp = make([]byte, 8)
	_, err := io.ReadFull(reader, timestamp)
	if err != nil {
		return Message{}, err
	}

	var chunkId = make([]byte, 1)
	_, err = io.ReadFull(reader, chunkId)
	if err != nil {
		return Message{}, err
	}

	var chunkLine = make([]byte, 8)
	_, err = io.ReadFull(reader, chunkLine)
	if err != nil {
		return Message{}, err
	}

	var userId = make([]byte, 16)
	_, err = io.ReadFull(reader, userId)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Timestamp: binary.BigEndian.Uint64(timestamp),
		ChunkId:   chunkId[0],
		ChunkLine: binary.BigEndian.Uint64(chunkLine),
		UserId:    string(bytes.Trim(userId, "\x00")),
	}, nil
}
