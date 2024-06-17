package model

import (
	"time"

	"gorm.io/gorm"
)

type Trips struct {
	ID          int
	Origin      string
	Destination string
	UserID      int
	Price       float64

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"-"`

	// Relations
	User *Users `json:",omitempty" gorm:"<-:false;foreignKey:UserID;references:ID;"`
}
