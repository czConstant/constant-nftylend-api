package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type Collection struct {
	gorm.Model
	Network               Network
	SeoURL                string
	Name                  string
	Description           string `gorm:"type:text collate utf8mb4_unicode_ci"`
	Creator               string
	ContractAddress       string
	OriginNetwork         Network
	OriginContractAddress string
	Enabled               bool `gorm:"default:0"`
	Verified              bool `gorm:"default:0"`
	ParasCollectionID     string
	CreatorURL            string
	TwitterID             string
	DiscordURL            string
	CoverURL              string
	ImageURL              string
	VolumeUsd             numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	FloorPrice            numeric.BigFloat `gorm:"type:decimal(48,24);default:0"`
	CurrencyID            uint
	Currency              *Currency
	NewLoanID             uint
	NewLoan               *Loan `gorm:"foreignKey:new_loan_id"`
}

type NftyRPTListingCollection struct {
	CollectionID uint
	Total        uint
}

type CollectionStats struct {
	TotalVolume  numeric.BigFloat
	TotalListed  uint
	Avg24hAmount numeric.BigFloat
	MinAmount    numeric.BigFloat
}

type BorrowerStats struct {
	TotalLoans     uint
	TotalDoneLoans uint
	TotalVolume    numeric.BigFloat
}

type PlatformStats struct {
	TotalLoans          uint
	TotalDoneLoans      uint
	TotalDefaultedLoans uint
	TotalVolume         numeric.BigFloat
}
