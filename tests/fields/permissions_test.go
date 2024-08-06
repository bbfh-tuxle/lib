package fields_test

import (
	"testing"

	"github.com/bbfh-tuxle/lib/tuxle/fields"
	"gotest.tools/assert"
)

func TestPermissionDefault(t *testing.T) {
	permissions := fields.EmptyPermissions(fields.CAN_READ)
	assert.DeepEqual(t, permissions.Server, fields.CAN_READ)
	assert.DeepEqual(t, permissions.Channels, fields.CAN_READ)
	assert.DeepEqual(t, permissions.Messages, fields.CAN_READ)
	assert.DeepEqual(t, permissions.Users, fields.CAN_READ)
}

func TestPermissionCan(t *testing.T) {
	permissions := fields.EmptyPermissions(fields.CAN_INTERACT)
	if !permissions.Server.Can(fields.CAN_READ) {
		t.Fatal("Should have permission")
	}
	if !permissions.Server.Can(fields.CAN_INTERACT) {
		t.Fatal("Should have permission")
	}
	if permissions.Server.Can(fields.CAN_MODIFY) {
		t.Fatal("Should not have permission")
	}
}

func TestPermissionParse(t *testing.T) {
	permissions, err := fields.ParsePermissions(fields.CAN_NOTHING, "server:rim*")
	if err != nil {
		t.Fatal(err)
	}

	if !permissions.Server.Can(fields.CAN_ALL) {
		t.Fatal("Should have permission")
	}
	if permissions.Users.Can(fields.CAN_READ) {
		t.Fatal("Should not have permission")
	}

	got, err := fields.ParsePermissions(fields.CAN_NOTHING, permissions.String())
	assert.DeepEqual(t, got, permissions)
}
