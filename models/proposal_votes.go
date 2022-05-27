package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalVote struct {
	gorm.Model
	Network          Network
	Type             string
	ProposalID       uint
	ProposalChoiceID uint
	Address          string
	Msg              string `gorm:"type:text"`
	Sig              string
	PowerVote        numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	Timestamp        *time.Time
	IpfsHash         string
	Status           string
}
