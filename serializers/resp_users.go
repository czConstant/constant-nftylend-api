package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type UserResp struct {
	ID              uint           `json:"id"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	Network         models.Network `json:"network"`
	Address         string         `json:"address"`
	Email           string         `json:"email"`
	UserName        string         `json:"user_name"`
	NewsNotiEnabled bool           `json:"news_noti_enabled"`
	LoanNotiEnabled bool           `json:"loan_noti_enabled"`
	SeenNotiID      uint           `json:"seen_noti_id"`
	IsVerified      bool           `json:"is_verified"`
}

func NewUserResp(m *models.User) *UserResp {
	if m == nil {
		return nil
	}
	resp := &UserResp{
		ID:              m.ID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		Network:         m.Network,
		Address:         m.Address,
		Email:           m.Email,
		UserName:        m.UserName,
		NewsNotiEnabled: m.NewsNotiEnabled,
		LoanNotiEnabled: m.LoanNotiEnabled,
		IsVerified:      m.IsVerified,
	}
	return resp
}

func NewMiniUserResp(m *models.User) *UserResp {
	if m == nil {
		return nil
	}
	resp := &UserResp{
		ID:         m.ID,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
		Network:    m.Network,
		Address:    m.Address,
		IsVerified: m.IsVerified,
	}
	return resp
}

func NewUserRespArr(arr []*models.User) []*UserResp {
	resps := []*UserResp{}
	for _, m := range arr {
		resps = append(resps, NewUserResp(m))
	}
	return resps
}
