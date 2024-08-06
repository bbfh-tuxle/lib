package fields

import (
	"bufio"
	"errors"

	"github.com/bbfh-tuxle/lib/internal/stream"
)

func ReadSessionKey(reader *bufio.Reader) (SessionKey, error) {
	userId, err := stream.ReadString(reader, '+')
	if err != nil {
		return SessionKey{}, err
	}

	token, err := stream.ReadString(reader, '@')
	if err != nil {
		return SessionKey{}, err
	}

	addr, _ := stream.ReadString(reader, '\n')
	if len(addr) == 0 {
		return SessionKey{}, errors.New("Server address is empty!")
	}

	return SessionKey{
		UserId:  userId,
		Token:   token,
		Address: addr,
	}, nil
}
