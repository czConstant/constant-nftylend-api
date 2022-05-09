package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Network Network
	Address string
	Email   string
}
