package services

import (
	"context"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
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

func (s *NftLend) GetNearApprovedCreators(ctx context.Context) ([]string, error) {
	ms, err := s.clsd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?": []interface{}{models.NetworkNEAR},
			"status in (?)": []interface{}{[]models.CollectionSubmittedStatus{
				models.CollectionSubmittedStatusApproved,
			}},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	rets := []string{}
	for _, m := range ms {
		rets = append(rets, m.Creator)
	}
	return rets, nil
}

func (s *NftLend) GetNearApprovedCollections(ctx context.Context) ([]*models.CollectionSubmitted, error) {
	ms, err := s.clsd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?": []interface{}{models.NetworkNEAR},
			"status in (?)": []interface{}{[]models.CollectionSubmittedStatus{
				models.CollectionSubmittedStatusApproved,
			}},
		},
		map[string][]interface{}{},
		[]string{},
		0,
		999999,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return ms, nil
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

func (s *NftLend) JobUpdateStatsCollection(ctx context.Context) error {
	var retErr error
	collections, err := s.cld.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?":               []interface{}{models.NetworkNEAR},
			"paras_collection_id != ''": []interface{}{},
			`exists (
				select 1 from collection_submitteds
				where collection_submitteds.network = collections.network
					and collection_submitteds.creator = collections.creator
					and collection_submitteds.status = ?
			)`: []interface{}{
				models.CollectionSubmittedStatusApproved,
			},
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
		err = s.UpdateStatsCollection(ctx, collection.ID)
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	return retErr
}

func (s *NftLend) UpdateStatsCollection(ctx context.Context, collectionID uint) error {
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
		parasStats, err := s.stc.GetParasCollectionStats(collection.ParasCollectionID)
		if err != nil {
			return errs.NewError(err)
		}
		if parasStats.ID != "" {
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
					collection.VolumeUsd = parasStats.VolumeUsd
					saleCurrency, err := s.getLendCurrencyBySymbol(tx, models.NetworkNEAR, "NEAR")
					if err != nil {
						return errs.NewError(err)
					}
					floorPrice := models.ConvertWeiToBigFloat(&parasStats.FloorPrice.Int, saleCurrency.Decimals)
					collection.FloorPrice = numeric.BigFloat{*floorPrice}
					collection.CurrencyID = saleCurrency.ID
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
