package models

import (
	"github.com/jinzhu/gorm"
)

type Currency struct {
	gorm.Model
	Network         Network
	ContractAddress string
	Decimals        uint
	Symbol          string
	Name            string
	IconURL         string
	AdminFeeAddress string
	Enabled         float64 `gorm:"default:0"`
	Price           float64 `gorm:"type:decimal(16,8);default:0"`
	WithdrawEnabled bool    `gorm:"default:0"`
}
