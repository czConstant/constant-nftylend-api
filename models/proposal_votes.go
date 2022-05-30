package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalVoteStatus string

const (
	ProposalVoteStatusCreated   ProposalVoteStatus = "created"
	ProposalVoteStatusCancelled ProposalVoteStatus = "cancelled"
)

type ProposalVote struct {
	gorm.Model
	Network          Network
	Type             string
	ProposalID       uint
	Proposal         *Proposal
	ProposalChoiceID uint
	ProposalChoice   *ProposalChoice
	Address          string
	Msg              string `gorm:"type:text"`
	Sig              string
	PowerVote        numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	Timestamp        *time.Time
	IpfsHash         string
	CancelledHash    string
	Status           ProposalVoteStatus
}
