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
	if rptDate == nil {
		t := helpers.GetStartDayOfMonth(time.Now())
		rptDate = &t
	} else {
		t := helpers.GetStartDayOfMonth(*rptDate)
		rptDate = &t
	}
	m, err := s.lbd.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"rpt_date = ?": []interface{}{rptDate},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetPrevLeaderBoardDetail(ctx context.Context, network models.Network, rptDate *time.Time) (*models.Leaderboard, error) {
	if rptDate == nil {
		t := helpers.GetStartDayOfMonth(time.Now())
		rptDate = &t
	} else {
		t := helpers.GetStartDayOfMonth(*rptDate)
		rptDate = &t
	}
	t := rptDate.AddDate(0, -1, 0)
	rptDate = &t
	m, err := s.lbd.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"rpt_date = ?": []interface{}{rptDate},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}
