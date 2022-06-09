package serializers

import "github.com/czConstant/constant-nftylend-api/models"

type CreateProposalReq struct {
	Network   models.Network `json:"network"`
	Address   string         `json:"address"`
	Message   string         `json:"message"`
	Signature string         `json:"signature"`
}

type CreateProposalVoteReq struct {
	Network   models.Network `json:"network"`
	Address   string         `json:"address"`
	Message   string         `json:"message"`
	Signature string         `json:"signature"`
}
