package serializers

import (
	"encoding/json"
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type AssetResp struct {
	ID                    uint            `json:"id"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	CollectionID          uint            `json:"collection_id"`
	Collection            *CollectionResp `json:"collection"`
	SeoURL                string          `json:"seo_url"`
	ContractAddress       string          `json:"contract_address"`
	TokenURL              string          `json:"token_url"`
	TokenID               string          `json:"token_id"`
	Name                  string          `json:"name"`
	SellerFeeRate         float64         `json:"seller_fee_rate"`
	Attributes            interface{}     `json:"attributes"`
	OriginNetwork         models.Network  `json:"origin_network"`
	OriginContractAddress string          `json:"origin_contract_address"`
	OriginTokenID         string          `json:"origin_token_id"`
	NewLoan               *LoanResp       `json:"new_loan"`
}

func NewAssetResp(m *models.Asset) *AssetResp {
	if m == nil {
		return nil
	}
	attr := []interface{}{}
	json.Unmarshal([]byte(m.Attributes), &attr)
	resp := &AssetResp{
		ID:                    m.ID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		CollectionID:          m.CollectionID,
		Collection:            NewCollectionResp(m.Collection),
		SeoURL:                m.SeoURL,
		ContractAddress:       m.ContractAddress,
		TokenURL:              m.TokenURL,
		TokenID:               m.TokenID,
		Name:                  m.Name,
		SellerFeeRate:         m.SellerFeeRate,
		Attributes:            attr,
		OriginNetwork:         m.OriginNetwork,
		OriginContractAddress: m.OriginContractAddress,
		OriginTokenID:         m.OriginTokenID,
		NewLoan:               NewLoanResp(m.NewLoan),
	}
	return resp
}

func NewAssetRespArr(arr []*models.Asset) []*AssetResp {
	resps := []*AssetResp{}
	for _, m := range arr {
		resps = append(resps, NewAssetResp(m))
	}
	return resps
}
