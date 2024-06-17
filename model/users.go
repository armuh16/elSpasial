package model

import (
	"time"

	"github.com/elspasial/enum"
	"gorm.io/gorm"
)

type Users struct {
	ID        int
	Name      string
	Email     string
	Password  string `json:"-"`
	Address   string
	RoleID    enum.RoleType `json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `json:"-"`
}
