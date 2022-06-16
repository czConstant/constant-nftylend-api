package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type Currency struct {
	gorm.Model
	Network               Network
	ContractAddress       string
	Decimals              uint
	Symbol                string
	Name                  string
	IconURL               string
	AdminFeeAddress       string
	Enabled               float64 `gorm:"default:0"`
	Price                 float64 `gorm:"type:decimal(16,8);default:0"`
	PoolAddress           string
	ClaimEnabled          bool             `gorm:"default:0"`
	ProposalThreshold     numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	ProposalPowerRequired numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	ProposalPwpRequired   numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
}
