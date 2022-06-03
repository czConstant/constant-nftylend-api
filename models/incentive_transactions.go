package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type IncentiveTransactionStatus string

const (
	IncentiveTransactionStatusLocked   IncentiveTransactionStatus = "locked"
	IncentiveTransactionStatusUnlocked IncentiveTransactionStatus = "unlocked"
	IncentiveTransactionStatusRevoked  IncentiveTransactionStatus = "revoked"
	IncentiveTransactionStatusDone     IncentiveTransactionStatus = "done"
)

type IncentiveTransaction struct {
	gorm.Model
	Network            Network
	IncentiveProgramID uint
	Type               IncentiveTransactionType
	UserID             uint
	CurrencyID         uint
	LoanID             uint
	Amount             numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	LockUntilAt        *time.Time
	UnlockedAt         *time.Time
	Status             IncentiveTransactionStatus
}
