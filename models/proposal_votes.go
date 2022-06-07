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
	UserID           uint
	User             *User
	Type             string
	ProposalID       uint
	Proposal         *Proposal
	ProposalChoiceID uint
	ProposalChoice   *ProposalChoice
	Message          string `gorm:"type:text"`
	Signature        string
	PowerVote        numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	Timestamp        *time.Time
	IpfsHash         string
	CancelledHash    string
	Status           ProposalVoteStatus
}
