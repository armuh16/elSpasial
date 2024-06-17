package dto

import (
	"fmt"
	"github.com/elspasial/enum"

	"github.com/elspasial/static"
)

//	type RegisterRequest struct {
//		Name     string
//		Email    string
//		Password string
//		Address   string
//		RoleID   enum.RoleType
//	}

// RegisterRequest represents the request body for user registration
// @Description Request body for user registration
type RegisterRequest struct {
	Name     string        `json:"name" example:"John Doe"`
	Email    string        `json:"email" example:"john.doe@example.com"`
	Password string        `json:"password" example:"secret"`
	Address  string        `json:"address" example:"123 Main St"`
	RoleID   enum.RoleType `json:"role_id" example:"2"`
}

func (d *RegisterRequest) Validate() error {
	if d.Name == "" {
		return fmt.Errorf(static.EmptyValue, "Nama Lengkap")
	}
	if d.Email == "" {
		return fmt.Errorf(static.EmptyValue, "Email")
	}
	if !static.EmailRegex.MatchString(d.Email) {
		return fmt.Errorf(static.ValueNotValid, "Email")
	}
	if d.Password == "" {
		return fmt.Errorf(static.EmptyValue, "Kata Sandi")
	}
	if d.Address == "" {
		return fmt.Errorf(static.EmptyValue, "Alamat")
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	return nil
}

//type LoginRequest struct {
//	Email    string
//	Password string
//}

// LoginRequest represents the request body for user login
// @Description Request body for user login
type LoginRequest struct {
	Email    string `json:"email" example:"john.doe@example.com"`
	Password string `json:"password" example:"secret"`
}

func (d *LoginRequest) Validate() error {
	if d.Email == "" {
		return fmt.Errorf(static.EmptyValue, "email")
	}
	if d.Password == "" {
		return fmt.Errorf(static.EmptyValue, "password")
	}
	return nil
}

//type Response struct {
//	Token        string
//	RefreshToken string
//}

// Response represents the response body for authentication
// @Description Response body for authentication
type Response struct {
	Token        string `json:"token" example:"jwt-token"`
	RefreshToken string `json:"refresh_token" example:"refresh-jwt-token"`
}
