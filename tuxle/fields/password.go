package fields

import (
	"fmt"
)

type Password struct {
	hash string // 128 characters long SHA512
	salt string
}

func (password Password) String() string {
	return fmt.Sprintf("%s.%s", password.salt, password.hash)
}

func (password Password) MatchesWith(str string) bool {
	return password.hash == encodeString(password.salt+str)
}
