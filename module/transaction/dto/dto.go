package dto

import (
	"errors"
	"fmt"

	"github.com/elspasial/enum"
	"github.com/elspasial/static"
)

//type CreateOrderRequest struct {
//	UserID int
//	RoleID enum.RoleType
//	Trip   []int
//}

// CreateOrderRequest represents the request body for creating an order
// @Description Request body for creating an order
type CreateOrderRequest struct {
	UserID int           `json:"-"`
	RoleID enum.RoleType `json:"-"`
	Trip   []int         `json:"trip" example:"[1, 2, 3]"`
}

func (d *CreateOrderRequest) Validate() error {
	if d.UserID <= 0 {
		return fmt.Errorf(static.EmptyValue, "UserID")
	}
	if len(d.Trip) == 0 {
		return fmt.Errorf(static.EmptyValue, "Trip")
	}
	if d.RoleID != enum.RoleTypeUser {
		return errors.New(static.Authorization)
	}
	return nil
}

//type FindAllRequest struct {
//	UserID int
//	Trip   []int
//	RoleID enum.RoleType
//}

// FindAllRequest represents the request body for finding transactions
// @Description Request body for finding transactions
type FindAllRequest struct {
	UserID int           `json:"-"`
	Trip   []int         `json:"-"`
	RoleID enum.RoleType `json:"-"`
}

func (d *FindAllRequest) Validate() error {
	if d.UserID <= 0 {
		return fmt.Errorf(static.EmptyValue, "UserID")
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	return nil
}

//type AcceptOrderRequest struct {
//	DriverID      int
//	TransactionID int
//	RoleID        enum.RoleType
//}

// AcceptOrderRequest represents the request body for accepting an order
// @Description Request body for accepting an order
type AcceptOrderRequest struct {
	DriverID      int           `json:"-"`
	TransactionID int           `json:"transaction_id" example:"1"`
	RoleID        enum.RoleType `json:"-"`
}

func (d *AcceptOrderRequest) Validate() error {
	if d.DriverID <= 0 {
		return fmt.Errorf(static.EmptyValue, "DriverID")
	}
	if d.TransactionID <= 0 {
		return fmt.Errorf(static.EmptyValue, "TransactionID")
	}
	if err := d.RoleID.IsValid(); err != nil {
		return err
	}
	if d.RoleID != enum.RoleTypeDriver {
		return errors.New(static.Authorization)
	}
	return nil
}
