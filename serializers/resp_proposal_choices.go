package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type ProposalChoiceResp struct {
	ID         uint             `json:"id"`
	CreatedAt  time.Time        `json:"created_at"`
	UpdatedAt  time.Time        `json:"updated_at"`
	Network    models.Network   `json:"network"`
	ProposalID uint             `json:"proposal_id"`
	Choice     int              `json:"choice"`
	Name       string           `json:"name"`
	PowerVote  numeric.BigFloat `json:"power_vote"`
}

func NewProposalChoiceResp(m *models.ProposalChoice) *ProposalChoiceResp {
	if m == nil {
		return nil
	}
	resp := &ProposalChoiceResp{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Network:    m.Network,
		ProposalID: m.ProposalID,
		Choice:     m.Choice,
		Name:       m.Name,
		PowerVote:  m.PowerVote,
	}
	return resp
}

func NewProposalChoiceRespArr(arr []*models.ProposalChoice) []*ProposalChoiceResp {
	resps := []*ProposalChoiceResp{}
	for _, m := range arr {
		resps = append(resps, NewProposalChoiceResp(m))
	}
	return resps
}
