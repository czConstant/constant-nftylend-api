package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
)

func (s *NftLend) CreateCollectionSubmitted(ctx context.Context, req *serializers.CollectionSubmittedReq) error {
	err := s.clsd.Create(
		daos.GetDBMainCtx(ctx),
		&models.CollectionSubmitted{
			Network:         req.Network,
			Name:            req.Name,
			Description:     req.Description,
			Creator:         req.Creator,
			ContractAddress: req.ContractAddress,
			ContactInfo:     req.ContactInfo,
			Verified:        req.Verified,
			WhoVerified:     req.WhoVerified,
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) JobVolumeCollections(ctx context.Context) error {
	var retErr error
	return retErr
}

func (s *NftLend) UpdateVolumeCollection(ctx context.Context, collectionID uint) error {
	collection, err := s.cld.FirstByID(
		daos.GetDBMainCtx(ctx),
		collectionID,
		map[string][]interface{}{},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	if collection == nil {
		return errs.NewError(errs.ErrBadRequest)
	}
	if collection.ParasCollectionID != "" {
		
	}
	return nil
}
