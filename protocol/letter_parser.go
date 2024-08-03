package protocol

import (
	"bufio"
	"encoding/binary"
	"io"

	"github.com/bbfh-tuxle/lib/field"
)

func readString(reader *bufio.Reader, delimiter byte) (string, error) {
	data, err := reader.ReadString(delimiter)
	if err != nil {
		return "", err
	}

	return data[:len(data)-1], nil
}

func readUInt32(reader *bufio.Reader, delimiter byte) (uint32, error) {
	data, err := reader.ReadBytes(delimiter)
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint32(data[:len(data)-1]), nil
}

// Reads (parses) Letter.Type (and only that).
//
// Use to avoid parsing the entire letter when only type is required (e.g. routing).
func ReadLetterType(stream io.Reader) (string, error) {
	return readString(bufio.NewReader(stream), ' ')
}

// Reads (parses) the entire Letter.
func ReadLetter(stream io.Reader) (Letter, error) {
	reader := bufio.NewReader(stream)

	letterType, err := readString(reader, ' ')
	if err != nil {
		return Letter{}, err
	}

	letterEndpoint, err := readString(reader, '\n')
	if err != nil {
		return Letter{}, err
	}

	numberOfParameters, err := readUInt32(reader, '\n')
	if err != nil {
		return Letter{}, err
	}

	var parameters = field.NewParameters()
	for range numberOfParameters {
		key, err := readString(reader, '=')
		if err != nil {
			return Letter{}, err
		}

		value, err := readString(reader, '\n')
		if err != nil {
			return Letter{}, err
		}

		parameters[key] = value
	}

	body, err := readString(reader, '\r')
	if err != nil {
		return Letter{}, err
	}

	return Letter{
		Type:       letterType,
		Endpoint:   letterEndpoint,
		Parameters: parameters,
		Body:       body,
	}, nil
}
