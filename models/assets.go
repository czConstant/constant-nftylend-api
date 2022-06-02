package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type Asset struct {
	gorm.Model
	Network                   Network
	CollectionID              uint
	Collection                *Collection
	SeoURL                    string
	ContractAddress           string
	TokenID                   string
	TestContractAddress       string
	TestTokenID               string
	TokenURL                  string `gorm:"type:text"`
	ExternalUrl               string `gorm:"type:text"`
	Name                      string
	Symbol                    string
	Description               string `gorm:"type:text"`
	SearchText                string `gorm:"type:text collate utf8mb4_unicode_ci"`
	MimeType                  string
	SellerFeeRate             float64 `gorm:"type:decimal(6,4);default:0"`
	Attributes                string  `gorm:"type:text"`
	MetaJson                  string  `gorm:"type:text"`
	MetaJsonUrl               string
	NewLoan                   *Loan
	OriginNetwork             Network
	OriginContractAddress     string
	OriginTokenID             string
	TestOriginContractAddress string
	TestOriginTokenID         uint
	MagicEdenCrawAt           *time.Time
	SolanartCrawAt            *time.Time
	SolSeaCrawAt              *time.Time
	OpenseaCrawAt             *time.Time
	ParasCrawAt               *time.Time
	NftbankCrawAt             *time.Time
	FloorPriceAt              *time.Time
	FloorPrice                numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
}

func (m *Asset) GetContractAddress() string {
	if m.TestContractAddress != "" {
		return m.TestContractAddress
	}
	return m.ContractAddress
}

func (m *Asset) GetTokenID() string {
	if m.TestTokenID != "" {
		return m.TestTokenID
	}
	return m.TokenID
}
