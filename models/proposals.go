package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

const (
	ProposalTypeProposal = "proposal"

	ProposalStatusCreated   = "created"
	ProposalStatusCancelled = "cancelled"
	ProposalStatusSucceeded = "succeeded"
	ProposalStatusDefeated  = "defeated"
	ProposalStatusQueued    = "queued"
	ProposalStatusExecuted  = "executed"

	ProposalChoiceTypeSingleChoice   = "single-choice"
	ProposalChoiceTypeMultipleChoice = "multiple-choice"
)

type Proposal struct {
	gorm.Model
	Network    Network
	Address    string
	Type       string
	ChoiceType string
	Msg        string `gorm:"type:text"`
	Sig        string
	Snapshot   int64
	Name       string `gorm:"type:text"`
	Body       string `gorm:"type:text"`
	Timestamp  *time.Time
	Start      *time.Time
	End        *time.Time
	IpfsHash   string
	Status     string
}
