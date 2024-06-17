package enum

import (
	"fmt"

	"github.com/elspasial/static"
)

type RoleType int

const (
	RoleTypeDriver RoleType = 1
	RoleTypeUser   RoleType = 2
)

func (t RoleType) String() string {
	switch t {
	case RoleTypeDriver:
		return "Driver"
	case RoleTypeUser:
		return "User"
	default:
		return "Unknown"
	}
}

func (t RoleType) IsValid() error {
	switch t {
	case RoleTypeDriver, RoleTypeUser:
		return nil
	}
	return fmt.Errorf(static.DataNotFound, "Role")
}
