package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Network         Network
	Address         string
	AddressChecked  string
	UserName        string
	Email           string
	NewsNotiEnabled bool `gorm:"default:0"`
	LoanNotiEnabled bool `gorm:"default:0"`
	SeenNotiID      uint `gorm:"default:0"`
	ReferrerCode    string
	ReferrerUserID  uint
}

type UserBorrowStats struct {
	TotalLoans  uint             `json:"total_loans"`
	TotalVolume numeric.BigFloat `json:"total_volume"`
}

type UserLendStats struct {
	TotalLoans  uint             `json:"total_loans"`
	TotalVolume numeric.BigFloat `json:"total_volume"`
}
