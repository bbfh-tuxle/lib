package fields

import (
	"fmt"
	"strings"
)

type Permissions struct {
	Server   permMask
	Messages permMask
	Channels permMask
	Users    permMask
}

func EmptyPermissions(all permMask) Permissions {
	return Permissions{
		Server:   all,
		Messages: all,
		Channels: all,
		Users:    all,
	}
}

func (permissions Permissions) String() string {
	return strings.Join([]string{
		fmt.Sprintf("server=%s", permissions.Server),
		fmt.Sprintf("channels=%s", permissions.Channels),
		fmt.Sprintf("messages=%s", permissions.Messages),
		fmt.Sprintf("users=%s", permissions.Users),
	}, ";")
}
