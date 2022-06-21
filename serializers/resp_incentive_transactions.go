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

type AffiliateStatsResp struct {
	CommisionsRate    float64          `json:"commisions_rate"`
	TotalCommisions   numeric.BigFloat `json:"total_commisions"`
	TotalUsers        uint             `json:"total_users"`
	TotalTransactions uint             `json:"total_transactions"`
}

func NewAffiliateStatsRespResp(m *models.AffiliateStats, commisionsRate float64) *AffiliateStatsResp {
	if m == nil {
		return nil
	}
	resp := &AffiliateStatsResp{
		CommisionsRate:    commisionsRate,
		TotalCommisions:   m.TotalCommisions,
		TotalUsers:        m.TotalUsers,
		TotalTransactions: m.TotalTransactions,
	}
	return resp
}

type AffiliateVolumesResp struct {
	RptDate         *time.Time       `json:"rpt_date"`
	TotalCommisions numeric.BigFloat `json:"total_commisions"`
}

func NewAffiliateVolumesResp(m *models.AffiliateVolumes) *AffiliateVolumesResp {
	if m == nil {
		return nil
	}
	resp := &AffiliateVolumesResp{
		RptDate:         m.RptDate,
		TotalCommisions: m.TotalCommisions,
	}
	return resp
}

func NewAffiliateVolumesRespArr(arr []*models.AffiliateVolumes) []*AffiliateVolumesResp {
	resps := []*AffiliateVolumesResp{}
	for _, m := range arr {
		resps = append(resps, NewAffiliateVolumesResp(m))
	}
	return resps
}
