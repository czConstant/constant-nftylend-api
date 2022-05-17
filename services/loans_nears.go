package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/serializers"
	"github.com/czConstant/constant-nftylend-api/services/3rd/saletrack"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) NearUpdateLoan(ctx context.Context, req *serializers.CreateLoanNearReq, lastUpdatedClient string) (*models.Loan, bool, error) {
	emailQueue := []*models.EmailQueue{}
	var isUpdated bool
	var loan *models.Loan
	if req.ContractAddress == "" ||
		req.TokenID == "" {
		return nil, false, errs.NewError(errs.ErrBadRequest)
	}
	req.ContractAddress = strings.ToLower(req.ContractAddress)
	asset, err := s.NearSynAsset(ctx, req.ContractAddress, req.TokenID)
	if err != nil {
		return nil, false, errs.NewError(err)
	}
	if asset == nil {
		return nil, false, errs.NewError(errs.ErrBadRequest)
	}
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			saleInfo, err := s.bcs.Near.GetNftpawnSale(s.conf.Contract.NearNftypawnAddress, fmt.Sprintf("%s||%s", asset.ContractAddress, asset.TokenID))
			if err != nil {
				return errs.NewError(err)
			}
			currency, err := s.getLendCurrency(tx, saleInfo.LoanCurrency)
			if err != nil {
				return errs.NewError(err)
			}
			loan, err = s.ld.First(
				tx,
				map[string][]interface{}{
					"network = ?":   []interface{}{models.NetworkNEAR},
					"asset_id = ?":  []interface{}{asset.ID},
					"nonce_hex = ?": []interface{}{saleInfo.CreatedAt},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			principalAmount := models.ConvertWeiToCollateralFloatAmount(&saleInfo.LoanPrincipalAmount.Int, currency.Decimals)
			interestRate, _ := models.ConvertWeiToBigFloat(big.NewInt(int64(saleInfo.LoanInterestRate)), 4).Float64()
			if loan == nil {
				v, err := models.ConvertString2BigInt(saleInfo.CreatedAt)
				if err != nil {
					return errs.NewError(err)
				}
				createdAt := helpers.TimeFromUnix(int64(v.Uint64()))
				loan = &models.Loan{
					Network:         models.NetworkNEAR,
					Owner:           saleInfo.OwnerID,
					PrincipalAmount: numeric.BigFloat{*big.NewFloat(principalAmount)},
					InterestRate:    interestRate,
					Duration:        uint(saleInfo.LoanDuration),
					StartedAt:       createdAt,
					ExpiredAt:       helpers.TimeAdd(*createdAt, time.Duration(saleInfo.LoanDuration)*time.Second),
					ValidAt:         helpers.TimeFromUnix(int64(saleInfo.AvailableAt)),
					Config:          saleInfo.LoanConfig,
					CurrencyID:      currency.ID,
					AssetID:         asset.ID,
					Status:          models.LoanStatusNew,
					NonceHex:        saleInfo.CreatedAt,
					DataLoanAddress: fmt.Sprintf("%d", saleInfo.ApprovalID),
				}
				err = s.ld.Create(
					tx,
					loan,
				)
				if err != nil {
					return errs.NewError(err)
				}
				err = s.ltd.Create(
					tx,
					&models.LoanTransaction{
						Network:         loan.Network,
						Type:            models.LoanTransactionTypeListed,
						LoanID:          loan.ID,
						Borrower:        loan.Owner,
						PrincipalAmount: loan.PrincipalAmount,
						InterestRate:    loan.InterestRate,
						StartedAt:       loan.StartedAt,
						Duration:        loan.Duration,
						ExpiredAt:       loan.ExpiredAt,
					},
				)
				if err != nil {
					return errs.NewError(err)
				}
				isUpdated = true
			}
			loanPrevStatus := loan.Status
			var eqLoan *models.EmailQueue
			var eqOffer *models.EmailQueue
			switch saleInfo.Status {
			case 0:
				{
					loan.Status = models.LoanStatusNew
				}
			case 1:
				{
					loan.Status = models.LoanStatusCreated
					eqLoan = &models.EmailQueue{
						EmailType: models.EMAIL_BORROWER_LOAN_STARTED,
						ObjectID:  loan.ID,
					}
				}
			case 2:
				{
					loan.Status = models.LoanStatusDone
				}
			case 3:
				{
					loan.Status = models.LoanStatusLiquidated
					eqLoan = &models.EmailQueue{
						EmailType: models.EMAIL_BORROWER_LOAN_LIQUIDATED,
						ObjectID:  loan.ID,
					}
				}
			case 4:
				{
					loan.Status = models.LoanStatusDone
				}
			case 5:
				{
					loan.Status = models.LoanStatusCancelled
				}
			default:
				{
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			if loanPrevStatus != loan.Status {
				isUpdated = true
				if eqLoan != nil {
					emailQueue = append(emailQueue, eqLoan)
				}
			}
			for _, saleOffer := range saleInfo.Offers {
				offer, err := s.lod.First(
					tx,
					map[string][]interface{}{
						"network = ?":   []interface{}{models.NetworkNEAR},
						"loan_id = ?":   []interface{}{loan.ID},
						"nonce_hex = ?": []interface{}{fmt.Sprintf("%d", saleOffer.OfferID)},
					},
					map[string][]interface{}{},
					[]string{},
				)
				if err != nil {
					return errs.NewError(err)
				}
				if offer == nil {
					offerPrincipalAmount := models.ConvertWeiToCollateralFloatAmount(&saleOffer.LoanPrincipalAmount.Int, currency.Decimals)
					offerInterestRate, _ := models.ConvertWeiToBigFloat(big.NewInt(int64(saleOffer.LoanInterestRate)), 4).Float64()
					offer = &models.LoanOffer{
						Network:         loan.Network,
						LoanID:          loan.ID,
						Lender:          saleInfo.Lender,
						PrincipalAmount: numeric.BigFloat{*big.NewFloat(offerPrincipalAmount)},
						InterestRate:    offerInterestRate,
						Duration:        uint(saleOffer.LoanDuration),
						Status:          models.LoanOfferStatusNew,
						NonceHex:        fmt.Sprintf("%d", saleOffer.OfferID),
						ValidAt:         helpers.TimeFromUnix(int64(saleOffer.AvailableAt)),
					}
					err = s.lod.Create(
						tx,
						offer,
					)
					if err != nil {
						return errs.NewError(err)
					}
					isUpdated = true
					eqLoan = &models.EmailQueue{
						EmailType: models.EMAIL_BORROWER_OFFER_NEW,
						ObjectID:  offer.ID,
					}
				}
				var isOffered bool
				offerPrevStatus := offer.Status
				switch saleOffer.Status {
				case 0:
					{
						offer.Status = models.LoanOfferStatusNew
					}
				case 1:
					{
						v, err := models.ConvertString2BigInt(saleOffer.StartedAt)
						if err != nil {
							return errs.NewError(err)
						}
						offer.StartedAt = helpers.TimeFromUnix(int64(v.Uint64()))
						offer.Duration = uint(saleOffer.LoanDuration)
						offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(offer.Duration))
						offer.Status = models.LoanOfferStatusApproved
						isOffered = true
						eqLoan = &models.EmailQueue{
							EmailType: models.EMAIL_LENDER_OFFER_STARTED,
							ObjectID:  offer.ID,
						}
					}
				case 2:
					{
						v, err := models.ConvertString2BigInt(saleOffer.StartedAt)
						if err != nil {
							return errs.NewError(err)
						}
						offer.StartedAt = helpers.TimeFromUnix(int64(v.Uint64()))
						offer.Duration = uint(saleOffer.LoanDuration)
						offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(offer.Duration))
						offer.Status = models.LoanOfferStatusDone
						isOffered = true
					}
				case 3:
					{
						v, err := models.ConvertString2BigInt(saleOffer.StartedAt)
						if err != nil {
							return errs.NewError(err)
						}
						offer.StartedAt = helpers.TimeFromUnix(int64(v.Uint64()))
						offer.Duration = uint(saleOffer.LoanDuration)
						offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(offer.Duration))
						offer.Status = models.LoanOfferStatusLiquidated
						isOffered = true
					}
				case 4:
					{
						offer.Status = models.LoanOfferStatusDone
					}
				case 5:
					{
						offer.Status = models.LoanOfferStatusCancelled
					}
				default:
					{
						return errs.NewError(errs.ErrBadRequest)
					}
				}
				if offerPrevStatus != offer.Status {
					isUpdated = true
				}
				err = s.lod.Save(
					tx,
					offer,
				)
				if err != nil {
					return errs.NewError(err)
				}
				if isOffered {
					loan.Lender = offer.Lender
					loan.OfferStartedAt = offer.StartedAt
					loan.OfferDuration = offer.Duration
					loan.OfferExpiredAt = offer.ExpiredAt
					loan.OfferPrincipalAmount = offer.PrincipalAmount
					loan.OfferInterestRate = offer.InterestRate
				}
			}
			if loan.UpdatedAt.After(time.Now().Add(30*time.Second)) &&
				loan.LastUpdatedClient == "worker" {
				isUpdated = true
				if eqOffer != nil {
					emailQueue = append(emailQueue, eqOffer)
				}
			}
			loan.LastUpdatedClient = lastUpdatedClient
			err = s.ld.Save(
				tx,
				loan,
			)
			if err != nil {
				return errs.NewError(err)
			}
			return nil
		},
	)
	if err != nil {
		return nil, false, errs.NewError(err)
	}
	{
		s.EmailForReference(ctx, emailQueue)
	}
	return loan, isUpdated, nil
}

func (s *NftLend) NearCreateLoanOffer(ctx context.Context, loanID uint, req *serializers.CreateLoanOfferReq) (*models.LoanOffer, error) {
	var offer *models.LoanOffer
	if req.PrincipalAmount.Float.Cmp(big.NewFloat(0)) <= 0 ||
		req.Duration <= 0 ||
		req.Lender == "" ||
		req.NonceHex == "" ||
		req.Signature == "" {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	req.Lender = strings.ToLower(req.Lender)
	req.NonceHex = strings.ToLower(req.NonceHex)
	req.Signature = strings.ToLower(req.Signature)
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			loan, err := s.ld.FirstByID(
				tx,
				loanID,
				map[string][]interface{}{},
				false,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if loan == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			if loan.Status != models.LoanStatusNew {
				return errs.NewError(errs.ErrBadRequest)
			}
			currency, err := s.GetCurrencyByID(tx, loan.CurrencyID, loan.Network)
			if err != nil {
				return errs.NewError(err)
			}
			asset, err := s.ad.FirstByID(tx, loan.AssetID, map[string][]interface{}{}, false)
			if err != nil {
				return errs.NewError(err)
			}
			if asset == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			msgHex := helpers.AppendHexStrings(
				helpers.ParseBigInt2Hex(models.Number2BigInt(req.PrincipalAmount.String(), int(currency.Decimals))),
				helpers.ParseBigInt2Hex(models.Number2BigInt(asset.TokenID, 0)),
				helpers.ParseBigInt2Hex(big.NewInt(int64(req.Duration))),
				helpers.ParseBigInt2Hex(models.Number2BigInt(fmt.Sprintf("%f", req.InterestRate), 4)),
				helpers.ParseBigInt2Hex(big.NewInt(s.getEvmAdminFee(loan.Network))),
				helpers.ParseHex2Hex(req.NonceHex),
				helpers.ParseAddress2Hex(asset.ContractAddress),
				helpers.ParseAddress2Hex(currency.ContractAddress),
				helpers.ParseAddress2Hex(req.Lender),
				helpers.ParseBigInt2Hex(big.NewInt(s.getEvmClientByNetwork(loan.Network).ChainID)),
			)
			err = s.getEvmClientByNetwork(loan.Network).ValidateSignature(
				msgHex,
				req.Signature,
				req.Lender,
			)
			if err != nil {
				return errs.NewError(err)
			}
			offer, err = s.lod.First(
				tx,
				map[string][]interface{}{
					"lender = ?":    []interface{}{req.Lender},
					"nonce_hex = ?": []interface{}{req.NonceHex},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if offer != nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			offer = &models.LoanOffer{
				Network:         loan.Network,
				LoanID:          loan.ID,
				Lender:          req.Lender,
				PrincipalAmount: req.PrincipalAmount,
				InterestRate:    req.InterestRate,
				Duration:        req.Duration,
				Status:          models.LoanOfferStatusNew,
				NonceHex:        req.NonceHex,
				Signature:       req.Signature,
			}
			err = s.lod.Create(
				tx,
				offer,
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
	return offer, nil
}

func (s *NftLend) NearSynAsset(ctx context.Context, contractAddress string, tokenID string) (*models.Asset, error) {
	var asset *models.Asset
	var err error
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			asset, err = s.ad.First(
				tx,
				map[string][]interface{}{
					"network = ?":          []interface{}{models.NetworkNEAR},
					"contract_address = ?": []interface{}{contractAddress},
					"token_id = ?":         []interface{}{tokenID},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if asset == nil {
				var collectionName, description, assetName string
				metaData, err := s.bcs.Near.GetNftMetadata(contractAddress)
				if err != nil {
					return errs.NewError(err)
				}
				if metaData == nil {
					return errs.NewError(errs.ErrBadRequest)
				}
				collectionName = metaData.Name
				tokenData, err := s.bcs.Near.GetNftToken(contractAddress, tokenID)
				if err != nil {
					return errs.NewError(err)
				}
				if tokenData == nil {
					return errs.NewError(errs.ErrBadRequest)
				}
				if tokenData.Metadata.Description != "" {
					description = tokenData.Metadata.Description
				}
				if assetName == "" {
					assetName = tokenData.Metadata.Title
				}
				var tokenURL string
				mediaURL := helpers.MergeMetaInfoURL(metaData.BaseUri, tokenData.Metadata.Media)
				var metaInfo *saletrack.EvmNftMetaResp
				if tokenData.Metadata.Reference != "" {
					tokenURL = helpers.MergeMetaInfoURL(metaData.BaseUri, tokenData.Metadata.Reference)
					metaInfo, err = s.stc.GetEvmNftMetaResp(helpers.ConvertImageDataURL(tokenURL))
					if err != nil {
						return errs.NewError(err)
					}
					if description == "" {
						description = metaInfo.Description
					}
					if assetName == "" {
						assetName = metaInfo.Name
					}
				}
				if collectionName == "" {
					return errs.NewError(errs.ErrBadRequest)
				}
				collection, err := s.cld.First(
					tx,
					map[string][]interface{}{
						"network = ?":          []interface{}{models.NetworkNEAR},
						"contract_address = ?": []interface{}{contractAddress},
					},
					map[string][]interface{}{},
					[]string{},
				)
				if err != nil {
					return errs.NewError(err)
				}
				if collection == nil {
					var isVerified bool
					// parasProfiles, err := s.stc.GetParasProfile(asset.GetContractAddress())
					// if err != nil {
					// 	return errs.NewError(err)
					// }
					// if len(parasProfiles) > 0 {
					// 	isVerified = parasProfiles[0].IsCreator
					// }
					collection = &models.Collection{
						Network:         models.NetworkNEAR,
						SeoURL:          helpers.MakeSeoURL(fmt.Sprintf("%s-%s", models.NetworkNEAR, contractAddress)),
						ContractAddress: contractAddress,
						Name:            collectionName,
						Description:     description,
						Enabled:         true,
						Verified:        isVerified,
					}
					err = s.cld.Create(
						tx,
						collection,
					)
					if err != nil {
						return errs.NewError(err)
					}
				}
				asset = &models.Asset{
					Network:               models.NetworkNEAR,
					CollectionID:          collection.ID,
					SeoURL:                helpers.MakeSeoURL(fmt.Sprintf("%s-%s", models.NetworkNEAR, fmt.Sprintf("%s-%s", contractAddress, tokenID))),
					ContractAddress:       collection.ContractAddress,
					TokenID:               tokenID,
					Symbol:                "",
					Name:                  assetName,
					TokenURL:              mediaURL,
					ExternalUrl:           tokenURL,
					SellerFeeRate:         0,
					MetaJsonUrl:           tokenURL,
					OriginNetwork:         "",
					OriginContractAddress: "",
					OriginTokenID:         "",
				}
				if metaInfo != nil {
					attributes, err := json.Marshal(metaInfo.Attributes)
					if err != nil {
						return errs.NewError(err)
					}
					asset.Attributes = string(attributes)
					metaJson, err := json.Marshal(metaInfo)
					if err != nil {
						return errs.NewError(err)
					}
					asset.MetaJson = string(metaJson)
				}
				err = s.ad.Create(
					tx,
					asset,
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
	return asset, nil
}
