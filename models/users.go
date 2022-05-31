package models

import (
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Network         Network
	Address         string
	AddressChecked  string
	Email           string
	NewsNotiEnabled bool `gorm:"default:0"`
	LoanNotiEnabled bool `gorm:"default:0"`
	SeenNotiID      uint `gorm:"default:0"`
}
