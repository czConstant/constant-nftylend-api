package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type UserBalanceResp struct {
	ID             uint             `json:"id"`
	CreatedAt      time.Time        `json:"created_at"`
	UpdatedAt      time.Time        `json:"updated_at"`
	UserID         uint             `json:"user_id"`
	User           *UserResp        `json:"user"`
	Network        models.Network   `json:"network"`
	CurrencyID     uint             `json:"currency_id"`
	Currency       *CurrencyResp    `json:"currency"`
	Balance        numeric.BigFloat `json:"balance"`
	LockedBalance  numeric.BigFloat `json:"locked_balance"`
	ClaimedBalance numeric.BigFloat `json:"claimed_balance"`
}

func NewUserBalanceResp(m *models.UserBalance) *UserBalanceResp {
	if m == nil {
		return nil
	}
	resp := &UserBalanceResp{
		ID:             m.ID,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
		Network:        m.Network,
		UserID:         m.UserID,
		User:           NewMiniUserResp(m.User),
		CurrencyID:     m.CurrencyID,
		Currency:       NewCurrencyResp(m.Currency),
		Balance:        m.Balance,
		LockedBalance:  m.LockedBalance,
		ClaimedBalance: m.ClaimedBalance,
	}
	return resp
}

func NewUserBalanceRespArr(arr []*models.UserBalance) []*UserBalanceResp {
	resps := []*UserBalanceResp{}
	for _, m := range arr {
		resps = append(resps, NewUserBalanceResp(m))
	}
	return resps
}
