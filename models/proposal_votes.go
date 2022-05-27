package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalVote struct {
	gorm.Model
	Network          Network
	ProposalID       uint
	ProposalChoiceID uint
	Address          string
	Msg              string `gorm:"type:text"`
	Sig              string
	PowerVote        numeric.BigFloat
	Timestamp        *time.Time
	IpfsHash         string
	Status           string
}
