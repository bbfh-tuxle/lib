package channels

import (
	"io"

	"github.com/bbfh-tuxle/lib/internal/stream"
)

func ReadEntry(reader io.Reader) (*Entry, error) {
	timestamp, err := stream.ReadUint64(reader)
	if err != nil {
		return nil, err
	}

	var chunkId = make([]byte, 1)
	_, err = io.ReadFull(reader, chunkId)
	if err != nil {
		return nil, err
	}

	chunkLine, err := stream.ReadUint64(reader)
	if err != nil {
		return nil, err
	}

	userId, err := stream.ReadChars47(reader)
	if err != nil {
		return nil, err
	}

	return &Entry{
		Timestamp: timestamp,
		ChunkId:   chunkId[0],
		ChunkLine: chunkLine,
		UserId:    userId,
	}, nil
}
