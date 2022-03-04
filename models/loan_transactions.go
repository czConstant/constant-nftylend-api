package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type LoanTransactionType string

const (
	LoanTransactionTypeListed     LoanTransactionType = "listed"
	LoanTransactionTypeCancelled  LoanTransactionType = "cancelled"
	LoanTransactionTypeOffered    LoanTransactionType = "offered"
	LoanTransactionTypeRepaid     LoanTransactionType = "repaid"
	LoanTransactionTypeLiquidated LoanTransactionType = "liquidated"
)

type LoanTransaction struct {
	gorm.Model
	Network         Chain
	LoanID          uint
	Loan            *Loan
	Type            LoanTransactionType
	Borrower        string
	Lender          string
	StartedAt       *time.Time
	Duration        uint
	ExpiredAt       *time.Time
	PrincipalAmount numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	InterestRate    float64          `gorm:"type:decimal(6,4);default:0"`
	TxHash          string
}
