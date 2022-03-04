package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type AssetTransactionResp struct {
	ID            uint                        `json:"id"`
	CreatedAt     time.Time                   `json:"created_at"`
	UpdatedAt     time.Time                   `json:"updated_at"`
	Source        string                      `json:"source"`
	Network       models.Chain                `json:"network"`
	AssetID       uint                        `json:"asset_id"`
	Asset         *AssetResp                  `json:"asset"`
	Type          models.AssetTransactionType `json:"type"`
	Seller        string                      `json:"seller"`
	Buyer         string                      `json:"buyer"`
	TransactionAt *time.Time                  `json:"transaction_at"`
	Amount        numeric.BigFloat            `json:"amount"`
	CurrencyID    uint                        `json:"currency_id"`
	Currency      *CurrencyResp               `json:"currency"`
	TransactionID string                      `json:"transaction_id"`
}

func NewAssetTransactionResp(m *models.AssetTransaction) *AssetTransactionResp {
	if m == nil {
		return nil
	}
	resp := &AssetTransactionResp{
		ID:            m.ID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		Source:        m.Source,
		Network:       m.Network,
		AssetID:       m.AssetID,
		Asset:         NewAssetResp(m.Asset),
		Type:          m.Type,
		Seller:        m.Seller,
		Buyer:         m.Buyer,
		TransactionAt: m.TransactionAt,
		Amount:        m.Amount,
		CurrencyID:    m.CurrencyID,
		Currency:      NewCurrencyResp(m.Currency),
		TransactionID: m.TransactionID,
	}
	return resp
}

func NewAssetTransactionRespArr(arr []*models.AssetTransaction) []*AssetTransactionResp {
	resps := []*AssetTransactionResp{}
	for _, m := range arr {
		resps = append(resps, NewAssetTransactionResp(m))
	}
	return resps
}
