package services

import (
	"context"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
)

func (s *NftLend) GetAffiliateStats(ctx context.Context, network models.Network, address string) (*models.AffiliateStats, error) {
	user, err := s.GetUser(
		ctx,
		network,
		address,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	nToken, err := s.getLendCurrencyBySymbol(
		daos.GetDBMainCtx(ctx),
		network,
		models.SymbolNEARToken,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	m, err := s.itd.GetAffiliateStats(
		daos.GetDBMainCtx(ctx),
		user.ID,
		nToken.ID,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetAffiliateShareRate(ctx context.Context, network models.Network, address string) (float64, error) {
	user, err := s.GetUser(
		ctx,
		network,
		address,
	)
	if err != nil {
		return 0, errs.NewError(err)
	}
	var shareRate float64
	m, err := s.ipdd.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			`exists(
				select 1
				from incentive_programs
				where 1 = 1
				  and incentive_programs.network = ?
				  and incentive_program_details.incentive_program_id = incentive_programs.id
				  and (? between incentive_programs.start and incentive_programs.end)
				  and incentive_programs.status = ?
			)`: []interface{}{
				network,
				time.Now(),
				models.IncentiveProgramStatusActived,
			},
			"user_rank = ? or user_rank = ?": []interface{}{
				models.UserRankAll,
				user.Rank,
			},
			"type in (?)": []interface{}{
				[]models.IncentiveTransactionType{
					models.IncentiveTransactionTypeAffiliateBorrowerLoanDone,
				},
			},
		},
		map[string][]interface{}{},
		[]string{
			"amount desc",
		},
	)
	if err != nil {
		return 0, errs.NewError(err)
	}
	if m != nil {
		if m.RewardType == models.IncentiveTransactionRewardTypeRateOfLoan {
			shareRate, _ = m.Amount.Float.Float64()
		}
	}
	return shareRate, nil
}

func (s *NftLend) GetAffiliateVolumes(ctx context.Context, network models.Network, address string, rptBy string, limit uint) ([]*models.AffiliateVolumes, error) {
	user, err := s.GetUser(
		ctx,
		network,
		address,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	nToken, err := s.getLendCurrencyBySymbol(
		daos.GetDBMainCtx(ctx),
		network,
		models.SymbolNEARToken,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	m, err := s.itd.GetAffiliateVolumes(
		daos.GetDBMainCtx(ctx),
		user.ID,
		nToken.ID,
		rptBy,
		limit,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}
