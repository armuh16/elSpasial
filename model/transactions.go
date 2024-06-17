package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/elspasial/enum"
	"github.com/elspasial/static"
	"gorm.io/gorm"
)

type Transactions struct {
	ID         int
	DriverID   int `json:"-"`
	UserID     int `json:"-"`
	Trip       TripsTransaction
	GrandTotal float64
	Status     enum.TransactionStatusType `json:"-"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `json:"-"`

	// Relations
	User *Users `json:",omitempty" gorm:"<-:false;foreignKey:UserID;references:ID;"`

	// Attribute
	StatusTransaction string `gorm:"<-:false;-;"`
}

type TripsTransaction []TripTransaction

type TripTransaction struct {
	ID          int
	Origin      string
	Destination string
	Price       float64
}

func (j TripsTransaction) Value() (driver.Value, error) {
	return json.Marshal(j)
}

func (j *TripsTransaction) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf(static.SomethingWrong)
	}

	result := TripsTransaction{}
	if err := json.Unmarshal(bytes, &result); err != nil {
		return err
	}

	*j = result

	return nil
}
