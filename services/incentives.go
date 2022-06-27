package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) GetIncentiveTransactions(ctx context.Context, network models.Network, address string, types []string, statuses []string, page int, limit int) ([]*models.IncentiveTransaction, uint, error) {
	user, err := s.GetUser(
		ctx,
		network,
		address,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	filters := map[string][]interface{}{
		"user_id = ?": []interface{}{user.ID},
	}
	if len(types) > 0 {
		filters["type in (?)"] = []interface{}{types}
	}
	if len(statuses) > 0 {
		filters["status in (?)"] = []interface{}{statuses}
	}
	txns, count, err := s.itd.Find4Page(
		daos.GetDBMainCtx(ctx),
		filters,
		map[string][]interface{}{
			"User":     []interface{}{},
			"Currency": []interface{}{},
			"Loan":     []interface{}{},
		},
		[]string{"id desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return txns, count, nil
}

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
	var userID uint
	var refUserID uint
	switch incentiveTransactionType {
	case models.IncentiveTransactionTypeBorrowerLoanListed:
		{
			checkIncentiveTime = &loan.CreatedAt
			userID = loan.BorrowerUserID
		}
	case models.IncentiveTransactionTypeBorrowerLoanDelisted:
		{
			userID = loan.BorrowerUserID
			checkIncentiveTime = &loan.CreatedAt
		}
	case models.IncentiveTransactionTypeLenderLoanMatched:
		{
			userID = loan.LenderUserID
			checkIncentiveTime = loan.OfferStartedAt
		}
	case models.IncentiveTransactionTypeAffiliateBorrowerLoanDone:
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
						if referrer.Type == models.UserTypeAffiliate {
							userID = referrer.ID
						}
					}
				}
			}
			refUserID = borrower.ID
			checkIncentiveTime = loan.OfferStartedAt
		}
	case models.IncentiveTransactionTypeAffiliateLenderLoanDone:
		{
			lender, err := s.ud.FirstByID(
				tx,
				loan.LenderUserID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if lender != nil {
				if lender.ReferrerUserID > 0 {
					referrer, err := s.ud.FirstByID(
						tx,
						lender.ReferrerUserID,
						map[string][]interface{}{},
						false,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if referrer != nil {
						if referrer.Type == models.UserTypeAffiliate {
							userID = referrer.ID
						}
					}
				}
			}
			refUserID = lender.ID
			checkIncentiveTime = loan.OfferStartedAt
		}
	default:
		{
			return errs.NewError(errs.ErrBadRequest)
		}
	}
	if userID > 0 {
		user, err := s.ud.FirstByID(
			tx,
			userID,
			map[string][]interface{}{},
			false,
		)
		if err != nil {
			return errs.NewError(err)
		}
		ipdM, err := s.ipdd.First(
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
				"user_rank = ? or user_rank = ?": []interface{}{
					models.UserRankAll,
					user.Rank,
				},
				"type = ?":                 []interface{}{incentiveTransactionType},
				"loan_valid_duration <= ?": []interface{}{uint(loan.ValidAt.Sub(*loan.StartedAt).Seconds())},
			},
			map[string][]interface{}{
				"IncentiveProgram": []interface{}{},
			},
			[]string{
				"amount desc",
			},
		)
		if err != nil {
			return errs.NewError(err)
		}
		if ipdM != nil {
			if !(ipdM.UserRank == models.UserRankAll ||
				ipdM.UserRank == user.Rank) {
				return errs.NewError(errs.ErrBadRequest)
			}
			ipM := ipdM.IncentiveProgram
			itM, err := s.itd.First(
				tx,
				map[string][]interface{}{
					"incentive_program_id = ?": []interface{}{ipM.ID},
					"type = ?":                 []interface{}{ipdM.Type},
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
				if ipdM.RevokeTypes != "" {
					revokeTypes := strings.Split(string(ipdM.RevokeTypes), ",")
					for _, revokeType := range revokeTypes {
						itM, err = s.itd.First(
							tx,
							map[string][]interface{}{
								"incentive_program_id = ?":  []interface{}{ipdM.IncentiveProgramID},
								"type = ?":                  []interface{}{revokeType},
								"loan_id = ?":               []interface{}{loan.ID},
								"status = ?":                []interface{}{models.IncentiveTransactionStatusLocked},
								"lock_until_at is not null": []interface{}{},
								"lock_until_at <= ?":        []interface{}{time.Now()},
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
					}
				}
				if ipdM.UnlockTypes != "" {
					unlockTypes := strings.Split(string(ipdM.UnlockTypes), ",")
					for _, unlockType := range unlockTypes {
						itM, err = s.itd.First(
							tx,
							map[string][]interface{}{
								"incentive_program_id = ?": []interface{}{ipdM.IncentiveProgramID},
								"type = ?":                 []interface{}{unlockType},
								"loan_id = ?":              []interface{}{loan.ID},
								"status = ?":               []interface{}{models.IncentiveTransactionStatusLocked},
							},
							map[string][]interface{}{},
							[]string{},
						)
						if err != nil {
							return errs.NewError(err)
						}
						if itM != nil {
							err = s.incentiveForUnlock(
								tx,
								itM.ID,
								true,
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
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
					transactionStatus := models.IncentiveTransactionStatusLocked
					lockUntilAt := helpers.TimeAdd(*checkIncentiveTime, time.Duration(ipdM.LockDuration)*time.Second)
					if !lockUntilAt.After(time.Now()) {
						transactionStatus = models.IncentiveTransactionStatusDone
					}
					itM = &models.IncentiveTransaction{
						Network:            ipM.Network,
						IncentiveProgramID: ipM.ID,
						Type:               ipdM.Type,
						UserID:             user.ID,
						CurrencyID:         currencyID,
						LoanID:             loanID,
						Amount:             amount,
						LockUntilAt:        lockUntilAt,
						UnlockedAt:         nil,
						Status:             transactionStatus,
						RefUserID:          refUserID,
					}
					err = s.itd.Create(
						tx,
						itM,
					)
					if err != nil {
						return errs.NewError(err)
					}
					reference := fmt.Sprintf("it_%d_locked", itM.ID)
					switch itM.Status {
					case models.IncentiveTransactionStatusDone:
						{
							reference = fmt.Sprintf("it_%d_done", itM.ID)
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
					var isLock bool
					switch itM.Status {
					case models.IncentiveTransactionStatusLocked:
						{
							isLock = true
						}
					}
					switch itM.Type {
					case models.IncentiveTransactionTypeBorrowerLoanDelisted:
						{
							isLock = true
						}
					}
					err = s.transactionUserBalance(
						tx,
						ipM.Network,
						itM.UserID,
						itM.CurrencyID,
						itM.Amount,
						isLock,
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
	return nil
}

func (s *NftLend) JobIncentiveStatus(ctx context.Context) error {
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
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, itM := range itMs {
		err = s.IncentiveForReward(ctx, itM.ID)
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
		999999,
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
			err := s.incentiveForUnlock(
				tx,
				transactionID,
				false,
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

func (s *NftLend) incentiveForUnlock(tx *gorm.DB, transactionID uint, checked bool) error {
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
		models.IncentiveTransactionTypeLenderLoanMatched,
		models.IncentiveTransactionTypeUserAirdropReward,
		models.IncentiveTransactionTypeUserAmaReward:
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
	if !checked && itM.LockUntilAt.After(time.Now()) {
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
}

func (s *NftLend) IncentiveForReward(ctx context.Context, transactionID uint) error {
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
			if itM.UserID <= 0 {
				user, err := s.getUser(
					tx,
					itM.Network,
					itM.Address,
					false,
				)
				if err != nil {
					return errs.NewError(err)
				}
				itM.UserID = user.ID
			}
			var reference string
			if itM.LockUntilAt != nil {
				itM.Status = models.IncentiveTransactionStatusLocked
				reference = fmt.Sprintf("it_%d_locked", itM.ID)
			} else {
				itM.Status = models.IncentiveTransactionStatusDone
				reference = fmt.Sprintf("it_%d_done", itM.ID)
			}
			err = s.itd.Save(
				tx,
				itM,
			)
			if err != nil {
				return errs.NewError(err)
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
			var isLocked bool
			switch itM.Status {
			case models.IncentiveTransactionStatusLocked:
				{
					isLocked = true
				}
			}
			err = s.transactionUserBalance(
				tx,
				itM.Network,
				itM.UserID,
				itM.CurrencyID,
				itM.Amount,
				isLocked,
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
