package serializers

import (
	"time"

	"github.com/czConstant/constant-nftylend-api/models"
)

type LeaderboardResp struct {
	Network  models.Network `json:"network"`
	RptDate  *time.Time     `json:"rpt_date"`
	Rewards  string         `json:"rewards"`
	ImageURL string         `json:"image_url"`
}

func NewLeaderboardResp(m *models.Leaderboard) *LeaderboardResp {
	if m == nil {
		return nil
	}
	resp := &LeaderboardResp{
		Network:  m.Network,
		RptDate:  m.RptDate,
		Rewards:  m.Rewards,
		ImageURL: m.ImageURL,
	}
	return resp
}

func NewLeaderboardRespArr(arr []*models.Leaderboard) []*LeaderboardResp {
	resps := []*LeaderboardResp{}
	for _, m := range arr {
		resps = append(resps, NewLeaderboardResp(m))
	}
	return resps
}
