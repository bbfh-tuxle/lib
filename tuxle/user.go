package tuxle

import (
	"io"

	"github.com/bbfh-tuxle/lib/internal/escape"
	"github.com/bbfh-tuxle/lib/tuxle/fields"
)

type User struct {
	Id          string
	Name        string
	PictureURI  string
	Description string
	Password    fields.Password
	Permissions fields.Permissions
}

func (user User) Parameters() fields.Parameters {
	return fields.Parameters{
		"Name":        user.Name,
		"PictureURI":  user.PictureURI,
		"Description": escape.EscapeString(user.Description),
		"Password":    user.Password.String(),
		"Permissions": user.Permissions.String(),
	}
}

func (user User) Write(buffer io.Writer) {
	fields.Parameters{
		"Name":        user.Name,
		"PictureURI":  user.PictureURI,
		"Description": escape.EscapeString(user.Description),
		"Password":    user.Password.String(),
		"Permissions": user.Permissions.String(),
	}.Write(buffer)
}
