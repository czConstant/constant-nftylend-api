package serializers

type UpdateUserSettingReq struct {
	SignatureTimestampReq
	Username        string `json:"username"`
	NewsNotiEnabled *bool  `json:"news_noti_enabled"`
	LoanNotiEnabled *bool  `json:"loan_noti_enabled"`
}

type UserConnectedReq struct {
	SignatureTimestampReq
	ReferrerCode string `json:"referrer_code"`
}

type ClaimUserBalanceReq struct {
	SignatureTimestampReq
	CurrencyID uint `json:"currency_id"`
}

type UserVerifyEmailReq struct {
	SignatureTimestampReq
	Email string `json:"email"`
}

type UserVerifyEmailTokenReq struct {
	Email string `json:"email"`
	Token string `json:"token"`
}
