package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type UserBalance struct {
	gorm.Model
	UserID         uint
	Network        Network
	Address        string
	AddressChecked string
	CurrencyID     uint
	Currency       *Currency
	Balance        numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	LockedBalance  numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
}
