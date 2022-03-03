package serializers

import (
	"time"

	"github.com/czConstant/constant-nftlend-api/models"
	"github.com/czConstant/constant-nftlend-api/types/numeric"
)

type CollectionResp struct {
	ID           uint             `json:"id"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	SeoURL       string           `json:"seo_url"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	ListingAsset *AssetResp       `json:"listing_asset"`
	ListingTotal uint             `json:"listing_total"`
	TotalVolume  numeric.BigFloat `json:"total_volume"`
	TotalListed  uint             `json:"total_listed"`
	Avg24hAmount numeric.BigFloat `json:"avg24h_amount"`
}

func NewCollectionResp(m *models.Collection) *CollectionResp {
	if m == nil {
		return nil
	}
	resp := &CollectionResp{
		ID:           m.ID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		SeoURL:       m.SeoURL,
		Name:         m.Name,
		Description:  m.Description,
		ListingAsset: NewAssetResp(m.ListingAsset),
	}
	return resp
}

func NewCollectionRespArr(arr []*models.Collection) []*CollectionResp {
	resps := []*CollectionResp{}
	for _, m := range arr {
		resps = append(resps, NewCollectionResp(m))
	}
	return resps
}
