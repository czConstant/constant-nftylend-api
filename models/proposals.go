package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalStatus string

const (
	ProposalTypeProposal = "proposal"

	ProposalStatusCreated   ProposalStatus = "created"
	ProposalStatusCancelled ProposalStatus = "cancelled"
	ProposalStatusSucceeded ProposalStatus = "succeeded"
	ProposalStatusDefeated  ProposalStatus = "defeated"
	ProposalStatusQueued    ProposalStatus = "queued"
	ProposalStatusExecuted  ProposalStatus = "executed"

	ProposalChoiceTypeSingleChoice   = "single-choice"
	ProposalChoiceTypeMultipleChoice = "multiple-choice"
)

type Proposal struct {
	gorm.Model
	Network           Network
	Address           string
	Type              string
	ChoiceType        string
	Msg               string `gorm:"type:text"`
	Sig               string
	Snapshot          int64
	Name              string `gorm:"type:text"`
	Body              string `gorm:"type:text"`
	Timestamp         *time.Time
	Start             *time.Time
	End               *time.Time
	IpfsHash          string
	TotalVote         numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	ProposalThreshold numeric.BigFloat `gorm:"type:decimal(36,18);default:0"`
	Status            ProposalStatus
	Choices           []*ProposalChoice
}
