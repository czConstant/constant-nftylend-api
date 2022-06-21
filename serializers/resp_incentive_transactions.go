package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type IncentiveTransactionResp struct {
	ID                 uint                              `json:"id"`
	CreatedAt          time.Time                         `json:"created_at"`
	UpdatedAt          time.Time                         `json:"updated_at"`
	Network            models.Network                    `json:"network"`
	IncentiveProgramID uint                              `json:"incentive_program_id"`
	Type               models.IncentiveTransactionType   `json:"type"`
	UserID             uint                              `json:"user_id"`
	User               *UserResp                         `json:"user"`
	CurrencyID         uint                              `json:"currency_id"`
	Currency           *CurrencyResp                     `json:"currency"`
	LoanID             uint                              `json:"loan_id"`
	Loan               *LoanResp                         `json:"loan"`
	Amount             numeric.BigFloat                  `json:"amount"`
	LockUntilAt        *time.Time                        `json:"lock_until_at"`
	UnlockedAt         *time.Time                        `json:"unlocked_at"`
	Status             models.IncentiveTransactionStatus `json:"status"`
}

func NewIncentiveTransactionResp(m *models.IncentiveTransaction) *IncentiveTransactionResp {
	if m == nil {
		return nil
	}
	resp := &IncentiveTransactionResp{
		ID:                 m.ID,
		CreatedAt:          m.CreatedAt,
		UpdatedAt:          m.UpdatedAt,
		Network:            m.Network,
		IncentiveProgramID: m.IncentiveProgramID,
		Type:               m.Type,
		UserID:             m.UserID,
		User:               NewMiniUserResp(m.User),
		CurrencyID:         m.CurrencyID,
		Currency:           NewCurrencyResp(m.Currency),
		LoanID:             m.LoanID,
		Loan:               NewLoanResp(m.Loan),
		Amount:             m.Amount,
		LockUntilAt:        m.LockUntilAt,
		UnlockedAt:         m.UnlockedAt,
		Status:             m.Status,
	}
	return resp
}

func NewIncentiveTransactionRespArr(arr []*models.IncentiveTransaction) []*IncentiveTransactionResp {
	resps := []*IncentiveTransactionResp{}
	for _, m := range arr {
		resps = append(resps, NewIncentiveTransactionResp(m))
	}
	return resps
}
