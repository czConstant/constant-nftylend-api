package models

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type ProposalStatus string
type ProposalChoiceType string

const (
	ProposalTypeProposal = "proposal"

	ProposalStatusPending   ProposalStatus = "pending"
	ProposalStatusCreated   ProposalStatus = "created"
	ProposalStatusCancelled ProposalStatus = "cancelled"
	ProposalStatusSucceeded ProposalStatus = "succeeded"
	ProposalStatusDefeated  ProposalStatus = "defeated"
	ProposalStatusQueued    ProposalStatus = "queued"
	ProposalStatusExecuted  ProposalStatus = "executed"

	ProposalChoiceTypeSingleChoice   ProposalChoiceType = "single-choice"
	ProposalChoiceTypeMultipleChoice ProposalChoiceType = "multiple-choice"
)

type Proposal struct {
	gorm.Model
	Network           Network
	UserID            uint
	User              *User
	Type              string
	ChoiceType        ProposalChoiceType
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
