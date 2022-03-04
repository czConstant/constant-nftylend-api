package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type CurrencyResp struct {
	ID              uint         `json:"id"`
	CreatedAt       time.Time    `json:"created_at"`
	UpdatedAt       time.Time    `json:"updated_at"`
	Network         models.Chain `json:"network"`
	ContractAddress string       `json:"contract_address"`
	Decimals        uint         `json:"decimals"`
	Symbol          string       `json:"symbol"`
	Name            string       `json:"name"`
	IconURL         string       `json:"icon_url"`
	AdminFeeAddress string       `json:"admin_fee_address"`
}

func NewCurrencyResp(m *models.Currency) *CurrencyResp {
	if m == nil {
		return nil
	}
	resp := &CurrencyResp{
		ID:              m.ID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		Network:         m.Network,
		ContractAddress: m.ContractAddress,
		Decimals:        m.Decimals,
		Symbol:          m.Symbol,
		Name:            m.Name,
		IconURL:         m.IconURL,
		AdminFeeAddress: m.AdminFeeAddress,
	}
	return resp
}

func NewCurrencyRespArr(arr []*models.Currency) []*CurrencyResp {
	resps := []*CurrencyResp{}
	for _, m := range arr {
		resps = append(resps, NewCurrencyResp(m))
	}
	return resps
}
