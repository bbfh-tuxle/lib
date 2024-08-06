package fields

import (
	"fmt"
	"io"
)

// A map with extra convinient functions as a part of Tuxle format.
type Parameters map[string]string

// Write the Parameters to an io.Writter in the correct format.
//
// No-op when Parameters are empty.
func (params Parameters) Write(buffer io.Writer) {
	for key, value := range params {
		fmt.Fprintf(buffer, "%s=%s\n", key, value)
	}
}
