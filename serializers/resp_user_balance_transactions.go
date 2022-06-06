package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type UserBalanceTransactionResp struct {
	ID                     uint                                `json:"id"`
	CreatedAt              time.Time                           `json:"created_at"`
	UpdatedAt              time.Time                           `json:"updated_at"`
	Network                models.Network                      `json:"network"`
	UserBalanceID          uint                                `json:"user_balance_id"`
	Type                   models.UserBalanceTransactionType   `json:"type"`
	CurrencyID             uint                                `json:"currency_id"`
	Currency               *CurrencyResp                       `json:"currency"`
	ToAddress              string                              `json:"to_address"`
	Amount                 numeric.BigFloat                    `json:"amount"`
	Signature              string                              `json:"signature"`
	TxHash                 string                              `json:"tx_hash"`
	Status                 models.UserBalanceTransactionStatus `json:"status"`
	IncentiveTransactionID uint                                `json:"incentive_transaction_id"`
	IncentiveTransaction   *IncentiveTransactionResp           `json:"incentive_transaction"`
}

func NewUserBalanceTransactionResp(m *models.UserBalanceTransaction) *UserBalanceTransactionResp {
	if m == nil {
		return nil
	}
	resp := &UserBalanceTransactionResp{
		ID:                     m.ID,
		CreatedAt:              m.CreatedAt,
		UpdatedAt:              m.UpdatedAt,
		Network:                m.Network,
		UserBalanceID:          m.UserBalanceID,
		Type:                   m.Type,
		CurrencyID:             m.CurrencyID,
		Currency:               NewCurrencyResp(m.Currency),
		ToAddress:              m.ToAddress,
		Amount:                 m.Amount,
		Signature:              m.Signature,
		TxHash:                 m.TxHash,
		Status:                 m.Status,
		IncentiveTransactionID: m.IncentiveTransactionID,
		IncentiveTransaction:   NewIncentiveTransactionResp(m.IncentiveTransaction),
	}
	return resp
}

func NewUserBalanceTransactionRespArr(arr []*models.UserBalanceTransaction) []*UserBalanceTransactionResp {
	resps := []*UserBalanceTransactionResp{}
	for _, m := range arr {
		resps = append(resps, NewUserBalanceTransactionResp(m))
	}
	return resps
}
