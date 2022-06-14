package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type IncentiveTransactionType string

const (
	IncentiveTransactionTypeBorrowerLoanListed   IncentiveTransactionType = "borrower_loan_listed"
	IncentiveTransactionTypeBorrowerLoanDelisted IncentiveTransactionType = "borrower_loan_delisted"
	IncentiveTransactionTypeLenderLoanMatched    IncentiveTransactionType = "lender_loan_matched"
)

type IncentiveProgramDetail struct {
	gorm.Model
	Network            Network
	IncentiveProgramID uint
	IncentiveProgram   *IncentiveProgram
	Type               IncentiveTransactionType
	Amount             numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	Description        string           `gorm:"type:text"`
}
