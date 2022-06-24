package services

import (
	"context"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/services/3rd/mailer"
)

func (s *NftLend) JobEmailSchedule(ctx context.Context) error {
	var retErr error
	err := s.JobEmailScheduleBorrowerLoanRemind7(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.JobEmailScheduleBorrowerLoanRemind3(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.JobEmailScheduleBorrowerLoanRemind1(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.JobEmailScheduleBorrowerLoanLiquidated(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	err = s.JobEmailScheduleLenderLoanLiquidated(ctx)
	if err != nil {
		retErr = errs.MergeError(retErr, err)
	}
	return retErr
}

func (s *NftLend) JobEmailScheduleBorrowerLoanRemind7(ctx context.Context) error {
	var retErr error
	loans, err := s.ld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":            []interface{}{models.LoanStatusCreated},
			"offer_expired_at < ?":  []interface{}{time.Now().Add((7 * 24) * time.Hour)},
			"offer_expired_at >= ?": []interface{}{time.Now().Add((7*24 - 1) * time.Hour)},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, loan := range loans {
		err = s.EmailForBorrowerLoanRemind7(ctx, loan.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, errs.NewErrorWithId(err, loan.ID))
		}
	}
	return retErr
}

func (s *NftLend) JobEmailScheduleBorrowerLoanRemind3(ctx context.Context) error {
	var retErr error
	loans, err := s.ld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":            []interface{}{models.LoanStatusCreated},
			"offer_expired_at < ?":  []interface{}{time.Now().Add((3 * 24) * time.Hour)},
			"offer_expired_at >= ?": []interface{}{time.Now().Add((3*24 - 1) * time.Hour)},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, loan := range loans {
		err = s.EmailForBorrowerLoanRemind3(ctx, loan.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, errs.NewErrorWithId(err, loan.ID))
		}
	}
	return retErr
}

func (s *NftLend) JobEmailScheduleBorrowerLoanRemind1(ctx context.Context) error {
	var retErr error
	loans, err := s.ld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":            []interface{}{models.LoanStatusCreated},
			"offer_expired_at < ?":  []interface{}{time.Now().Add((1 * 24) * time.Hour)},
			"offer_expired_at >= ?": []interface{}{time.Now().Add((1*24 - 1) * time.Hour)},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, loan := range loans {
		err = s.EmailForBorrowerLoanRemind1(ctx, loan.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, errs.NewErrorWithId(err, loan.ID))
		}
	}
	return retErr
}

func (s *NftLend) JobEmailScheduleBorrowerLoanLiquidated(ctx context.Context) error {
	var retErr error
	loans, err := s.ld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":            []interface{}{models.LoanStatusCreated},
			"offer_overdue_at >= ?": []interface{}{time.Now().Add((-1) * time.Hour)},
			"offer_overdue_at < ?":  []interface{}{time.Now()},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, loan := range loans {
		err = s.EmailForBorrowerLoanLiquidated(ctx, loan.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, errs.NewErrorWithId(err, loan.ID))
		}
	}
	return retErr
}

func (s *NftLend) JobEmailScheduleLenderLoanLiquidated(ctx context.Context) error {
	var retErr error
	loans, err := s.ld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"status = ?":            []interface{}{models.LoanStatusCreated},
			"offer_overdue_at >= ?": []interface{}{time.Now().Add((-1) * time.Hour)},
			"offer_overdue_at < ?":  []interface{}{time.Now()},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, loan := range loans {
		err = s.EmailForLenderLoanLiquidated(ctx, loan.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, errs.NewErrorWithId(err, loan.ID))
		}
	}
	return retErr
}

func (s *NftLend) getCurrencyMap(ctx context.Context, m *models.Currency) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"id":               m.ID,
		"created_at":       m.CreatedAt,
		"updated_at":       m.UpdatedAt,
		"network":          m.Network,
		"contract_address": m.ContractAddress,
		"decimals":         m.Decimals,
		"symbol":           m.Symbol,
		"name":             m.Name,
		"icon_url":         m.IconURL,
	}
}

func (s *NftLend) getAssetMap(ctx context.Context, m *models.Asset) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"id":                      m.ID,
		"created_at":              m.CreatedAt,
		"updated_at":              m.UpdatedAt,
		"network":                 m.Network,
		"collection_id":           m.CollectionID,
		"seo_url":                 m.SeoURL,
		"contract_address":        m.ContractAddress,
		"token_url":               m.TokenURL,
		"token_id":                m.TokenID,
		"name":                    m.Name,
		"seller_fee_rate":         m.SellerFeeRate,
		"origin_network":          m.OriginNetwork,
		"origin_contract_address": m.OriginContractAddress,
		"origin_token_id":         m.OriginTokenID,
	}
}

func (s *NftLend) getLoanMap(ctx context.Context, m *models.Loan) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"id":                           m.ID,
		"created_at":                   m.CreatedAt,
		"updated_at":                   m.UpdatedAt,
		"network":                      m.Network,
		"owner":                        m.Owner,
		"lender":                       m.Lender,
		"asset":                        s.getAssetMap(ctx, m.Asset),
		"currency":                     s.getCurrencyMap(ctx, m.Currency),
		"started_at":                   models.FormatEmailTime(m.StartedAt),
		"duration":                     models.FormatFloatNumber("%.2f", models.DivFloats(float64(m.Duration), 60*60*24)),
		"expired_at":                   models.FormatEmailTime(m.ExpiredAt),
		"finished_at":                  models.FormatEmailTime(m.FinishedAt),
		"principal_amount":             models.FormatStringNumber(m.PrincipalAmount.Float.Text('f', 8)),
		"interest_rate":                models.FormatFloatNumber("%.18f", m.InterestRate*100),
		"fee_rate":                     models.FormatFloatNumber("%.18f", m.FeeRate),
		"fee_amount":                   m.FeeAmount,
		"status":                       m.Status,
		"offer_principal_amount":       models.FormatStringNumber(m.OfferPrincipalAmount.Float.Text('f', 8)),
		"offer_interest_rate":          models.FormatFloatNumber("%.18f", m.OfferInterestRate*100),
		"offer_duration":               models.FormatFloatNumber("%.2f", models.DivFloats(float64(m.OfferDuration), 60*60*24)),
		"offer_started_at":             models.FormatEmailTime(m.OfferStartedAt),
		"offer_expired_at":             models.FormatEmailTime(m.OfferExpiredAt),
		"matured_offer_payment_amount": models.FormatStringNumber(m.MaturedOfferPaymentAmount().Text('f', 8)),
		"early_offer_payment_amount":   models.FormatStringNumber(m.EarlyOfferPaymentAmount().Text('f', 8)),
	}
}

func (s *NftLend) getLoanOfferMap(ctx context.Context, m *models.LoanOffer) map[string]interface{} {
	if m == nil {
		return map[string]interface{}{}
	}
	return map[string]interface{}{
		"id":                     m.ID,
		"created_at":             m.CreatedAt,
		"updated_at":             m.UpdatedAt,
		"network":                m.Network,
		"lender":                 m.Lender,
		"started_at":             models.FormatEmailTime(m.StartedAt),
		"duration":               models.FormatFloatNumber("%.2f", models.DivFloats(float64(m.Duration), 60*60*24)),
		"expired_at":             models.FormatEmailTime(m.ExpiredAt),
		"finished_at":            models.FormatEmailTime(m.FinishedAt),
		"principal_amount":       models.FormatStringNumber(m.PrincipalAmount.Float.Text('f', 8)),
		"interest_rate":          models.FormatFloatNumber("%.18f", m.InterestRate*100),
		"status":                 m.Status,
		"matured_payment_amount": models.FormatStringNumber(m.MaturedOfferPaymentAmount().Text('f', 8)),
		"early_payment_amount":   models.FormatStringNumber(m.EarlyOfferPaymentAmount().Text('f', 8)),
		"loan":                   s.getLoanMap(ctx, m.Loan),
	}
}

func (s *NftLend) sendEmailToUser(ctx context.Context, address string, network models.Network, emailType string, reqMap map[string]interface{}) error {
	err := s.CreateNotification(
		ctx,
		network,
		address,
		models.NotificationType(emailType),
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	user, err := s.GetUser(ctx, network, address)
	if err != nil {
		return errs.NewError(err)
	}
	if user.LoanNotiEnabled {
		if user.Email != "" {
			reqMap["web_url"] = s.conf.WebUrl
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
	}
	return nil
}

func (s *NftLend) sendEmailToEmail(ctx context.Context, email string, emailType string, reqMap map[string]interface{}) error {
	reqMap["web_url"] = s.conf.WebUrl
	err := mailer.Send(
		"hello@nftpawn.financial",
		"Admin",
		email,
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
	return nil
}

func (s *NftLend) EmailForReference(ctx context.Context, emailQuueue []*models.EmailQueue) error {
	var retErr error
	for _, q := range emailQuueue {
		var err error
		switch q.EmailType {
		case models.EMAIL_BORROWER_OFFER_NEW:
			{
				err = s.EmailForBorrowerOfferNew(ctx, q.ObjectID)
			}
		case models.EMAIL_BORROWER_LOAN_STARTED:
			{
				err = s.EmailForBorrowerLoanStarted(ctx, q.ObjectID)
			}
		case models.EMAIL_BORROWER_LOAN_REMIND7:
			{
				err = s.EmailForBorrowerLoanRemind7(ctx, q.ObjectID)
			}
		case models.EMAIL_BORROWER_LOAN_REMIND3:
			{
				err = s.EmailForBorrowerLoanRemind7(ctx, q.ObjectID)
			}
		case models.EMAIL_BORROWER_LOAN_REMIND1:
			{
				err = s.EmailForBorrowerLoanRemind7(ctx, q.ObjectID)
			}
		case models.EMAIL_BORROWER_LOAN_LIQUIDATED:
			{
				err = s.EmailForBorrowerLoanLiquidated(ctx, q.ObjectID)
			}
		case models.EMAIL_LENDER_OFFER_STARTED:
			{
				err = s.EmailForLenderOfferStarted(ctx, q.ObjectID)
			}
		case models.EMAIL_LENDER_LOAN_REPAID:
			{
				err = s.EmailForLenderLoanRepaid(ctx, q.ObjectID)
			}
		case models.EMAIL_LENDER_LOAN_LIQUIDATED:
			{
				err = s.EmailForLenderLoanLiquidated(ctx, q.ObjectID)
			}
		}
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	return retErr
}

func (s *NftLend) EmailForBorrowerOfferNew(ctx context.Context, offerID uint) error {
	offer, err := s.lod.FirstByID(
		daos.GetDBMainCtx(ctx),
		offerID,
		map[string][]interface{}{
			"Loan.Asset":    []interface{}{},
			"Loan.Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"offer": s.getLoanOfferMap(ctx, offer),
	}
	network := offer.Loan.Network
	address := offer.Loan.Owner
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_BORROWER_OFFER_NEW,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForBorrowerLoanStarted(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Owner
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_BORROWER_LOAN_STARTED,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForBorrowerLoanRemind7(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Owner
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_BORROWER_LOAN_REMIND7,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForBorrowerLoanRemind3(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Owner
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_BORROWER_LOAN_REMIND3,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForBorrowerLoanRemind1(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Owner
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_BORROWER_LOAN_REMIND1,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForBorrowerLoanLiquidated(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Owner
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_BORROWER_LOAN_LIQUIDATED,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForLenderOfferStarted(ctx context.Context, offerID uint) error {
	offer, err := s.lod.FirstByID(
		daos.GetDBMainCtx(ctx),
		offerID,
		map[string][]interface{}{
			"Loan.Asset":    []interface{}{},
			"Loan.Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"offer": s.getLoanOfferMap(ctx, offer),
	}
	network := offer.Loan.Network
	address := offer.Loan.Lender
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_LENDER_OFFER_STARTED,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForLenderLoanRepaid(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Lender
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_LENDER_LOAN_REPAID,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForLenderLoanLiquidated(ctx context.Context, loanID uint) error {
	loan, err := s.ld.FirstByID(
		daos.GetDBMainCtx(ctx),
		loanID,
		map[string][]interface{}{
			"Asset":    []interface{}{},
			"Currency": []interface{}{},
		},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"loan": s.getLoanMap(ctx, loan),
	}
	network := loan.Network
	address := loan.Lender
	err = s.sendEmailToUser(
		ctx,
		address,
		network,
		models.EMAIL_LENDER_LOAN_LIQUIDATED,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) EmailForEmailVerification(ctx context.Context, vID uint) error {
	vM, err := s.vd.FirstByID(
		daos.GetDBMainCtx(ctx),
		vID,
		map[string][]interface{}{},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	reqMap := map[string]interface{}{
		"email": vM.Email,
		"token": vM.Token,
	}
	err = s.sendEmailToEmail(
		ctx,
		vM.Email,
		models.EMAIL_USER_VERIFY_EMAIL,
		reqMap,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}
