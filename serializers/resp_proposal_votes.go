package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type ProposalVoteResp struct {
	ID               uint                      `json:"id"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
	Network          models.Network            `json:"network"`
	ProposalID       uint                      `json:"proposal_id"`
	Proposal         *ProposalResp             `json:"proposal"`
	ProposalChoiceID uint                      `json:"proposal_choice_id"`
	ProposalChoice   *ProposalChoiceResp       `json:"proposal_choice"`
	Address          string                    `json:"address"`
	PowerVote        numeric.BigFloat          `json:"power_vote"`
	Timestamp        *time.Time                `json:"timestamp"`
	IpfsHash         string                    `json:"ipfs_hash"`
	Status           models.ProposalVoteStatus `json:"status"`
}

func NewProposalVoteResp(m *models.ProposalVote) *ProposalVoteResp {
	if m == nil {
		return nil
	}
	resp := &ProposalVoteResp{
		ID:               m.ID,
		CreatedAt:        m.CreatedAt,
		UpdatedAt:        m.UpdatedAt,
		Network:          m.Network,
		ProposalID:       m.ProposalID,
		Proposal:         NewProposalResp(m.Proposal),
		ProposalChoiceID: m.ProposalChoiceID,
		ProposalChoice:   NewProposalChoiceResp(m.ProposalChoice),
		Address:          m.Address,
		PowerVote:        m.PowerVote,
		Timestamp:        m.Timestamp,
		IpfsHash:         m.IpfsHash,
		Status:           m.Status,
	}
	return resp
}

func NewProposalVoteRespArr(arr []*models.ProposalVote) []*ProposalVoteResp {
	resps := []*ProposalVoteResp{}
	for _, m := range arr {
		resps = append(resps, NewProposalVoteResp(m))
	}
	return resps
}
