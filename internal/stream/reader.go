package stream

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"io"
	"strings"
)

func ReadString(reader *bufio.Reader, delimiter byte) (string, error) {
	data, err := reader.ReadString(delimiter)
	if err != nil {
		return data, err
	}

	return data[:len(data)-1], nil
}

func NewReader(str string) *bufio.Reader {
	return bufio.NewReader(strings.NewReader(str))
}

func ReadUint32(reader io.Reader) (uint32, error) {
	var data = make([]byte, 4)

	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(data), nil
}

func ReadUint64(reader io.Reader) (uint64, error) {
	var data = make([]byte, 8)

	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(data), nil
}

func ReadInt64(reader io.Reader) (int64, error) {
	var data = make([]byte, 8)

	_, err := io.ReadFull(reader, data)
	if err != nil {
		return 0, err
	}

	return int64(binary.BigEndian.Uint64(data)), nil
}

func ReadChars47(reader io.Reader) (string, error) {
	var data = make([]byte, 47)

	_, err := io.ReadFull(reader, data)
	if err != nil {
		return "", err
	}

	return string(bytes.Trim(data, "\x00")), nil
}
