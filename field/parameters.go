package field

import "io"

// A map with extra convinient functions as a part of Tuxle format.
type Parameters map[string]string

// Get a new empty Parameters map.
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
