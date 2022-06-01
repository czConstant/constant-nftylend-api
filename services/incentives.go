package services

import (
	"context"
	"fmt"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
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
	default:
		{
			return errs.NewError(errs.ErrBadRequest)
		}
	}
	ipdMs, err := s.ipdd.Find(
		tx,
		map[string][]interface{}{
			`exists(
				select 1
				from incentive_programs
				where incentive_program_details.incentive_program_id = incentive_programs.id
				  and (? between incentive_programs.start and incentive_programs.end)
				  and incentive_programs.status = ?
			)`: []interface{}{checkIncentiveTime, models.IncentiveProgramStatusActived},
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
	for _, ipdM := range ipdMs {
		ipM := ipdM.IncentiveProgram
		if uint(loan.ValidAt.Sub(*loan.StartedAt).Seconds()) >= ipM.LoanValidDuration {
			itM, err := s.itd.First(
				tx,
				map[string][]interface{}{
					"incentive_program_id = ?": []interface{}{ipM.ID},
					"type = ?":                 []interface{}{ipdM.Type},
					"address = ?":              []interface{}{address},
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
								"incentive_program_id = ?": []interface{}{ipdM.IncentiveProgramID},
								"type = ?":                 []interface{}{models.IncentiveTransactionTypeBorrowerLoanListed},
								"address = ?":              []interface{}{address},
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
					itM = &models.IncentiveTransaction{
						Network:            ipM.Network,
						IncentiveProgramID: ipM.ID,
						Type:               ipdM.Type,
						Address:            address,
						CurrencyID:         ipdM.IncentiveProgram.CurrencyID,
						LoanID:             loanID,
						Amount:             ipdM.Amount,
						LockUntilAt:        helpers.TimeAdd(*checkIncentiveTime, time.Duration(loan.OfferDuration)*time.Second),
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
					err = s.transactionUserBalance(
						tx,
						ipM.Network,
						itM.Address,
						itM.CurrencyID,
						itM.Amount,
						true,
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

func (s *NftLend) JobIncentiveForUnlock(ctx context.Context) error {
	var retErr error
	itMs, err := s.itd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"type in (?)": []interface{}{[]models.IncentiveTransactionType{
				models.IncentiveTransactionTypeBorrowerLoanListed,
				models.IncentiveTransactionTypeLenderLoanMatched,
			}},
			"status = ?":         []interface{}{models.IncentiveTransactionStatusLocked},
			"lock_until_at <= ?": []interface{}{time.Now()},
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
				itM.Network,
				itM.Address,
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