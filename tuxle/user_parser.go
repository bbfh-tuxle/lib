package tuxle

import (
	"bufio"

	"github.com/bbfh-tuxle/lib/internal/escape"
	"github.com/bbfh-tuxle/lib/internal/stream"
	"github.com/bbfh-tuxle/lib/tuxle/channels"
	"github.com/bbfh-tuxle/lib/tuxle/fields"
)

func ReadUser(id string, file channels.File) (User, error) {
	parameters, err := fields.ReadAllParameters(bufio.NewReader(file))
	if err != nil {
		return User{}, err
	}

	errs := parameters.Validate(map[string]fields.ValidateFunc{
		"Name":        fields.Exists,
		"PictureURI":  fields.Exists,
		"Description": fields.Exists,
		"Password":    fields.Exists,
		"Permissions": fields.Exists,
	})
	if len(errs) != 0 {
		return User{}, stream.CombineErrors(errs...)
	}

	password, err := fields.ParsePassword(parameters["Password"])
	if err != nil {
		return User{}, err
	}

	permissions, err := fields.ParsePermissions(fields.CAN_READ, parameters["Permissions"])
	if err != nil {
		return User{}, err
	}

	return User{
		Id:          id,
		Name:        parameters["Name"],
		PictureURI:  parameters["PictureURI"],
		Description: escape.UnescapeString(parameters["Description"]),
		Password:    password,
		Permissions: permissions,
	}, nil
}
