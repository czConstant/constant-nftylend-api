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
	RefUserID          uint
	RefUser            *User
	CurrencyID         uint
	Currency           *Currency
	LoanID             uint
	Loan               *Loan
	Amount             numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	LockUntilAt        *time.Time
	UnlockedAt         *time.Time
	Status             IncentiveTransactionStatus
}

type AffiliateStats struct {
	CurrencyID        uint
	Currency          *Currency
	TotalCommissions  numeric.BigFloat
	TotalUsers        uint
	TotalTransactions uint
}

type AffiliateVolumes struct {
	RptDate          *time.Time
	TotalCommissions numeric.BigFloat
}
