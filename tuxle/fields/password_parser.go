package fields

import (
	"errors"
	"strings"
)

func ParsePassword(str string) (Password, error) {
	partitions := strings.SplitN(str, ".", 2)
	if len(partitions) != 2 {
		return Password{}, errors.New("Invalid password format. Requires `salt.hash`")
	}

	return Password{
		hash: partitions[1],
		salt: partitions[0],
	}, nil
}

func GenPassword(password string) Password {
	salt := RandomString(SALT_LENGTH)
	return Password{
		hash: encodeString(salt + password),
		salt: salt,
	}
}
