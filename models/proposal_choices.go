package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalChoice struct {
	gorm.Model
	Network    Network
	ProposalID uint
	Choice     int
	Name       string
	PowerVote  numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
}
