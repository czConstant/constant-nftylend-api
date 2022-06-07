package models

import (
	"math"
	"math/big"
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
	LenderUserID        uint
	LenderUser          *User
	StartedAt           *time.Time
	Duration            uint
	ExpiredAt           *time.Time
	FinishedAt          *time.Time
	PrincipalAmount     numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	InterestRate        float64          `gorm:"type:decimal(6,4);default:0"`
	ValidAt             *time.Time
	NonceHex            string
	Signature           string
	Status              LoanOfferStatus
	DataOfferAddress    string
	DataCurrencyAddress string
	RepaidAt            *time.Time
	RepaidAmount        numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	MakeTxHash          string
	AcceptTxHash        string
	CancelTxHash        string
	CloseTxHash         string
}

func (m *LoanOffer) MaturedOfferPaymentAmount() *big.Float {
	var amount big.Float
	amount = m.PrincipalAmount.Float
	duration := float64(m.Duration) / (24 * 60 * 60)
	amount = *MulBigFloats(&amount, big.NewFloat(m.InterestRate), big.NewFloat(duration/365))
	amount = *AddBigFloats(&amount, &m.PrincipalAmount.Float)
	return &amount
}

func (m *LoanOffer) EarlyOfferPaymentAmount() *big.Float {
	var amount big.Float
	if m.StartedAt != nil {
		amount = m.PrincipalAmount.Float
		duration := float64(m.Duration) / (24 * 60 * 60)
		if time.Now().Before(*m.ExpiredAt) {
			duration = math.Ceil(float64(time.Since(*m.StartedAt)) / float64(24*time.Hour))
		}
		amount = *MulBigFloats(&amount, big.NewFloat(m.InterestRate), big.NewFloat(duration/365))
	}
	amount = *AddBigFloats(&amount, &m.PrincipalAmount.Float)
	amount = *QuoBigFloats(AddBigFloats(&amount, m.MaturedOfferPaymentAmount()), big.NewFloat(2))
	return &amount
}
