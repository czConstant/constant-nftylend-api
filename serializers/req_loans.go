package serializers

import (
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type CreateLoanReq struct {
	Network         models.Network   `json:"chain"`
	Borrower        string           `json:"borrower"`
	CurrencyID      uint             `json:"currency_id"`
	PrincipalAmount numeric.BigFloat `json:"principal_amount"`
	InterestRate    float64          `json:"interest_rate"`
	Duration        uint             `json:"duration"`
	ContractAddress string           `json:"contract_address"`
	TokenID         string           `json:"token_id"`
	Signature       string           `json:"signature"`
	NonceHex        string           `json:"nonce_hex"`
}

type CreateLoanOfferReq struct {
	Lender          string           `json:"lender"`
	PrincipalAmount numeric.BigFloat `json:"principal_amount"`
	InterestRate    float64          `json:"interest_rate"`
	Duration        uint             `json:"duration"`
	Signature       string           `json:"signature"`
	NonceHex        string           `json:"nonce_hex"`
}

type CreateLoanNearReq struct {
	ContractAddress string `json:"contract_address"`
	TokenID         string `json:"token_id"`
}
