package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
)

type CollectionResp struct {
	ID                    uint             `json:"id"`
	CreatedAt             time.Time        `json:"created_at"`
	UpdatedAt             time.Time        `json:"updated_at"`
	Network               models.Network   `json:"network"`
	SeoURL                string           `json:"seo_url"`
	Name                  string           `json:"name"`
	Description           string           `json:"description"`
	Verified              bool             `json:"verified"`
	ListingAsset          *AssetResp       `json:"listing_asset"`
	ListingTotal          uint             `json:"listing_total"`
	TotalVolume           numeric.BigFloat `json:"total_volume"`
	TotalListed           uint             `json:"total_listed"`
	Avg24hAmount          numeric.BigFloat `json:"avg24h_amount"`
	OriginNetwork         models.Network   `json:"origin_network"`
	OriginContractAddress string           `json:"origin_contract_address"`
	CreatorURL            string           `json:"creator_url"`
	TwitterID             string           `json:"twitter_id"`
	DiscordURL            string           `json:"discord_url"`
	CoverURL              string           `json:"cover_url"`
	RandAsset             *AssetResp       `json:"rand_asset"`
}

func NewCollectionResp(m *models.Collection) *CollectionResp {
	if m == nil {
		return nil
	}
	resp := &CollectionResp{
		ID:                    m.ID,
		CreatedAt:             m.CreatedAt,
		UpdatedAt:             m.UpdatedAt,
		Network:               m.Network,
		SeoURL:                m.SeoURL,
		Name:                  m.Name,
		Verified:              m.Verified,
		Description:           m.Description,
		OriginNetwork:         m.OriginNetwork,
		OriginContractAddress: m.OriginContractAddress,
		CreatorURL:            m.CreatorURL,
		TwitterID:             m.TwitterID,
		DiscordURL:            m.DiscordURL,
		CoverURL:              m.CoverURL,
		ListingAsset:          NewAssetResp(m.ListingAsset),
		RandAsset:             NewAssetResp(m.RandAsset),
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
