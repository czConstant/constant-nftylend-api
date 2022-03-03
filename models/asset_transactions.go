package models

import (
	"time"

	"github.com/czConstant/constant-nftlend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type AssetTransactionType string

const (
	AssetTransactionTypeExchange AssetTransactionType = "exchange"
)

type AssetTransaction struct {
	gorm.Model
	Network       Chain
	AssetID       uint
	Asset         *Asset
	Type          AssetTransactionType
	Seller        string
	Buyer         string
	TransactionAt *time.Time
	Amount        numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	CurrencyID    uint
	Currency      *Currency
}