package serializers

import (
	"fmt"
	"time"

	"github.com/czConstant/constant-nftylend-api/configs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type ProposalResp struct {
	ID                uint                      `json:"id"`
	CreatedAt         time.Time                 `json:"created_at"`
	UpdatedAt         time.Time                 `json:"updated_at"`
	Network           models.Network            `json:"network"`
	UserID            uint                      `json:"user_id"`
	User              *UserResp                 `json:"user"`
	Type              models.ProposalType       `json:"type"`
	ChoiceType        models.ProposalChoiceType `json:"choice_type"`
	Message           string                    `json:"message"`
	Signature         string                    `json:"signature"`
	Snapshot          int64                     `json:"snapshot"`
	Name              string                    `json:"name"`
	Body              string                    `json:"body"`
	ProjectName       string                    `json:"project_name"`
	Contact           string                    `json:"contact"`
	Timestamp         *time.Time                `json:"timestamp"`
	Start             *time.Time                `json:"start"`
	End               *time.Time                `json:"end"`
	IpfsHash          string                    `json:"ipfs_hash"`
	ProposalThreshold numeric.BigFloat          `json:"proposal_threshold"`
	Status            models.ProposalStatus     `json:"status"`
	Choices           []*ProposalChoiceResp     `json:"choices"`
}

func NewProposalResp(m *models.Proposal) *ProposalResp {
	if m == nil {
		return nil
	}
	resp := &ProposalResp{
		ID:                m.ID,
		CreatedAt:         m.CreatedAt,
		UpdatedAt:         m.UpdatedAt,
		Network:           m.Network,
		UserID:            m.UserID,
		User:              NewMiniUserResp(m.User),
		Type:              m.Type,
		ChoiceType:        m.ChoiceType,
		Message:           m.Message,
		Signature:         m.Signature,
		Snapshot:          m.Snapshot,
		Name:              m.Name,
		Body:              m.Body,
		Timestamp:         m.Timestamp,
		Start:             m.Start,
		End:               m.End,
		IpfsHash:          fmt.Sprintf("%s/api/ipfs/%s", configs.GetConfig().WebUrl, m.IpfsHash),
		ProposalThreshold: m.ProposalThreshold,
		Status:            m.Status,
		Choices:           NewProposalChoiceRespArr(m.Choices),
	}
	return resp
}

func NewProposalRespArr(arr []*models.Proposal) []*ProposalResp {
	resps := []*ProposalResp{}
	for _, m := range arr {
		resps = append(resps, NewProposalResp(m))
	}
	return resps
}
