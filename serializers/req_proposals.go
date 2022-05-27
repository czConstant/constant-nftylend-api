package serializers

import "github.com/czConstant/constant-nftylend-api/models"

type CreateProposalReq struct {
	Network models.Network `json:"network"`
	Address string         `json:"address"`
	Msg     string         `json:"msg"`
	Sig     string         `json:"sig"`
}
