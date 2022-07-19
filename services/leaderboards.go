package services

import (
	"context"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
)

func (s *NftLend) GetLeaderBoardAtNow(ctx context.Context, network models.Network, rptDate *time.Time) ([]*models.LeaderBoardData, error) {
	if rptDate == nil {
		t := helpers.GetStartDayOfMonth(time.Now())
		rptDate = &t
	} else {
		t := helpers.GetStartDayOfMonth(*rptDate)
		rptDate = &t
	}
	m, err := s.ld.GetLeaderBoardByMonth(
		daos.GetDBMainCtx(ctx),
		network,
		*rptDate,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if m == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return m, nil
}

func (s *NftLend) GetLeaderBoardDetail(ctx context.Context, network models.Network, rptDate *time.Time) (*models.Leaderboard, error) {
	return nil, nil
}
