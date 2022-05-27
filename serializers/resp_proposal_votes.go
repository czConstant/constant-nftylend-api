package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type ProposalVoteResp struct {
	ID               uint             `json:"id"`
	CreatedAt        time.Time        `json:"created_at"`
	UpdatedAt        time.Time        `json:"updated_at"`
	Network          models.Network   `json:"network"`
	ProposalID       uint             `json:"proposal_id"`
	ProposalChoiceID uint             `json:"proposal_choice_id"`
	Address          string           `json:"address"`
	PowerVote        numeric.BigFloat `json:"power_vote"`
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
		ProposalChoiceID: m.ProposalChoiceID,
		Address:          m.Address,
		PowerVote:        m.PowerVote,
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
