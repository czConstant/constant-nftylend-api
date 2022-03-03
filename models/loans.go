package models

import (
	"time"

	"github.com/czConstant/constant-nftlend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type Chain string
type LoanStatus string

const (
	LoanStatusNew        LoanStatus = "new"
	LoanStatusCreated    LoanStatus = "created"
	LoanStatusCancelled  LoanStatus = "cancelled"
	LoanStatusDone       LoanStatus = "done"
	LoanStatusLiquidated LoanStatus = "liquidated"
	LoanStatusExpired    LoanStatus = "expired"

	ChainSOL   Chain = "SOL"
	ChainMATIC Chain = "MATIC"
	ChainETH   Chain = "ETH"
)

type Loan struct {
	gorm.Model
	Network              Chain
	Owner                string
	Lender               string
	AssetID              uint
	Asset                *Asset
	CurrencyID           uint `gorm:"default:0"`
	Currency             *Currency
	StartedAt            *time.Time
	Duration             uint `gorm:"default:0"`
	ExpiredAt            *time.Time
	FinishedAt           *time.Time
	PrincipalAmount      numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	InterestRate         float64          `gorm:"type:decimal(6,4);default:0"`
	OfferStartedAt       *time.Time
	OfferDuration        uint `gorm:"default:0"`
	OfferExpiredAt       *time.Time
	OfferPrincipalAmount numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	OfferInterestRate    float64          `gorm:"type:decimal(6,4);default:0"`
	FeeRate              float64          `gorm:"type:decimal(6,4);default:0"`
	FeeAmount            numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	NonceHex             string
	ImageUrl             string
	Signature            string
	Status               LoanStatus
	DataLoanAddress      string
	DataAssetAddress     string
	Offers               []*LoanOffer
	ApprovedOffer        *LoanOffer
	RepaidAmount         numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	InitTxHash           string
	CancelTxHash         string
	PayTxHash            string
	LiquidateTxHash      string
}
