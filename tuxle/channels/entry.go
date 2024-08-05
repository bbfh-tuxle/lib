package channels

import (
	"io"

	"github.com/bbfh-tuxle/lib/internal/stream"
)

const ENTRY_SIZE = 8 + 1 + 8 + 47

type Entry struct {
	Timestamp uint64
	ChunkId   byte
	ChunkLine uint64
	UserId    string
}

func (entry Entry) Write(buffer io.Writer) error {
	stream.WriteUint64(buffer, entry.Timestamp)
	buffer.Write([]byte{entry.ChunkId})
	stream.WriteUint64(buffer, entry.ChunkLine)
	stream.WriteChars47(buffer, entry.UserId)

	return nil
}
