package services

import (
	"context"
	"strings"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/jinzhu/gorm"
)

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
