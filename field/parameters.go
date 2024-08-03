package field

import "io"

type Parameters map[string]string

func NewParameters() Parameters {
	return Parameters{}
}

// Write the Parameters to an io.Writter in the correct format.
//
// No-op when Parameters are empty.
func (params Parameters) Write(buffer io.Writer) {
	for key, value := range params {
		buffer.Write([]byte(key))
		buffer.Write([]byte{'='})
		buffer.Write([]byte(value))
		buffer.Write([]byte{'\n'})
	}
}
