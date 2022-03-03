package models

import (
	"github.com/czConstant/constant-nftlend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type Collection struct {
	gorm.Model
	Network      Chain
	SeoURL       string
	Name         string
	Description  string
	ListingAsset *Asset
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
