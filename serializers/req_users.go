package serializers

import (
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type UpdateUserSettingReq struct {
	Network         models.Network `json:"network"`
	Address         string         `json:"address"`
	Email           string         `json:"email"`
	NewsNotiEnabled *bool          `json:"news_noti_enabled"`
	LoanNotiEnabled *bool          `json:"loan_noti_enabled"`
}

type UserConnectedReq struct {
	SignatureTimestampReq
	ReferrerCode string `json:"referrer_code"`
}

type ClaimUserBalanceReq struct {
	UserID     uint             `json:"user_id"`
	CurrencyID uint             `json:"currency_id"`
	ToAddress  string           `json:"to_address"`
	Amount     numeric.BigFloat `json:"amount"`
	Timestamp  int64            `json:"timestamp"`
	Signature  string           `json:"signature"`
}

type UserVerifyEmailReq struct {
	SignatureTimestampReq
	Email string `json:"email"`
}

type UserVerifyEmailTokenReq struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
