package serializers

import (
	"time"

	"github.com/czConstant/constant-nftlend-api/models"
	"github.com/czConstant/constant-nftlend-api/types/numeric"
)

type LoanOfferResp struct {
	ID                  uint                   `json:"id"`
	CreatedAt           time.Time              `json:"created_at"`
	UpdatedAt           time.Time              `json:"updated_at"`
	LoanID              uint                   `json:"loan_id"`
	Loan                *LoanResp              `json:"loan"`
	Lender              string                 `json:"lender"`
	PrincipalAmount     numeric.BigFloat       `json:"principal_amount"`
	InterestRate        float64                `json:"interest_rate"`
	Duration            uint                   `json:"duration"`
	NonceHex            string                 `json:"nonce_hex"`
	Signature           string                 `json:"signature"`
	Status              models.LoanOfferStatus `json:"status"`
	DataOfferAddress    string                 `json:"data_offer_address"`
	DataCurrencyAddress string                 `json:"data_currency_address"`
	MakeTxHash          string                 `json:"make_tx_hash"`
	AcceptTxHash        string                 `json:"accept_tx_hash"`
	CancelTxHash        string                 `json:"cancel_tx_hash"`
	CloseTxHash         string                 `json:"close_tx_hash"`
}

func NewLoanOfferResp(m *models.LoanOffer) *LoanOfferResp {
	if m == nil {
		return nil
	}
	resp := &LoanOfferResp{
		ID:                  m.ID,
		CreatedAt:           m.CreatedAt,
		UpdatedAt:           m.UpdatedAt,
		Lender:              m.Lender,
		PrincipalAmount:     m.PrincipalAmount,
		InterestRate:        m.InterestRate,
		Duration:            m.Duration,
		NonceHex:            m.NonceHex,
		Signature:           m.Signature,
		Status:              m.Status,
		DataOfferAddress:    m.DataOfferAddress,
		DataCurrencyAddress: m.DataCurrencyAddress,
		LoanID:              m.LoanID,
		Loan:                NewLoanResp(m.Loan),
		MakeTxHash:          m.MakeTxHash,
		AcceptTxHash:        m.AcceptTxHash,
		CancelTxHash:        m.CancelTxHash,
		CloseTxHash:         m.CloseTxHash,
	}
	return resp
}

func NewLoanOfferRespArr(arr []*models.LoanOffer) []*LoanOfferResp {
	resps := []*LoanOfferResp{}
	for _, m := range arr {
		resps = append(resps, NewLoanOfferResp(m))
	}
	return resps
}
