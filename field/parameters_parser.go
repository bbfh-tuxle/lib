package field

import (
	"bufio"
	"io"

	"github.com/bbfh-tuxle/lib/field/parser"
)

func readEntireParameters(reader *bufio.Reader) (Parameters, error) {
	parameters := NewParameters()

	for {
		key, err := parser.ReadString(reader, '=')
		if err != nil {
			if err == io.EOF {
				break
			}
			return parameters, err
		}

		value, err := parser.ReadString(reader, '\n')
		if err != nil {
			return parameters, err
		}

		parameters[key] = value
	}

	return parameters, nil
}

// Read (parse) Parameters.
//
// `count` — The amount of lines to expect. Set to -1 to parse until io.EOF
func ReadParameters(reader *bufio.Reader, count int) (Parameters, error) {
	if count == -1 {
		return readEntireParameters(reader)
	}

	parameters := NewParameters()
	for range count {
		key, err := parser.ReadString(reader, '=')
		if err != nil {
			return parameters, err
		}

		value, err := parser.ReadString(reader, '\n')
		if err != nil {
			return parameters, err
		}

		parameters[key] = value
	}

	return parameters, nil
}
