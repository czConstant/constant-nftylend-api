package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/services/3rd/mailer"
)

func (s *NftLend) JobEmailSchedule(ctx context.Context) error {
	var retErr error
	return retErr
}

func (s *NftLend) sendEmailToUser(ctx context.Context, address string, network models.Network, emailType string, reqMap interface{}) error {
	user, err := s.GetUser(ctx, address, network)
	if err != nil {
		return errs.NewError(err)
	}
	if user.Email != "" {
		err := mailer.Send(
			"hello@nftpawn.financial",
			"Admin",
			user.Email,
			"",
			emailType,
			"en",
			reqMap,
			[]string{},
			[]string{},
		)
		if err != nil {
			return errs.NewError(err)
		}
	}
	return nil
}

func (s *NftLend) EmailForBorrowerOfferNew(ctx context.Context, offerID uint) error {
	reqMap := map[string]interface{}{}
	err := s.sendEmailToUser(
		ctx,
		"",
		models.NetworkBSC,
		models.EMAIL_BORROWER_NEW_OFFER,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForBorrowerLoanRemind(ctx context.Context, loanID uint) error {
	return nil
}

func (s *NftLend) EmailForBorrowerLoanStarted(ctx context.Context, loanID uint) error {
	return nil
}

func (s *NftLend) EmailForBorrowerLoanLiquidated(ctx context.Context, loanID uint) error {
	return nil
}

func (s *NftLend) EmailForLenderOfferStarted(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) EmailForLenderLoanRepaid(ctx context.Context, loanID uint) error {
	return nil
}

func (s *NftLend) EmailForLenderLoanLiquidated(ctx context.Context, loanID uint) error {
	return nil
}
