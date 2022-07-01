package services

import (
	"context"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) UserVerifyEmail(ctx context.Context, req *serializers.UserVerifyEmailReq) error {
	req.Email = strings.ToLower(req.Email)
	var vM *models.Verification
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err := s.getUser(tx, req.Network, req.Address, false)
			if err != nil {
				return errs.NewError(err)
			}
			if user.IsVerified {
				if user.Email == req.Email {
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			vMs, err := s.vd.Find(
				tx,
				map[string][]interface{}{
					"network = ?": []interface{}{req.Network},
					"user_id = ?": []interface{}{user.ID},
					"type = ?":    []interface{}{models.VerificationTypeEmail},
					"status = ?":  []interface{}{models.VerificationStatusVerifying},
				},
				map[string][]interface{}{},
				[]string{},
				0,
				999999,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for _, vM := range vMs {
				if vM.ExpiredAt.After(time.Now().Add(-10 * time.Minute)) {
					return errs.NewError(errs.ErrVerificationLimited)
				}
				vM.Status = models.VerificationStatusCancelled
				err = s.vd.Save(
					tx,
					vM,
				)
				if err != nil {
					return errs.NewError(err)
				}
			}
			vM = &models.Verification{
				Network:   user.Network,
				UserID:    user.ID,
				Type:      models.VerificationTypeEmail,
				Email:     req.Email,
				Token:     helpers.RandomStringWithLength(32),
				ExpiredAt: helpers.TimeAdd(time.Now(), 15*time.Minute),
				Status:    models.VerificationStatusVerifying,
			}
			err = s.vd.Create(
				tx,
				vM,
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
	if vM != nil {
		s.EmailForEmailVerification(
			ctx,
			vM.ID,
		)
	}
	return nil
}

func (s *NftLend) UserVerifyEmailToken(ctx context.Context, req *serializers.UserVerifyEmailTokenReq) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			vM, err := s.vd.First(
				tx,
				map[string][]interface{}{
					"email = ?":  []interface{}{req.Email},
					"type = ?":   []interface{}{models.VerificationTypeEmail},
					"token = ?":  []interface{}{req.Token},
					"status = ?": []interface{}{models.VerificationStatusVerifying},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if vM == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			if vM.ExpiredAt.Before(time.Now()) {
				return errs.NewError(errs.ErrVerificationExpired)
			}
			vM.Status = models.VerificationStatusVerified
			err = s.vd.Save(
				tx,
				vM,
			)
			if err != nil {
				return errs.NewError(err)
			}
			user, err := s.ud.FirstByID(
				tx,
				vM.UserID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			user.IsVerified = true
			user.LoanNotiEnabled = true
			user.NewsNotiEnabled = true
			user.Email = vM.Email
			err = s.ud.Save(
				tx,
				user,
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
