package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/czConstant/blockchain-api/bcclient/solana"
	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) CreateCollectionSubmission(ctx context.Context, req *serializers.CollectionSubmissionReq) error {
	err := s.clsd.Create(
		daos.GetDBMainCtx(ctx),
		&models.CollectionSubmission{
			Network:         req.Network,
			Name:            req.Name,
			Description:     req.Description,
			Creator:         req.Creator,
			ContractAddress: req.ContractAddress,
			ContactInfo:     req.ContactInfo,
			Verified:        req.Verified,
			WhoVerified:     req.WhoVerified,
			Status:          models.CollectionSubmissionStatusSubmitted,
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
			"status in (?)": []interface{}{[]models.CollectionSubmissionStatus{
				models.CollectionSubmissionStatusApproved,
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

func (s *NftLend) GetNearApprovedCollections(ctx context.Context) ([]*models.CollectionSubmission, error) {
	ms, err := s.clsd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network = ?": []interface{}{models.NetworkNEAR},
			"status in (?)": []interface{}{[]models.CollectionSubmissionStatus{
				models.CollectionSubmissionStatusApproved,
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
				select 1 from collection_submissions
				where collection_submissions.network = collections.network
					and collection_submissions.creator = collections.creator
					and collection_submissions.status = ?
			)`: []interface{}{
				models.CollectionSubmissionStatusApproved,
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
					saleCurrency, err := s.getCurrencyByNetworkSymbol(tx, models.NetworkNEAR, "NEAR")
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

func (s *NftLend) GetRPTListingCollection(ctx context.Context) ([]*models.NftyRPTListingCollection, error) {
	ms, err := s.ad.GetRPTListingCollection(
		daos.GetDBMainCtx(ctx),
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return ms, nil
}

func (s *NftLend) GetCollectionVerified(ctx context.Context, network models.Network, contractAddress string, tokenID string) (*models.Collection, error) {
	var m *models.Collection
	var err error
	switch network {
	case models.NetworkSOL:
		{
			m, _, err = s.getSolanaCollectionVerified(
				daos.GetDBMainCtx(ctx),
				contractAddress,
				nil,
				nil,
			)
			if err != nil {
				return nil, errs.NewError(err)
			}
		}
	case models.NetworkNEAR:
		{
			asset, err := s.CreateNearAsset(ctx, contractAddress, tokenID)
			if err != nil {
				return nil, errs.NewError(err)
			}
			m, err = s.cld.FirstByID(
				daos.GetDBMainCtx(ctx),
				asset.CollectionID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return nil, errs.NewError(err)
			}
			csM, err := s.clsd.First(
				daos.GetDBMainCtx(ctx),
				map[string][]interface{}{
					"network = ?": []interface{}{m.Network},
					"creator = ?": []interface{}{m.Creator},
					"status = ?":  []interface{}{models.CollectionSubmissionStatusApproved},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return nil, errs.NewError(err)
			}
			if csM == nil {
				m = nil
			}
		}
	}
	return m, nil
}

func (s *NftLend) getSolanaCollectionVerified(tx *gorm.DB, mintAddress string, meta *solana.MetadataResp, metaInfo *solana.MetadataInfoResp) (*models.Collection, string, error) {
	vrs, err := s.bcs.SolanaNftVerifier.GetNftVerifier(mintAddress)
	if err != nil {
		return nil, "", errs.NewError(err)
	}
	if vrs.IsWrapped {
		chain := s.bcs.SolanaNftVerifier.ParseChain(vrs.ChainID)
		m, err := s.cld.First(
			tx,
			map[string][]interface{}{
				"origin_network = ?":          []interface{}{chain},
				"origin_contract_address = ?": []interface{}{vrs.AssetAddress},
				"enabled = ?":                 []interface{}{true},
			},
			map[string][]interface{}{},
			[]string{"id desc"},
		)
		if err != nil {
			return nil, "", errs.NewError(err)
		}
		if m != nil {
			return m, vrs.TokenID, nil
		}
	} else {
		if meta == nil {
			meta, err = s.bcs.Solana.GetMetadata(mintAddress)
			if err != nil {
				return nil, "", errs.NewError(err)
			}
		}
		if metaInfo == nil {
			metaInfo, err = s.bcs.Solana.GetMetadataInfo(meta.Data.Uri)
			if err != nil {
				return nil, "", errs.NewError(err)
			}
		}
		collectionName := metaInfo.Collection.Name
		if collectionName == "" {
			collectionName = metaInfo.Collection.Family
			if collectionName == "" {
				names := strings.Split(metaInfo.Name, "#")
				if len(names) >= 2 {
					collectionName = strings.TrimSpace(names[0])
				}
			}
		}
		if collectionName == "" {
			return nil, "", errs.NewError(err)
		}
		for _, creator := range meta.Data.Creators {
			m, err := s.cld.First(
				tx,
				map[string][]interface{}{
					"name = ?":       []interface{}{collectionName},
					"creator like ?": []interface{}{fmt.Sprintf("%%%s%%", creator.Address)},
					"enabled = ?":    []interface{}{true},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return nil, "", errs.NewError(err)
			}
			if m != nil {
				return m, "", nil
			}
		}
	}
	return nil, "", nil
}

func (s *NftLend) GetCollections(ctx context.Context, page int, limit int) ([]*models.Collection, uint, error) {
	categories, count, err := s.cld.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network in (?)":  []interface{}{s.getSupportedNetworks()},
			"new_loan_id > ?": []interface{}{0},
			`approved = true or exists(
				select 1
				from collection_submissions
				where collections.network = collection_submissions.network
				  and collections.creator = collection_submissions.creator
				  and collection_submissions.status = ?
			)`: []interface{}{models.CollectionSubmissionStatusApproved},
		},
		map[string][]interface{}{
			"Currency":      []interface{}{},
			"NewLoan":       []interface{}{},
			"NewLoan.Asset": []interface{}{},
		},
		[]string{"id desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return categories, count, nil
}

func (s *NftLend) GetCollectionDetail(ctx context.Context, seoURL string) (*models.Collection, error) {
	m, err := s.cld.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"seo_url = ?":    []interface{}{seoURL},
			"network in (?)": []interface{}{s.getSupportedNetworks()},
		},
		map[string][]interface{}{
			"Currency":      []interface{}{},
			"NewLoan":       []interface{}{},
			"NewLoan.Asset": []interface{}{},
		},
		[]string{"id desc"},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) updateCollectionForLoan(tx *gorm.DB, collectionID uint) error {
	collection, err := s.cld.FirstByID(
		tx,
		collectionID,
		map[string][]interface{}{},
		true,
	)
	if err != nil {
		return errs.NewError(err)
	}
	loan, err := s.ld.First(
		tx,
		map[string][]interface{}{
			"collection_id = ?": []interface{}{collection.ID},
			"status = ?":        []interface{}{models.LoanStatusNew},
		},
		map[string][]interface{}{},
		[]string{"id desc"},
	)
	if err != nil {
		return errs.NewError(err)
	}
	collection.NewLoanID = 0
	if loan != nil {
		collection.NewLoanID = loan.ID
	}
	err = s.cld.Save(
		tx,
		collection,
	)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}
