package models

import (
	"math"
	"math/big"
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type Network string
type LoanStatus string

const (
	LoanStatusNew        LoanStatus = "new"
	LoanStatusCreated    LoanStatus = "created"
	LoanStatusCancelled  LoanStatus = "cancelled"
	LoanStatusDone       LoanStatus = "done"
	LoanStatusLiquidated LoanStatus = "liquidated"
	LoanStatusExpired    LoanStatus = "expired"

	NetworkSOL    Network = "SOL"
	NetworkMATIC  Network = "MATIC"
	NetworkETH    Network = "ETH"
	NetworkAVAX   Network = "AVAX"
	NetworkBSC    Network = "BSC"
	NetworkBOBA   Network = "BOBA"
	NetworkNEAR   Network = "NEAR"
	NetworkAURORA Network = "AURORA"
)

type Loan struct {
	gorm.Model
	Network              Network
	Owner                string
	BorrowerUserID       uint
	BorrowerUser         *User
	Lender               string
	LenderUserID         uint
	LenderUser           *User
	AssetID              uint
	Asset                *Asset
	CollectionID         uint
	Collection           *Collection
	CurrencyID           uint `gorm:"default:0"`
	Currency             *Currency
	CurrencyPrice        float64 `gorm:"type:decimal(16,8);default:0"`
	StartedAt            *time.Time
	Duration             uint `gorm:"default:0"`
	ExpiredAt            *time.Time
	FinishedAt           *time.Time
	PrincipalAmount      numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	InterestRate         float64          `gorm:"type:decimal(6,4);default:0"`
	ValidAt              *time.Time
	Config               uint `gorm:"default:0"`
	OfferStartedAt       *time.Time
	OfferDuration        uint `gorm:"default:0"`
	OfferExpiredAt       *time.Time
	OfferPrincipalAmount numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	OfferInterestRate    float64          `gorm:"type:decimal(6,4);default:0"`
	FeeRate              float64          `gorm:"type:decimal(6,4);default:0"`
	FeeAmount            numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	OfferOverdueAt       *time.Time
	NonceHex             string
	ImageUrl             string
	Signature            string
	Status               LoanStatus
	DataLoanAddress      string
	DataAssetAddress     string
	Offers               []*LoanOffer
	ApprovedOffer        *LoanOffer
	RepaidAmount         numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	InitTxHash           string
	CancelTxHash         string
	PayTxHash            string
	LiquidateTxHash      string
	SynchronizedAt       *time.Time
}

func (m *Loan) MaturedOfferPaymentAmount() *big.Float {
	var amount big.Float
	if m.OfferStartedAt != nil {
		amount = m.OfferPrincipalAmount.Float
		duration := float64(m.OfferDuration) / (24 * 60 * 60)
		amount = *MulBigFloats(&amount, big.NewFloat(m.InterestRate), big.NewFloat(duration/365))
	}
	amount = *AddBigFloats(&amount, &m.OfferPrincipalAmount.Float)
	return &amount
}

func (m *Loan) EarlyOfferPaymentAmount() *big.Float {
	var amount big.Float
	if m.OfferStartedAt != nil {
		amount = m.OfferPrincipalAmount.Float
		duration := float64(m.OfferDuration) / (24 * 60 * 60)
		if time.Now().Before(*m.OfferExpiredAt) {
			duration = math.Ceil(float64(time.Since(*m.OfferStartedAt)) / float64(24*time.Hour))
		}
		amount = *MulBigFloats(&amount, big.NewFloat(m.InterestRate), big.NewFloat(duration/365))
	}
	amount = *AddBigFloats(&amount, &m.OfferPrincipalAmount.Float)
	amount = *QuoBigFloats(AddBigFloats(&amount, m.MaturedOfferPaymentAmount()), big.NewFloat(2))
	return &amount
}
