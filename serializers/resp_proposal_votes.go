package serializers

import (
	"fmt"
	"time"

	"github.com/czConstant/constant-nftylend-api/configs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type ProposalVoteResp struct {
	ID               uint                      `json:"id"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
	Network          models.Network            `json:"network"`
	UserID           uint                      `json:"user_id"`
	User             *UserResp                 `json:"user"`
	ProposalID       uint                      `json:"proposal_id"`
	Proposal         *ProposalResp             `json:"proposal"`
	ProposalChoiceID uint                      `json:"proposal_choice_id"`
	ProposalChoice   *ProposalChoiceResp       `json:"proposal_choice"`
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
		UserID:           m.UserID,
		User:             NewMiniUserResp(m.User),
		ProposalID:       m.ProposalID,
		Proposal:         NewProposalResp(m.Proposal),
		ProposalChoiceID: m.ProposalChoiceID,
		ProposalChoice:   NewProposalChoiceResp(m.ProposalChoice),
		PowerVote:        m.PowerVote,
		Timestamp:        m.Timestamp,
		IpfsHash:         fmt.Sprintf("%s/api/ipfs/%s", configs.GetConfig().WebUrl, m.IpfsHash),
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
