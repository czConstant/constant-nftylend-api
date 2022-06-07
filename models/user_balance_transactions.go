package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type UserBalanceTransactionType string
type UserBalanceTransactionStatus string

const (
	UserBalanceTransactionTypeDeposit   UserBalanceTransactionType = "deposit"
	UserBalanceTransactionTypeClaim     UserBalanceTransactionType = "claim"
	UserBalanceTransactionTypeIncentive UserBalanceTransactionType = "incentive"

	UserBalanceTransactionStatusPending UserBalanceTransactionStatus = "pending"
	UserBalanceTransactionStatusDone    UserBalanceTransactionStatus = "done"
)

type UserBalanceTransaction struct {
	gorm.Model
	Network                Network
	UserID                 uint
	User                   *User
	UserBalanceID          uint
	UserBalance            *UserBalance
	Type                   UserBalanceTransactionType
	CurrencyID             uint
	Currency               *Currency
	ToAddress              string
	Amount                 numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	Signature              string
	TxHash                 string
	Status                 UserBalanceTransactionStatus
	IncentiveTransactionID uint
	IncentiveTransaction   *IncentiveTransaction
}
