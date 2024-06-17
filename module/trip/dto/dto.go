package dto

import (
	"errors"
	"fmt"

	"github.com/elspasial/enum"
	"github.com/elspasial/model"
	"github.com/elspasial/static"
)

//type CreateRequest struct {
//	Origin      string
//	Destination string
//	UserID      int
//	Price       float64
//	RoleID      enum.RoleType
//}

// CreateRequest represents the request body for creating a trip
// @Description Request body for creating a trip
type CreateRequest struct {
	Origin      string        `json:"origin" example:"City A"`
	Destination string        `json:"destination" example:"City B"`
	UserID      int           `json:"-"`
	Price       float64       `json:"price" example:"100.000"`
	RoleID      enum.RoleType `json:"-"`
}

func (d *CreateRequest) Validate() error {
	if d.UserID <= 0 {
		return fmt.Errorf(static.EmptyValue, "SellerID")
	}
	if d.Origin == "" {
		return fmt.Errorf(static.EmptyValue, "Description")
	}
	if d.Destination == "" {
		return fmt.Errorf(static.EmptyValue, "Description")
	}
	if d.Price <= 0 {
		return fmt.Errorf(static.MinValue, "Price", 0)
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	if d.RoleID != enum.RoleTypeUser {
		return errors.New(static.Authorization)
	}
	return nil
}

//type FindAllRequest struct {
//	UserID int
//}

// FindAllRequest represents the request body for finding all trips
// @Description Request body for finding all trips
type FindAllRequest struct {
	UserID int `json:"-"`
}

func (d *FindAllRequest) Validate() error {
	if d.UserID <= 0 {
		return fmt.Errorf(static.EmptyValue, "UserID")
	}
	return nil
}

type FindRequest model.Trips
