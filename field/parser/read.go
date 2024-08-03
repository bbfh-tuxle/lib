package parser

import (
	"bufio"
	"encoding/binary"
)

func ReadString(reader *bufio.Reader, delimiter byte) (string, error) {
	data, err := reader.ReadString(delimiter)
	if err != nil {
		return "", err
	}

	return data[:len(data)-1], nil
}

func ReadUInt32(reader *bufio.Reader, delimiter byte) (uint32, error) {
	data, err := reader.ReadBytes(delimiter)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(data[:len(data)-1]), nil
}
