package services

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) GetUser(ctx context.Context, network models.Network, address string) (*models.User, error) {
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
			user, err = s.getUser(tx, network, address)
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

func (s *NftLend) getUser(tx *gorm.DB, network models.Network, address string) (*models.User, error) {
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
		addressChecked := strings.ToLower(strings.TrimSpace(address))
		user = &models.User{
			Network:         network,
			Address:         address,
			AddressChecked:  addressChecked,
			UserName:        addressChecked,
			NewsNotiEnabled: false,
			LoanNotiEnabled: false,
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

func (s *NftLend) UserConnected(ctx context.Context, network models.Network, address string, referrerCode string) (*models.User, error) {
	referrerCode = strings.ToLower(referrerCode)
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
			user, err = s.getUser(tx, network, address)
			if err != nil {
				return errs.NewError(err)
			}
			// check wallet connected
			if referrerCode != "" {
				user.ReferrerCode = referrerCode
				referrer, err := s.ud.First(
					tx,
					map[string][]interface{}{
						"network = ?":   []interface{}{network},
						"user_name = ?": []interface{}{referrerCode},
					},
					map[string][]interface{}{},
					[]string{},
				)
				if err != nil {
					return errs.NewError(err)
				}
				if referrer != nil {
					user.ReferrerUserID = referrer.ID
				}
				err = s.ud.Save(
					tx,
					user,
				)
				if err != nil {
					return errs.NewError(err)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return user, nil
}

func (s *NftLend) UserGetSettings(ctx context.Context, network models.Network, address string) (*models.User, error) {
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
			user, err = s.getUser(tx, network, address)
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
			user, err = s.getUser(tx, req.Network, req.Address)
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

func (s *NftLend) GetUserStats(ctx context.Context, network models.Network, address string) (*models.UserBorrowStats, *models.UserLendStats, error) {
	var user *models.User
	var err error
	if address == "" {
		return nil, nil, errs.NewError(errs.ErrBadRequest)
	}
	var lendStats *models.UserLendStats
	var borrowStats *models.UserBorrowStats
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
			return nil, nil, errs.NewError(errs.ErrBadRequest)
		}
	}
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err = s.getUser(tx, network, address)
			if err != nil {
				return errs.NewError(err)
			}
			borrowStats, err = s.ud.GetUserBorrowStats(
				tx,
				user.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			lendStats, err = s.ud.GetUserLendStats(
				tx,
				user.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, nil, errs.NewError(err)
	}
	return borrowStats, lendStats, nil
}

func (s *NftLend) GetUserBalances(ctx context.Context, network models.Network, address string) ([]*models.UserBalance, error) {
	user, err := s.GetUser(ctx, network, address)
	if err != nil {
		return nil, errs.NewError(err)
	}
	userBalances, err := s.ubd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"user_id = ?": []interface{}{user.ID},
		},
		map[string][]interface{}{},
		[]string{"id desc"},
		0,
		9999,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return userBalances, nil
}

func (s *NftLend) GetUserBalanceTransactions(ctx context.Context, network models.Network, address string, currencyID uint, page int, limit int) ([]*models.UserBalanceTransaction, uint, error) {
	user, err := s.GetUser(ctx, network, address)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	filters := map[string][]interface{}{
		"user_id = ?": []interface{}{user.ID},
	}
	if currencyID > 0 {
		filters["currency_id = ?"] = []interface{}{currencyID}
	}
	userBalanceTxns, count, err := s.ubtd.Find4Page(
		daos.GetDBMainCtx(ctx),
		filters,
		map[string][]interface{}{
			"User":                 []interface{}{},
			"Currency":             []interface{}{},
			"IncentiveTransaction": []interface{}{},
		},
		[]string{"id desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return userBalanceTxns, count, nil
}

func (s *NftLend) GetUserPWPTokenBalance(ctx context.Context, network models.Network, address string) (*models.UserBalance, error) {
	var userBalance *models.UserBalance
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err := s.getUser(tx, network, address)
			if err != nil {
				return errs.NewError(err)
			}
			userBalance, err = s.getUserPWPTokenBalance(
				tx,
				user.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			userBalance, err = s.ubd.FirstByID(
				tx,
				userBalance.ID,
				map[string][]interface{}{
					"User":     []interface{}{},
					"Currency": []interface{}{},
				},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return userBalance, nil
}

func (s *NftLend) getUserPWPTokenBalance(tx *gorm.DB, userID uint) (*models.UserBalance, error) {
	pwpToken, err := s.getLendCurrencyBySymbol(
		tx,
		models.SymbolPWPToken,
		models.NetworkNEAR,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if pwpToken == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	userBalance, err := s.getUserBalance(
		tx,
		userID,
		pwpToken.ID,
		false,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return userBalance, nil
}

func (s *NftLend) getUserBalance(tx *gorm.DB, userID uint, currencyID uint, forUpdate bool) (*models.UserBalance, error) {
	userBalance, err := s.ubd.First(
		tx,
		map[string][]interface{}{
			"user_id = ?":     []interface{}{userID},
			"currency_id = ?": []interface{}{currencyID},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if userBalance == nil {
		user, err := s.ud.FirstByID(
			tx,
			userID,
			map[string][]interface{}{},
			false,
		)
		if err != nil {
			return nil, errs.NewError(err)
		}
		userBalance = &models.UserBalance{
			Network:       user.Network,
			UserID:        user.ID,
			CurrencyID:    currencyID,
			Balance:       numeric.BigFloat{*big.NewFloat(0)},
			LockedBalance: numeric.BigFloat{*big.NewFloat(0)},
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

func (s *NftLend) transactionUserBalance(tx *gorm.DB, network models.Network, userID uint, currencyID uint, amount numeric.BigFloat, isLocked bool, isClaimed bool, reference string) error {
	userBalance, err := s.getUserBalance(
		tx,
		userID,
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
	if isClaimed {
		claimTx := &models.UserBalanceHistory{
			Network:       userBalance.Network,
			Type:          models.UserBalanceHistoryTypeClaimedBalance,
			UserBalanceID: userBalance.ID,
			CurrencyID:    currencyID,
			Amount:        numeric.BigFloat{*models.NegativeBigFloat(&amount.Float)},
			Reference:     reference,
		}
		err = s.ubhd.Create(
			tx,
			claimTx,
		)
		if err != nil {
			return errs.NewError(err)
		}
		newClaimedBalance := models.AddBigFloats(&userBalance.ClaimedBalance.Float, &claimTx.Amount.Float)
		userBalance.ClaimedBalance = numeric.BigFloat{*newClaimedBalance}
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

func (s *NftLend) unlockUserBalance(tx *gorm.DB, userID uint, currencyID uint, amount numeric.BigFloat, reference string) error {
	userBalance, err := s.getUserBalance(
		tx,
		userID,
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

func (s *NftLend) ClaimUserBalance(ctx context.Context, req *serializers.ClaimUserBalanceReq) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			userBalance, err := s.getUserBalance(
				tx,
				req.UserID,
				req.CurrencyID,
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if userBalance.UpdatedAt.Unix() != req.Timestamp {
				return errs.NewError(errs.ErrBadRequest)
			}
			if userBalance.GetAvaiableBalance().Cmp(&req.Amount.Float) < 0 {
				return errs.NewError(errs.ErrBadRequest)
			}
			user, err := s.ud.FirstByID(
				tx,
				userBalance.UserID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if !strings.EqualFold(user.Address, req.ToAddress) {
				return errs.NewError(errs.ErrBadRequest)
			}
			currency, err := s.cd.FirstByID(
				tx,
				userBalance.CurrencyID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if !currency.ClaimEnabled {
				return errs.NewError(errs.ErrBadRequest)
			}
			// validate request by sig
			message := fmt.Sprintf("%s-%s-%s-%d", strings.ToLower(user.Address), strings.ToLower(currency.ContractAddress), req.Amount.ToString(), req.Timestamp)
			err = s.bcs.Near.ValidateMessageSignature(
				s.conf.Contract.NearNftypawnAddress,
				message,
				req.Signature,
				user.Address,
			)
			if err != nil {
				return errs.NewError(err)
			}
			//
			userBalanceTransaction := &models.UserBalanceTransaction{
				Network:       userBalance.Network,
				UserID:        userBalance.UserID,
				UserBalanceID: userBalance.ID,
				CurrencyID:    userBalance.CurrencyID,
				Type:          models.UserBalanceTransactionTypeClaim,
				ToAddress:     req.ToAddress,
				Amount:        numeric.BigFloat{*models.NegativeBigFloat(&req.Amount.Float)},
				Message:       message,
				Signature:     req.Signature,
				Status:        models.UserBalanceTransactionStatusDone,
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
				userBalance.Network,
				userBalance.UserID,
				userBalance.CurrencyID,
				userBalanceTransaction.Amount,
				false,
				true,
				fmt.Sprintf("ubt_%d_claim", userBalanceTransaction.ID),
			)
			if err != nil {
				return errs.NewError(err)
			}
			hash, err := s.bcs.Near.FtTransfer(
				currency.ContractAddress,
				currency.PoolAddress,
				userBalanceTransaction.ToAddress,
				models.ConvertBigFloatToWei(&userBalanceTransaction.Amount.Float, currency.Decimals),
			)
			if err != nil {
				return errs.NewError(err)
			}
			userBalanceTransaction.TxHash = hash
			err = s.ubtd.Save(
				tx,
				userBalanceTransaction,
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
