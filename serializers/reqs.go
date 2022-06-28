package serializers

import "github.com/czConstant/constant-nftylend-api/models"

type SignatureReq struct {
	Network   models.Network `json:"network"`
	Address   string         `json:"address"`
	Timestamp int64          `json:"timestamp"`
	Message   string         `json:"message"`
	Signature string         `json:"signature"`
}

type SignatureTimestampReq struct {
	Network   models.Network `json:"network"`
	Address   string         `json:"address"`
	Timestamp int64          `json:"timestamp"`
	Signature string         `json:"signature"`
}

type BaseReq struct {
	Network models.Network `json:"network"`
	Address string         `json:"address"`
}
