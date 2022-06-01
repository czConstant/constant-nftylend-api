package services

import (
	"fmt"
	"time"

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
						},
						map[string][]interface{}{},
						[]string{},
					)
					if err != nil {
						return errs.NewError(err)
					}
					if itM == nil {
						isOk = false
					}
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
					Status:             models.IncentiveTransactionStatusLocked,
				}
				err = s.itd.Create(
					tx,
					itM,
				)
				if err != nil {
					return errs.NewError(err)
				}
				err = s.transactionUserBalance(
					tx,
					ipM.Network,
					itM.Address,
					itM.CurrencyID,
					itM.Amount,
					true,
					fmt.Sprintf("it_%d", itM.ID),
				)
				if err != nil {
					return errs.NewError(err)
				}
			}
		}
	}
	return nil
}
