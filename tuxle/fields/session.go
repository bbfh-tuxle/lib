package fields

import (
	"fmt"
	"io"
)

type SessionKey struct {
	UserId  string
	Token   string
	Address string
}

func (session SessionKey) Write(buffer io.Writer) {
	fmt.Fprintf(buffer, "%s+%s@%s", session.UserId, session.Token, session.Address)
}
