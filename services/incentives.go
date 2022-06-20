package services

import (
	"context"
	"fmt"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) IncentiveForLoan(tx *gorm.DB, incentiveTransactionType models.IncentiveTransactionType, loanID uint) error {
	loan, err := s.ld.FirstByID(
		tx,
		loanID,
		map[string][]interface{}{},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	if loan == nil {
		return errs.NewError(errs.ErrBadRequest)
	}
	var checkIncentiveTime *time.Time
	var address string
	switch incentiveTransactionType {
	case models.IncentiveTransactionTypeBorrowerLoanListed:
		{
			checkIncentiveTime = &loan.CreatedAt
			address = loan.Owner
		}
	case models.IncentiveTransactionTypeBorrowerLoanDelisted:
		{
			address = loan.Owner
			checkIncentiveTime = &loan.CreatedAt
		}
	case models.IncentiveTransactionTypeLenderLoanMatched:
		{
			address = loan.Lender
			checkIncentiveTime = loan.OfferStartedAt
		}
	case models.IncentiveTransactionTypeAffiliateLoanDone:
		{
			borrower, err := s.ud.FirstByID(
				tx,
				loan.BorrowerUserID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if borrower != nil {
				if borrower.ReferrerUserID > 0 {
					referrer, err := s.ud.FirstByID(
						tx,
						borrower.ReferrerUserID,
						map[string][]interface{}{},
						false,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if referrer != nil {
						return errs.NewError(errs.ErrBadRequest)
					}
					address = referrer.Address
				}
			}
			checkIncentiveTime = loan.OfferStartedAt
		}
	default:
		{
			return errs.NewError(errs.ErrBadRequest)
		}
	}
	if address != "" {
		ipdMs, err := s.ipdd.Find(
			tx,
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
					loan.Network,
					checkIncentiveTime,
					models.IncentiveProgramStatusActived,
				},
				"type = ?": []interface{}{incentiveTransactionType},
			},
			map[string][]interface{}{
				"IncentiveProgram": []interface{}{},
			},
			[]string{},
			0,
			999999,
		)
		if err != nil {
			return errs.NewError(err)
		}
		if len(ipdMs) > 0 {
			user, err := s.getUser(
				tx,
				loan.Network,
				address,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for _, ipdM := range ipdMs {
				ipM := ipdM.IncentiveProgram
				if uint(loan.ValidAt.Sub(*loan.StartedAt).Seconds()) >= ipM.LoanValidDuration {
					itM, err := s.itd.First(
						tx,
						map[string][]interface{}{
							"incentive_program_id = ?": []interface{}{ipM.ID},
							"type = ?":                 []interface{}{ipdM.Type},
							"user_id = ?":              []interface{}{user.ID},
							"loan_id = ?":              []interface{}{loan.ID},
						},
						map[string][]interface{}{},
						[]string{},
					)
					if err != nil {
						return errs.NewError(err)
					}
					if itM == nil {
						isOk := true
						txStatus := models.IncentiveTransactionStatusLocked
						switch incentiveTransactionType {
						case models.IncentiveTransactionTypeBorrowerLoanDelisted:
							{
								// check tx for listed
								itM, err = s.itd.First(
									tx,
									map[string][]interface{}{
										"user_id = ?":              []interface{}{user.ID},
										"incentive_program_id = ?": []interface{}{ipdM.IncentiveProgramID},
										"type = ?":                 []interface{}{models.IncentiveTransactionTypeBorrowerLoanListed},
										"loan_id = ?":              []interface{}{loan.ID},
										"status = ?":               []interface{}{models.IncentiveTransactionStatusLocked},
									},
									map[string][]interface{}{},
									[]string{},
								)
								if err != nil {
									return errs.NewError(err)
								}
								if itM == nil {
									isOk = false
								} else {
									itM.Status = models.IncentiveTransactionStatusRevoked
									err = s.itd.Save(
										tx,
										itM,
									)
									if err != nil {
										return errs.NewError(err)
									}
								}
								txStatus = models.IncentiveTransactionStatusDone
							}
						}
						if isOk {
							var amount numeric.BigFloat
							var currencyID uint
							switch ipdM.RewardType {
							case models.IncentiveTransactionRewardTypeAmount:
								{
									currencyID = ipM.CurrencyID
									amount = ipdM.Amount
								}
							case models.IncentiveTransactionRewardTypeRateOfLoan:
								{
									currencyID = loan.CurrencyID
									amount = numeric.BigFloat{*models.MulBigFloats(&loan.OfferPrincipalAmount.Float, &ipdM.Amount.Float)}
								}
							default:
								{
									return errs.NewError(errs.ErrBadRequest)
								}
							}
							itM = &models.IncentiveTransaction{
								Network:            ipM.Network,
								IncentiveProgramID: ipM.ID,
								Type:               ipdM.Type,
								UserID:             user.ID,
								CurrencyID:         currencyID,
								LoanID:             loanID,
								Amount:             amount,
								LockUntilAt:        helpers.TimeAdd(*checkIncentiveTime, time.Duration(ipM.LockDuration)*time.Second),
								UnlockedAt:         nil,
								Status:             txStatus,
							}
							err = s.itd.Create(
								tx,
								itM,
							)
							if err != nil {
								return errs.NewError(err)
							}
							reference := fmt.Sprintf("it_%d_locked", itM.ID)
							switch incentiveTransactionType {
							case models.IncentiveTransactionTypeBorrowerLoanDelisted:
								{
									reference = fmt.Sprintf("it_%d_revoked", itM.ID)
								}
							}
							userBalance, err := s.getUserBalance(
								tx,
								itM.UserID,
								itM.CurrencyID,
								true,
							)
							if err != nil {
								return errs.NewError(err)
							}
							userBalanceTransaction := &models.UserBalanceTransaction{
								Network:                userBalance.Network,
								UserID:                 userBalance.UserID,
								UserBalanceID:          userBalance.ID,
								CurrencyID:             userBalance.CurrencyID,
								Type:                   models.UserBalanceTransactionTypeIncentive,
								Amount:                 itM.Amount,
								Status:                 models.UserBalanceTransactionStatusDone,
								IncentiveTransactionID: itM.ID,
							}
							err = s.ubtd.Create(
								tx,
								userBalanceTransaction,
							)
							if err != nil {
								return errs.NewError(err)
							}
							err = s.transactionUserBalance(
								tx,
								ipM.Network,
								itM.UserID,
								itM.CurrencyID,
								itM.Amount,
								true,
								false,
								reference,
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					}
				}
			}
		}
	}
	return nil
}

func (s *NftLend) JobIncentiveForUnlock(ctx context.Context) error {
	var retErr error
	itMs, err := s.itd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"type in (?)": []interface{}{[]models.IncentiveTransactionType{
				models.IncentiveTransactionTypeUserAirdropReward,
				models.IncentiveTransactionTypeUserAmaReward,
			}},
			"status = ?": []interface{}{models.IncentiveTransactionStatusPending},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		9999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, itM := range itMs {
		err = s.IncentiveForLock(ctx, itM.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	itMs, err = s.itd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"type in (?)": []interface{}{[]models.IncentiveTransactionType{
				models.IncentiveTransactionTypeBorrowerLoanListed,
				models.IncentiveTransactionTypeLenderLoanMatched,
				models.IncentiveTransactionTypeUserAirdropReward,
				models.IncentiveTransactionTypeUserAmaReward,
			}},
			"status = ?":                []interface{}{models.IncentiveTransactionStatusLocked},
			"lock_until_at is not null": []interface{}{},
			"lock_until_at <= ?":        []interface{}{time.Now()},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		9999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, itM := range itMs {
		err = s.IncentiveForUnlock(ctx, itM.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	return retErr
}

func (s *NftLend) IncentiveForUnlock(ctx context.Context, transactionID uint) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(context.Background()),
		func(tx *gorm.DB) error {
			itM, err := s.itd.FirstByID(
				tx,
				transactionID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if itM == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			switch itM.Type {
			case models.IncentiveTransactionTypeBorrowerLoanListed,
				models.IncentiveTransactionTypeLenderLoanMatched:
				{
				}
			default:
				{
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			if itM.Status != models.IncentiveTransactionStatusLocked {
				return errs.NewError(errs.ErrBadRequest)
			}
			if itM.LockUntilAt.After(time.Now()) {
				return errs.NewError(errs.ErrBadRequest)
			}
			itM.UnlockedAt = helpers.TimeNow()
			itM.Status = models.IncentiveTransactionStatusUnlocked
			err = s.itd.Save(
				tx,
				itM,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.unlockUserBalance(
				tx,
				itM.UserID,
				itM.CurrencyID,
				itM.Amount,
				fmt.Sprintf("it_%d_unlocked", itM.ID),
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) IncentiveForLock(ctx context.Context, transactionID uint) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(context.Background()),
		func(tx *gorm.DB) error {
			itM, err := s.itd.FirstByID(
				tx,
				transactionID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if itM == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			switch itM.Type {
			case models.IncentiveTransactionTypeUserAirdropReward,
				models.IncentiveTransactionTypeUserAmaReward:
				{
				}
			default:
				{
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			if itM.Status != models.IncentiveTransactionStatusPending {
				return errs.NewError(errs.ErrBadRequest)
			}
			user, err := s.getUser(
				tx,
				itM.Network,
				itM.Address,
			)
			if err != nil {
				return errs.NewError(err)
			}
			itM.UserID = user.ID
			itM.Status = models.IncentiveTransactionStatusLocked
			err = s.itd.Save(
				tx,
				itM,
			)
			if err != nil {
				return errs.NewError(err)
			}
			reference := fmt.Sprintf("it_%d_locked", itM.ID)
			userBalance, err := s.getUserBalance(
				tx,
				itM.UserID,
				itM.CurrencyID,
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			userBalanceTransaction := &models.UserBalanceTransaction{
				Network:                userBalance.Network,
				UserID:                 userBalance.UserID,
				UserBalanceID:          userBalance.ID,
				CurrencyID:             userBalance.CurrencyID,
				Type:                   models.UserBalanceTransactionTypeIncentive,
				Amount:                 itM.Amount,
				Status:                 models.UserBalanceTransactionStatusDone,
				IncentiveTransactionID: itM.ID,
			}
			err = s.ubtd.Create(
				tx,
				userBalanceTransaction,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.transactionUserBalance(
				tx,
				itM.Network,
				itM.UserID,
				itM.CurrencyID,
				itM.Amount,
				true,
				false,
				reference,
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}
