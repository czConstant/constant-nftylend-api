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
	Description           string `gorm:"type:text"`
	Creator               string
	ContractAddress       string
	OriginNetwork         Network
	OriginContractAddress string
	Enabled               bool `gorm:"default:0"`
	Verified              bool `gorm:"default:0"`
	ListingAsset          *Asset
	RandAsset             *Asset
	ParasCollectionID     string
}

type NftyRPTListingCollection struct {
	CollectionID uint
	Total        uint
}

type NftyRPTCollectionLoan struct {
	TotalVolume  numeric.BigFloat
	TotalListed  uint
	Avg24hAmount numeric.BigFloat
}

type BorrowerStats struct {
	TotalLoans     uint
	TotalDoneLoans uint
	TotalVolume    numeric.BigFloat
}
