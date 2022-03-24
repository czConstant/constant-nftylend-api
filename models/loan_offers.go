package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type LoanOfferStatus string

const (
	LoanOfferStatusNew        LoanOfferStatus = "new"
	LoanOfferStatusApproved   LoanOfferStatus = "approved"
	LoanOfferStatusCancelled  LoanOfferStatus = "cancelled"
	LoanOfferStatusRejected   LoanOfferStatus = "rejected"
	LoanOfferStatusRepaid     LoanOfferStatus = "repaid"
	LoanOfferStatusLiquidated LoanOfferStatus = "liquidated"
	LoanOfferStatusDone       LoanOfferStatus = "done"
	LoanOfferStatusExpired    LoanOfferStatus = "expired"
)

type LoanOffer struct {
	gorm.Model
	Network             Network
	LoanID              uint
	Loan                *Loan
	Lender              string
	StartedAt           *time.Time
	Duration            uint
	ExpiredAt           *time.Time
	FinishedAt          *time.Time
	PrincipalAmount     numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	InterestRate        float64          `gorm:"type:decimal(6,4);default:0"`
	NonceHex            string
	Signature           string
	Status              LoanOfferStatus
	DataOfferAddress    string
	DataCurrencyAddress string
	RepaidAt            *time.Time
	RepaidAmount        numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	MakeTxHash          string
	AcceptTxHash        string
	CancelTxHash        string
	CloseTxHash         string
}
