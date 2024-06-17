package dto

import (
	"fmt"
	"github.com/elspasial/enum"
	"github.com/elspasial/model"
	"github.com/elspasial/static"
)

type FindRequest model.Users

type CreateRequest struct {
	Name     string
	Password string
	Email    string
	Adress   string
	RoleID   enum.RoleType
}

func (d *CreateRequest) Validate() error {
	if d.Name == "" {
		return fmt.Errorf(static.EmptyValue, "Nama Lengkap")
	}
	if d.Email == "" {
		return fmt.Errorf(static.EmptyValue, "Email")
	}
	if d.Password == "" {
		return fmt.Errorf(static.EmptyValue, "Password")
	}
	if !static.EmailRegex.MatchString(d.Email) {
		return fmt.Errorf(static.ValueNotValid, "Email")
	}
	if d.Adress == "" {
		return fmt.Errorf(static.EmptyValue, "Alamat")

	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	return nil
}
