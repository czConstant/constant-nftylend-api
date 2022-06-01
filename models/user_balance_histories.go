package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type UserBalanceHistoryType string

const (
	UserBalanceHistoryTypeBalance       UserBalanceHistoryType = "balance"
	UserBalanceHistoryTypeLockedBalance UserBalanceHistoryType = "locked_balance"
)

type UserBalanceHistory struct {
	gorm.Model
	Network       Network
	UserBalanceID uint
	Type          UserBalanceHistoryType
	CurrencyID    uint
	Currency      *Currency
	Amount        numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	Reference     string
}
