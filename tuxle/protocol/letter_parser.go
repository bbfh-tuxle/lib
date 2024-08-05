package protocol

import (
	"bufio"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/fields"
)

// Reads (parses) Letter.Type (and only that).
//
// Use to avoid parsing the entire letter when only type is required (e.g. routing).
//
// Returns EOF if Parameters format is invalid.
func ReadLetterType(reader *bufio.Reader) (string, error) {
	return stream.ReadString(reader, ' ')
}

// Reads (parses) the entire Letter.
//
// Returns EOF if Parameters format is invalid.
func ReadLetter(reader *bufio.Reader) (*Letter, error) {
	letterType, err := stream.ReadString(reader, ' ')
	if err != nil {
		return nil, err
	}

	letterEndpoint, err := stream.ReadString(reader, '\n')
	if err != nil {
		return nil, err
	}

	numberOfParameters, err := stream.ReadUint32(reader)
	if err != nil {
		return nil, err
	}

	parameters, err := fields.ReadParameters(reader, int(numberOfParameters))

	body, err := stream.ReadString(reader, '\r')
	if err != nil {
		return nil, err
	}

	return &Letter{
		Type:       letterType,
		Endpoint:   letterEndpoint,
		Parameters: parameters,
		Body:       body,
	}, nil
}
