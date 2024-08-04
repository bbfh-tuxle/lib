package field

import (
	"encoding/binary"
	"io"
)

const MESSAGE_SIZE int64 = 8 + 1 + 8 + 16

type Message struct {
	Timestamp uint64
	ChunkId   byte
	ChunkLine uint64
	UserId    string
}

func (message Message) Write(buffer io.Writer) error {
	byteTimestamp := make([]byte, 8)
	binary.BigEndian.PutUint64(byteTimestamp, message.Timestamp)
	buffer.Write(byteTimestamp)

	buffer.Write([]byte{message.ChunkId})

	byteLine := make([]byte, 8)
	binary.BigEndian.PutUint64(byteLine, message.ChunkLine)
	buffer.Write(byteLine)

	var byteUser [16]byte
	copy(byteUser[:], []byte(message.UserId))
	buffer.Write(byteUser[:])

	return nil
}
