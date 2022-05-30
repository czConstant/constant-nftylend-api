package models

import (
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalChoiceStatus string

const (
	ProposalChoiceStatusCreated   ProposalChoiceStatus = "created"
	ProposalChoiceStatusCancelled ProposalChoiceStatus = "cancelled"
	ProposalChoiceStatusSucceeded ProposalChoiceStatus = "succeeded"
	ProposalChoiceStatusDefeated  ProposalChoiceStatus = "defeated"
	ProposalChoiceStatusQueued    ProposalChoiceStatus = "queued"
	ProposalChoiceStatusExecuted  ProposalChoiceStatus = "executed"
)

type ProposalChoice struct {
	gorm.Model
	Network    Network
	ProposalID uint
	Proposal   *Proposal
	Choice     int
	Name       string
	PowerVote  numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
}
