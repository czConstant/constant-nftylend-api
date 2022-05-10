package serializers

import "github.com/czConstant/constant-nftylend-api/models"

type UpdateUserSettingReq struct {
	Address         string         `json:"address"`
	Network         models.Network `json:"network"`
	Email           string         `json:"email"`
	NewsNotiEnabled bool           `json:"news_noti_enabled"`
	LoanNotiEnabled bool           `json:"loan_noti_enabled"`
}
