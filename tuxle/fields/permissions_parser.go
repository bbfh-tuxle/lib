package fields

import (
	"errors"
	"strings"
)

func parseLine(permissions *Permissions, str string) error {
	parts := strings.SplitN(str, "=", 2)
	if len(parts) != 2 {
		return errors.New("Invalid permission format. Requires `[key]=[permission]`")
	}

	switch parts[0] {
	case "server":
		permissions.Server = permMask(parts[1])
	case "users":
		permissions.Users = permMask(parts[1])
	case "channels":
		permissions.Channels = permMask(parts[1])
	case "messages":
		permissions.Messages = permMask(parts[1])
	}

	return nil
}

func ParsePermissions(mask permMask, str string) (Permissions, error) {
	permissions := EmptyPermissions(mask)

	for _, item := range strings.Split(str, ";") {
		err := parseLine(&permissions, item)
		if err != nil {
			return permissions, err
		}
	}

	return permissions, nil
}
