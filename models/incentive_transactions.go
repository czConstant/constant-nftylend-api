package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type IncentiveTransactionStatus string

const (
	IncentiveTransactionStatusPending  IncentiveTransactionStatus = "pending"
	IncentiveTransactionStatusLocked   IncentiveTransactionStatus = "locked"
	IncentiveTransactionStatusUnlocked IncentiveTransactionStatus = "unlocked"
	IncentiveTransactionStatusRevoked  IncentiveTransactionStatus = "revoked"
	IncentiveTransactionStatusDone     IncentiveTransactionStatus = "done"
)

type IncentiveTransaction struct {
	gorm.Model
	Address            string
	Network            Network
	IncentiveProgramID uint
	Type               IncentiveTransactionType
	UserID             uint
	User               *User
	CurrencyID         uint
	Currency           *Currency
	LoanID             uint
	Loan               *Loan
	Amount             numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	LockUntilAt        *time.Time
	UnlockedAt         *time.Time
	Status             IncentiveTransactionStatus
}
