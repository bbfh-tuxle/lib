package fields

import (
	"bufio"
	"io"

	"github.com/bbfh-tuxle/lib/internal/stream"
)

// Read (parse) Parameters. until io.EOF.
//
// err == io.EOF if parameters have invalid format
func ReadAllParameters(reader *bufio.Reader) (Parameters, error) {
	parameters := Parameters{}

	for {
		key, err := stream.ReadString(reader, '=')
		if err != nil {
			if err == io.EOF {
				break
			}
			return parameters, err
		}

		value, err := stream.ReadString(reader, '\n')
		if err != nil {
			return parameters, err
		}

		parameters[key] = value
	}

	return parameters, nil
}

// Read (parse) Parameters.
//
// `count` — The amount of lines to expect. Use ReadAllParameters() to read until io.EOF.
//
// err == io.EOF if parameters have invalid format or count > number of parameters.
func ReadParameters(reader *bufio.Reader, count int) (Parameters, error) {
	parameters := Parameters{}

	for range count {
		key, err := stream.ReadString(reader, '=')
		if err != nil {
			return parameters, err
		}

		value, err := stream.ReadString(reader, '\n')
		if err != nil {
			return parameters, err
		}

		parameters[key] = value
	}

	return parameters, nil
}
