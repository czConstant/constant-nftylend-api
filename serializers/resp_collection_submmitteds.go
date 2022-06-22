package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type CollectionSubmittedResp struct {
	ID              uint                             `json:"id"`
	CreatedAt       time.Time                        `json:"created_at"`
	UpdatedAt       time.Time                        `json:"updated_at"`
	Network         models.Network                   `json:"network"`
	Name            string                           `json:"name"`
	Description     string                           `json:"description"`
	Creator         string                           `json:"creator"`
	ContractAddress string                           `json:"contract_address"`
	TokenSeriesID   string                           `json:"token_series_id"`
	ContactInfo     string                           `json:"contact_info"`
	Verified        bool                             `json:"verified"`
	WhoVerified     string                           `json:"who_verified"`
	Status          models.CollectionSubmittedStatus `json:"status"`
}

func NewCollectionSubmittedResp(m *models.CollectionSubmitted) *CollectionSubmittedResp {
	if m == nil {
		return nil
	}
	resp := &CollectionSubmittedResp{
		ID:            m.ID,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
		Network:       m.Network,
		Name:          m.Name,
		Description:   m.Description,
		TokenSeriesID: m.TokenSeriesID,
		Creator:       m.Creator,
		Verified:      m.Verified,
		WhoVerified:   m.WhoVerified,
		Status:        m.Status,
	}
	return resp
}

func NewCollectionSubmittedRespArr(arr []*models.CollectionSubmitted) []*CollectionSubmittedResp {
	resps := []*CollectionSubmittedResp{}
	for _, m := range arr {
		resps = append(resps, NewCollectionSubmittedResp(m))
	}
	return resps
}
