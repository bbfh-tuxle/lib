package protocol

import (
	"io"

	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/fields"
)

// A letter is used to communicate between programs.
type Letter struct {
	Type       string
	Endpoint   string
	Parameters fields.Parameters
	Body       string
}

// Write the letter to an io.Writter in the correct format.
func (letter *Letter) Write(buffer io.Writer) error {
	buffer.Write([]byte(letter.Type))
	buffer.Write([]byte{' '})
	buffer.Write([]byte(letter.Endpoint))
	buffer.Write([]byte{'\n'})
	stream.WriteUint32(buffer, uint32(len(letter.Parameters)))
	letter.Parameters.Write(buffer)
	buffer.Write([]byte(letter.Body))
	buffer.Write([]byte{'\r'})

	return nil
}
