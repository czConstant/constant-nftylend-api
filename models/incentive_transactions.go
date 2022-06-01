package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type IncentiveTransactionStatus string

const (
	IncentiveTransactionStatusLocked   IncentiveTransactionStatus = "locked"
	IncentiveTransactionStatusReleased IncentiveTransactionStatus = "released"
	IncentiveTransactionStatusRevoked  IncentiveTransactionStatus = "revoked"
	IncentiveTransactionStatusDone     IncentiveTransactionStatus = "done"
)

type IncentiveTransaction struct {
	gorm.Model
	Network            Network
	IncentiveProgramID uint
	Type               IncentiveTransactionType
	Address            string
	CurrencyID         uint
	LoanID             uint
	Amount             numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	LockUntilAt        *time.Time
	UnlockedAt         *time.Time
	Status             IncentiveTransactionStatus
}
