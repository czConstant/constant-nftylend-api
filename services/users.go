package services

import (
	"context"
	"math/big"
	"strings"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) GetUser(ctx context.Context, address string, network models.Network) (*models.User, error) {
	var user *models.User
	var err error
	if address == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	switch network {
	case models.NetworkSOL,
		models.NetworkAVAX,
		models.NetworkBOBA,
		models.NetworkBSC,
		models.NetworkETH,
		models.NetworkMATIC,
		models.NetworkNEAR:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err = s.getUser(tx, address, network)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return user, nil
}

func (s *NftLend) getUser(tx *gorm.DB, address string, network models.Network) (*models.User, error) {
	if address == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	switch network {
	case models.NetworkSOL,
		models.NetworkAVAX,
		models.NetworkBOBA,
		models.NetworkBSC,
		models.NetworkETH,
		models.NetworkMATIC,
		models.NetworkNEAR:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	user, err := s.ud.First(
		tx,
		map[string][]interface{}{
			"network = ?":         []interface{}{network},
			"address_checked = ?": []interface{}{strings.ToLower(strings.TrimSpace(address))},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if user == nil {
		user = &models.User{
			Network:         network,
			Address:         address,
			AddressChecked:  strings.ToLower(strings.TrimSpace(address)),
			NewsNotiEnabled: true,
			LoanNotiEnabled: true,
		}
		err = s.ud.Create(
			tx,
			user,
		)
		if err != nil {
			return nil, errs.NewError(err)
		}
	}
	return user, nil
}

func (s *NftLend) UserGetSettings(ctx context.Context, address string, network models.Network) (*models.User, error) {
	var user *models.User
	var err error
	if address == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	switch network {
	case models.NetworkSOL,
		models.NetworkAVAX,
		models.NetworkBOBA,
		models.NetworkBSC,
		models.NetworkETH,
		models.NetworkMATIC,
		models.NetworkNEAR:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err = s.getUser(tx, address, network)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return user, nil
}

func (s *NftLend) UserUpdateSetting(ctx context.Context, req *serializers.UpdateUserSettingReq) (*models.User, error) {
	var user *models.User
	var err error
	if req.Address == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	if req.Email == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	switch req.Network {
	case models.NetworkSOL,
		models.NetworkAVAX,
		models.NetworkBOBA,
		models.NetworkBSC,
		models.NetworkETH,
		models.NetworkMATIC,
		models.NetworkNEAR:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err = s.getUser(tx, req.Address, req.Network)
			if err != nil {
				return errs.NewError(err)
			}
			user.Email = req.Email
			user.NewsNotiEnabled = req.NewsNotiEnabled
			user.LoanNotiEnabled = req.LoanNotiEnabled
			if err != nil {
				return errs.NewError(err)
			}
			err = s.ud.Save(tx, user)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return user, nil
}

func (s *NftLend) getUserBalance(tx *gorm.DB, network models.Network, address string, currencyID uint, forUpdate bool) (*models.UserBalance, error) {
	if address == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	switch network {
	case models.NetworkSOL,
		models.NetworkAVAX,
		models.NetworkBOBA,
		models.NetworkBSC,
		models.NetworkETH,
		models.NetworkMATIC,
		models.NetworkNEAR:
		{
		}
	default:
		{
			return nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	userBalance, err := s.ubd.First(
		tx,
		map[string][]interface{}{
			"network = ?":         []interface{}{network},
			"address_checked = ?": []interface{}{strings.ToLower(strings.TrimSpace(address))},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if userBalance == nil {
		userBalance = &models.UserBalance{
			Network:        network,
			Address:        address,
			AddressChecked: strings.ToLower(strings.TrimSpace(address)),
			CurrencyID:     currencyID,
			Balance:        numeric.BigFloat{*big.NewFloat(0)},
			LockedBalance:  numeric.BigFloat{*big.NewFloat(0)},
		}
		err = s.ubd.Create(
			tx,
			userBalance,
		)
		if err != nil {
			return nil, errs.NewError(err)
		}
	}
	if forUpdate {
		userBalance, err = s.ubd.FirstByID(
			tx,
			userBalance.ID,
			map[string][]interface{}{},
			true,
		)
		if err != nil {
			return nil, errs.NewError(err)
		}
	}
	return userBalance, nil
}

func (s *NftLend) transactionUserBalance(tx *gorm.DB, network models.Network, address string, currencyID uint, amount numeric.BigFloat, isLocked bool, reference string) error {
	userBalance, err := s.getUserBalance(
		tx,
		network,
		address,
		currencyID,
		true,
	)
	if err != nil {
		return errs.NewError(err)
	}
	err = s.ubhd.Create(
		tx,
		&models.UserBalanceHistory{
			Network:       userBalance.Network,
			Type:          models.UserBalanceHistoryTypeBalance,
			UserBalanceID: userBalance.ID,
			CurrencyID:    currencyID,
			Amount:        amount,
			Reference:     reference,
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	newBalance := models.AddBigFloats(&userBalance.Balance.Float, &amount.Float)
	userBalance.Balance = numeric.BigFloat{*newBalance}
	if isLocked {
		err = s.ubhd.Create(
			tx,
			&models.UserBalanceHistory{
				Network:       userBalance.Network,
				Type:          models.UserBalanceHistoryTypeLockedBalance,
				UserBalanceID: userBalance.ID,
				CurrencyID:    currencyID,
				Amount:        amount,
				Reference:     reference,
			},
		)
		if err != nil {
			return errs.NewError(err)
		}
		newLockedBalance := models.AddBigFloats(&userBalance.LockedBalance.Float, &amount.Float)
		userBalance.LockedBalance = numeric.BigFloat{*newLockedBalance}
	}
	err = s.ubd.Save(
		tx,
		userBalance,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) unlockUserBalance(tx *gorm.DB, network models.Network, address string, currencyID uint, amount numeric.BigFloat, reference string) error {
	userBalance, err := s.getUserBalance(
		tx,
		network,
		address,
		currencyID,
		true,
	)
	if err != nil {
		return errs.NewError(err)
	}
	amount = numeric.BigFloat{*models.SubBigFloats(big.NewFloat(0), &amount.Float)}
	err = s.ubhd.Create(
		tx,
		&models.UserBalanceHistory{
			Network:       userBalance.Network,
			Type:          models.UserBalanceHistoryTypeLockedBalance,
			UserBalanceID: userBalance.ID,
			CurrencyID:    currencyID,
			Amount:        amount,
			Reference:     reference,
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	newLockedBalance := models.AddBigFloats(&userBalance.LockedBalance.Float, &amount.Float)
	userBalance.LockedBalance = numeric.BigFloat{*newLockedBalance}
	err = s.ubd.Save(
		tx,
		userBalance,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}
