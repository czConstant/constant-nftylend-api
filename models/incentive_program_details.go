package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type IncentiveTransactionType string
type IncentiveTransactionRewardType string

const (
	IncentiveTransactionTypeBorrowerLoanListed        IncentiveTransactionType = "borrower_loan_listed"
	IncentiveTransactionTypeBorrowerLoanDelisted      IncentiveTransactionType = "borrower_loan_delisted"
	IncentiveTransactionTypeLenderLoanMatched         IncentiveTransactionType = "lender_loan_matched"
	IncentiveTransactionTypeUserAirdropReward         IncentiveTransactionType = "user_airdrop_reward"
	IncentiveTransactionTypeUserAmaReward             IncentiveTransactionType = "user_ama_reward"
	IncentiveTransactionTypeAffiliateBorrowerLoanDone IncentiveTransactionType = "affiliate_borrower_loan_done"
	IncentiveTransactionTypeAffiliateLenderLoanDone   IncentiveTransactionType = "affiliate_lender_loan_done"

	IncentiveTransactionRewardTypeAmount     IncentiveTransactionRewardType = "amount"
	IncentiveTransactionRewardTypeRateOfLoan IncentiveTransactionRewardType = "rate_of_loan"
)

type IncentiveProgramDetail struct {
	gorm.Model
	Network            Network
	IncentiveProgramID uint
	IncentiveProgram   *IncentiveProgram
	Type               IncentiveTransactionType
	RewardType         IncentiveTransactionRewardType `gorm:"default:'amount'"`
	Amount             numeric.BigFloat               `gorm:"type:decimal(48,24);default:0"`
	Description        string                         `gorm:"type:text collate utf8mb4_unicode_ci"`
}
