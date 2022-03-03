package serializers

import (
	"time"

	"github.com/czConstant/constant-nftlend-api/models"
	"github.com/czConstant/constant-nftlend-api/types/numeric"
)

type LoanTransactionResp struct {
	ID              uint                       `json:"id"`
	CreatedAt       time.Time                  `json:"created_at"`
	UpdatedAt       time.Time                  `json:"updated_at"`
	Network         models.Chain               `json:"network"`
	LoanID          uint                       `json:"loan_id"`
	Loan            *LoanResp                  `json:"loan"`
	Type            models.LoanTransactionType `json:"type"`
	Borrower        string                     `json:"borrower"`
	Lender          string                     `json:"lender"`
	StartedAt       *time.Time                 `json:"started_at"`
	Duration        uint                       `json:"duration"`
	ExpiredAt       *time.Time                 `json:"expired_at"`
	PrincipalAmount numeric.BigFloat           `json:"principal_amount"`
	InterestRate    float64                    `json:"interest_rate"`
	TxHash          string                     `json:"tx_hash"`
}

func NewLoanTransactionResp(m *models.LoanTransaction) *LoanTransactionResp {
	if m == nil {
		return nil
	}
	resp := &LoanTransactionResp{
		ID:              m.ID,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
		Network:         m.Network,
		Type:            m.Type,
		LoanID:          m.LoanID,
		Loan:            NewLoanResp(m.Loan),
		Borrower:        m.Borrower,
		Lender:          m.Lender,
		PrincipalAmount: m.PrincipalAmount,
		InterestRate:    m.InterestRate,
		StartedAt:       m.StartedAt,
		Duration:        m.Duration,
		ExpiredAt:       m.ExpiredAt,
		TxHash:          m.TxHash,
	}
	return resp
}

func NewLoanTransactionRespArr(arr []*models.LoanTransaction) []*LoanTransactionResp {
	resps := []*LoanTransactionResp{}
	for _, m := range arr {
		resps = append(resps, NewLoanTransactionResp(m))
	}
	return resps
}
