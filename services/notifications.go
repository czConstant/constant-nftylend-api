package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) CreateNotification(ctx context.Context, network models.Network, notiType models.NotificationType, address string, reqMap map[string]interface{}) error {
	notiTpl, err := s.ntd.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"type = ?":    []interface{}{notiType},
			"enabled = ?": []interface{}{true},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return errs.NewError(err)
	}
	if notiTpl != nil {
		noti, err := notiTpl.Execute(network, address, reqMap)
		if err != nil {
			return errs.NewError(err)
		}
		s.nd.Create(
			daos.GetDBMainCtx(ctx),
			noti,
		)
		if err != nil {
			return errs.NewError(err)
		}
	}
	return nil
}

func (s *NftLend) GetNotifications(ctx context.Context, network models.Network, address string, page int, limit int) ([]*models.Notification, uint, error) {
	notifications, count, err := s.nd.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?": []interface{}{network},
			"address = ?": []interface{}{address},
		},
		map[string][]interface{}{},
		[]string{"id desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return notifications, count, nil
}

func (s *NftLend) SeenNotification(ctx context.Context, req *serializers.SeenNotificationReq) error {
	var user *models.User
	var err error
	if req.Address == "" {
		return errs.NewError(errs.ErrBadRequest)
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
			return errs.NewError(errs.ErrBadRequest)
		}
	}
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err = s.getUser(tx, req.Address, req.Network)
			if err != nil {
				return errs.NewError(err)
			}
			noti, err := s.nd.First(
				tx,
				map[string][]interface{}{
					"network = ?": []interface{}{req.Network},
					"address = ?": []interface{}{req.Address},
				},
				map[string][]interface{}{},
				[]string{"id desc"},
			)
			if noti == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			if noti.ID < req.SeenNotiID {
				return errs.NewError(errs.ErrBadRequest)
			}
			user.SeenNotiID = req.SeenNotiID
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
		return errs.NewError(err)
	}
	return nil
}