package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/jinzhu/gorm"
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

func (s *NftLend) JobUpdateProfileCollection(ctx context.Context) error {
	var retErr error
	collections, err := s.cld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?":   []interface{}{models.NetworkNEAR},
			"creator != ''": []interface{}{},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return errs.NewError(err)
	}
	for _, collection := range collections {
		err = s.UpdateProfileCollection(ctx, collection.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	return retErr
}

func (s *NftLend) UpdateProfileCollection(ctx context.Context, collectionID uint) error {
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
	if collection.Creator != "" {
		parasProfiles, err := s.stc.GetParasProfile(collection.Creator)
		if err != nil {
			return errs.NewError(err)
		}
		if len(parasProfiles) > 0 {
			err = daos.WithTransaction(
				daos.GetDBMainCtx(ctx),
				func(tx *gorm.DB) error {
					collection, err = s.cld.FirstByID(
						tx,
						collection.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					collection.Verified = parasProfiles[0].IsCreator
					collection.CoverURL = parasProfiles[0].CoverURL
					collection.ImageURL = parasProfiles[0].ImgURL
					collection.CreatorURL = parasProfiles[0].Website
					collection.TwitterID = parasProfiles[0].TwitterId
					err = s.cld.Save(
						tx,
						collection,
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
		}
	}
	return nil
}
