package services

import "context"

func (s *NftLend) JobEmailSchedule(ctx context.Context) error {
	var retErr error
	return retErr
}

func (s *NftLend) emailForBorrowerNewOffer(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForBorrowerRemindPaybackLoan(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForBorrowerAcceptLoan(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForBorrowerLiquidateLoanDone(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForLenderRequestExtendLoan(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForLenderAcceptOffer(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForLenderPaybackLoanDone(ctx context.Context, offerID uint) error {
	return nil
}

func (s *NftLend) emailForLenderLiquidateLoanDone(ctx context.Context, offerID uint) error {
	return nil
}
