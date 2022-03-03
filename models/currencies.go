package models

import "github.com/jinzhu/gorm"

type Currency struct {
	gorm.Model
	Network         Chain
	ContractAddress string
	Decimals        uint
	Symbol          string
	Name            string
	IconURL         string
	AdminFeeAddress string
	Enabled         float64 `gorm:"default:0"`
}
