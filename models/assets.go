package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Asset struct {
	gorm.Model
	Network             Chain
	CollectionID        uint
	Collection          *Collection
	SeoURL              string
	ContractAddress     string
	TestContractAddress string
	TokenID             uint
	TokenURL            string
	ExternalUrl         string
	Name                string
	Symbol              string
	SellerFeeRate       float64 `gorm:"type:decimal(6,4);default:0"`
	Attributes          string  `gorm:"type:text"`
	MetaJson            string  `gorm:"type:text"`
	MetaJsonUrl         string
	NewLoan             *Loan
	MagicEdenCrawAt     *time.Time
	SolMartCrawAt       *time.Time
	SolSeaCrawAt        *time.Time
}
