package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) CreateNotification(ctx context.Context, network models.Network, address string, notiType models.NotificationType, reqMap map[string]interface{}) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			user, err := s.getUser(tx, network, address, true)
			if err != nil {
				return errs.NewError(err)
			}
			notiTpl, err := s.ntd.First(
				tx,
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
				noti, err := notiTpl.Execute(network, user.ID, reqMap)
				if err != nil {
					return errs.NewError(err)
				}
				s.nd.Create(
					tx,
					noti,
				)
				if err != nil {
					return errs.NewError(err)
				}
				newNotiNum, err := s.nd.Count(
					tx,
					map[string][]interface{}{
						"user_id = ?": []interface{}{user.ID},
						"id > ?":      []interface{}{user.SeenNotiID},
					},
				)
				if err != nil {
					return errs.NewError(err)
				}
				user.NewNotiID = noti.ID
				user.NewNotiNum = newNotiNum
				s.ud.Save(
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
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) GetNotifications(ctx context.Context, network models.Network, address string, page int, limit int) ([]*models.Notification, uint, error) {
	user, err := s.GetUser(ctx, network, address)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	notifications, count, err := s.nd.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"user_id = ?": []interface{}{user.ID},
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
	if req.SeenNotiID <= 0 {
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
			user, err = s.getUser(tx, req.Network, req.Address, true)
			if err != nil {
				return errs.NewError(err)
			}
			if user.NewNotiID < req.SeenNotiID {
				return errs.NewError(errs.ErrBadRequest)
			}
			m, err := s.nd.FirstByID(
				tx,
				req.SeenNotiID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if m == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			if m.UserID != user.ID {
				return errs.NewError(errs.ErrBadRequest)
			}
			newNotiNum, err := s.nd.Count(
				tx,
				map[string][]interface{}{
					"user_id = ?": []interface{}{user.ID},
					"id > ?":      []interface{}{m.ID},
				},
			)
			if err != nil {
				return errs.NewError(err)
			}
			user.SeenNotiID = m.ID
			user.NewNotiNum = newNotiNum
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
