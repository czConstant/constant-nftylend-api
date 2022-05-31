package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
)

func (s *NftLend) CreateNotification(ctx context.Context, notiType models.NotificationType, address string, reqMap map[string]interface{}) error {
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
		noti, err := notiTpl.Execute(address, reqMap)
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
