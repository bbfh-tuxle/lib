package protocol

import (
	"encoding/binary"
	"io"

	"github.com/bbfh-tuxle/lib/field"
)

// A letter is used to communicate between programs.
type Letter struct {
	Type       string
	Endpoint   string
	Parameters field.Parameters
	Body       string
}

// Write the letter to an io.Writter in the correct format.
func (letter Letter) Write(buffer io.Writer) error {
	buffer.Write([]byte(letter.Type))
	buffer.Write([]byte{' '})
	buffer.Write([]byte(letter.Endpoint))
	buffer.Write([]byte{'\n'})

	byteSlice := make([]byte, 4)
	binary.BigEndian.PutUint32(byteSlice, uint32(len(letter.Parameters)))
	buffer.Write(byteSlice)
	buffer.Write([]byte{'\n'})

	for key, value := range letter.Parameters {
		buffer.Write([]byte(key))
		buffer.Write([]byte{'='})
		buffer.Write([]byte(value))
		buffer.Write([]byte{'\n'})
	}

	buffer.Write([]byte(letter.Body))
	buffer.Write([]byte{'\r'})

	return nil
}
