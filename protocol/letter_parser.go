package protocol

import (
	"bufio"

	"github.com/bbfh-tuxle/lib/field"
	"github.com/bbfh-tuxle/lib/field/parser"
)

// Reads (parses) Letter.Type (and only that).
//
// Use to avoid parsing the entire letter when only type is required (e.g. routing).
//
// Returns EOF if Parameters format is invalid.
func ReadLetterType(reader *bufio.Reader) (string, error) {
	return parser.ReadString(reader, ' ')
}

// Reads (parses) the entire Letter.
//
// Returns EOF if Parameters format is invalid.
func ReadLetter(reader *bufio.Reader) (Letter, error) {
	letterType, err := parser.ReadString(reader, ' ')
	if err != nil {
		return Letter{}, err
	}

	letterEndpoint, err := parser.ReadString(reader, '\n')
	if err != nil {
		return Letter{}, err
	}

	numberOfParameters, err := parser.ReadUInt32(reader, '\n')
	if err != nil {
		return Letter{}, err
	}

	parameters, err := field.ReadParameters(reader, int(numberOfParameters))

	body, err := parser.ReadString(reader, '\r')
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
