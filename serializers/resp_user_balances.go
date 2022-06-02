package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type UserBalanceResp struct {
	ID        uint           `json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	UserID    uint           `json:"user_id"`
	Network   models.Network `json:"network"`
	Address   string         `json:"address"`
}

func NewUserBalanceResp(m *models.UserBalance) *UserBalanceResp {
	if m == nil {
		return nil
	}
	resp := &UserBalanceResp{
		ID:        m.ID,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
		UserID:    m.UserID,
		Network:   m.Network,
		Address:   m.Address,
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
