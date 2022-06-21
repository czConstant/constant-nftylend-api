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
	asset, err := s.CreateNearAsset(ctx, req.ContractAddress, req.TokenID)
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
			principalAmount := models.ConvertWeiToBigFloat(&saleInfo.LoanPrincipalAmount.Int, currency.Decimals)
			interestRate, _ := models.ConvertWeiToBigFloat(big.NewInt(int64(saleInfo.LoanInterestRate)), 4).Float64()
			if loan == nil {
				v, err := models.ConvertString2BigInt(saleInfo.CreatedAt)
				if err != nil {
					return errs.NewError(err)
				}
				createdAt := helpers.TimeFromUnix(int64(v.Uint64()))
				borrower, err := s.getUser(
					tx,
					models.NetworkNEAR,
					saleInfo.OwnerID,
					false,
				)
				if err != nil {
					return errs.NewError(err)
				}
				loan = &models.Loan{
					Network:         models.NetworkNEAR,
					Owner:           saleInfo.OwnerID,
					BorrowerUserID:  borrower.ID,
					PrincipalAmount: numeric.BigFloat{*principalAmount},
					InterestRate:    interestRate,
					Duration:        uint(saleInfo.LoanDuration),
					StartedAt:       createdAt,
					ExpiredAt:       helpers.TimeAdd(*createdAt, time.Duration(saleInfo.LoanDuration)*time.Second),
					ValidAt:         helpers.TimeFromUnix(int64(saleInfo.AvailableAt)),
					Config:          saleInfo.LoanConfig,
					CurrencyID:      currency.ID,
					AssetID:         asset.ID,
					CollectionID:    asset.CollectionID,
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
			// check existsed loan
			{
				loanCheck, err := s.ld.First(
					tx,
					map[string][]interface{}{
						"id != ?":      []interface{}{loan.ID},
						"network = ?":  []interface{}{loan.Network},
						"asset_id = ?": []interface{}{loan.AssetID},
						"status = ?":   []interface{}{models.LoanStatusCreated},
					},
					map[string][]interface{}{},
					[]string{},
				)
				if err != nil {
					return errs.NewError(err)
				}
				if loanCheck != nil {
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			// cancel old pending loan
			{
				loans, err := s.ld.Find(
					tx,
					map[string][]interface{}{
						"id != ?":      []interface{}{loan.ID},
						"network = ?":  []interface{}{loan.Network},
						"asset_id = ?": []interface{}{loan.AssetID},
						"status = ?":   []interface{}{models.LoanStatusNew},
					},
					map[string][]interface{}{},
					[]string{},
					0,
					999999,
				)
				if err != nil {
					return errs.NewError(err)
				}
				for _, l := range loans {
					l, err = s.ld.FirstByID(
						tx,
						l.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if l.Status != models.LoanStatusNew {
						return errs.NewError(errs.ErrBadRequest)
					}
					l.Status = models.LoanStatusCancelled
					err = s.ld.Save(
						tx,
						l,
					)
					if err != nil {
						return errs.NewError(err)
					}
					err = s.updateIncentiveForLoan(
						tx,
						loan,
					)
					if err != nil {
						return errs.NewError(err)
					}
				}
			}
			loan, err = s.ld.FirstByID(
				tx,
				loan.ID,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			if loan.SynchronizedAt != nil &&
				loan.SynchronizedAt.After(time.Now().Add(-15*time.Second)) {
				isUpdated = true
				return nil
			}
			loanPrevStatus := loan.Status
			switch saleInfo.Status {
			case 0:
				{
					loan.Status = models.LoanStatusNew
				}
			case 1:
				{
					loan.Status = models.LoanStatusCreated
					if loanPrevStatus != loan.Status {
						emailQueue = append(emailQueue, &models.EmailQueue{
							EmailType: models.EMAIL_BORROWER_LOAN_STARTED,
							ObjectID:  loan.ID,
						})
					}
				}
			case 2:
				{
					loan.Status = models.LoanStatusDone
					if loanPrevStatus != loan.Status {
						emailQueue = append(emailQueue, &models.EmailQueue{
							EmailType: models.EMAIL_LENDER_LOAN_REPAID,
							ObjectID:  loan.ID,
						})
					}
				}
			case 3:
				{
					loan.Status = models.LoanStatusLiquidated
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
					offerPrincipalAmount := models.ConvertWeiToBigFloat(&saleOffer.LoanPrincipalAmount.Int, currency.Decimals)
					offerInterestRate, _ := models.ConvertWeiToBigFloat(big.NewInt(int64(saleOffer.LoanInterestRate)), 4).Float64()
					lender, err := s.getUser(
						tx,
						models.NetworkNEAR,
						saleInfo.Lender,
						false,
					)
					if err != nil {
						return errs.NewError(err)
					}
					offer = &models.LoanOffer{
						Network:         loan.Network,
						LoanID:          loan.ID,
						Lender:          saleInfo.Lender,
						LenderUserID:    lender.ID,
						PrincipalAmount: numeric.BigFloat{*offerPrincipalAmount},
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
					emailQueue = append(emailQueue, &models.EmailQueue{
						EmailType: models.EMAIL_BORROWER_OFFER_NEW,
						ObjectID:  offer.ID,
					})
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
						if offerPrevStatus != offer.Status {
							emailQueue = append(emailQueue, &models.EmailQueue{
								EmailType: models.EMAIL_LENDER_OFFER_STARTED,
								ObjectID:  offer.ID,
							})
							err = s.ltd.Create(
								tx,
								&models.LoanTransaction{
									Network:         loan.Network,
									Type:            models.LoanTransactionTypeOffered,
									LoanID:          loan.ID,
									Borrower:        loan.Owner,
									Lender:          offer.Lender,
									PrincipalAmount: offer.PrincipalAmount,
									InterestRate:    offer.InterestRate,
									StartedAt:       offer.StartedAt,
									Duration:        offer.Duration,
									ExpiredAt:       offer.ExpiredAt,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
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
						if offerPrevStatus != offer.Status {
							err = s.ltd.Create(
								tx,
								&models.LoanTransaction{
									Network:         loan.Network,
									Type:            models.LoanTransactionTypeRepaid,
									LoanID:          loan.ID,
									Borrower:        loan.Owner,
									Lender:          offer.Lender,
									PrincipalAmount: offer.PrincipalAmount,
									InterestRate:    offer.InterestRate,
									StartedAt:       offer.StartedAt,
									Duration:        offer.Duration,
									ExpiredAt:       offer.ExpiredAt,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
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
						if offerPrevStatus != offer.Status {
							err = s.ltd.Create(
								tx,
								&models.LoanTransaction{
									Network:         loan.Network,
									Type:            models.LoanTransactionTypeLiquidated,
									LoanID:          loan.ID,
									Borrower:        loan.Owner,
									Lender:          offer.Lender,
									PrincipalAmount: offer.PrincipalAmount,
									InterestRate:    offer.InterestRate,
									StartedAt:       offer.StartedAt,
									Duration:        offer.Duration,
									ExpiredAt:       offer.ExpiredAt,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					}
				case 4:
					{
						offer.Status = models.LoanOfferStatusDone
					}
				case 5:
					{
						offer.Status = models.LoanOfferStatusCancelled
						err = s.ltd.Create(
							tx,
							&models.LoanTransaction{
								Network:         loan.Network,
								Type:            models.LoanTransactionTypeCancelled,
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
					loan.LenderUserID = offer.LenderUserID
					loan.OfferStartedAt = offer.StartedAt
					loan.OfferDuration = offer.Duration
					loan.OfferExpiredAt = offer.ExpiredAt
					loan.OfferPrincipalAmount = offer.PrincipalAmount
					loan.OfferInterestRate = offer.InterestRate
					if loan.CurrencyPrice <= 0 {
						loan.CurrencyPrice = currency.Price
					}
				}
			}
			if isUpdated {
				loan.SynchronizedAt = helpers.TimeNow()
			}
			err = s.ld.Save(
				tx,
				loan,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.updateIncentiveForLoan(
				tx,
				loan,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.updateCollectionForLoan(
				tx,
				loan.CollectionID,
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

func (s *NftLend) UpdateIncentiveForLoanID(ctx context.Context, loanID uint) error {
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
			err = s.updateIncentiveForLoan(
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
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) updateIncentiveForLoan(tx *gorm.DB, loan *models.Loan) error {
	var err error
	switch loan.Status {
	case models.LoanStatusNew:
		{
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeBorrowerLoanListed,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
		}
	case models.LoanStatusDone:
		{
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeBorrowerLoanListed,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeLenderLoanMatched,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeAffiliateBorrowerLoanDone,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeAffiliateLenderLoanDone,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
		}
	case models.LoanStatusCreated,
		models.LoanStatusLiquidated,
		models.LoanStatusExpired:
		{
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeBorrowerLoanListed,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeLenderLoanMatched,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
		}
	case models.LoanStatusCancelled:
		{
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeBorrowerLoanListed,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
			err = s.IncentiveForLoan(
				tx,
				models.IncentiveTransactionTypeBorrowerLoanDelisted,
				loan.ID,
			)
			if err != nil {
				return errs.NewError(err)
			}
		}
	}
	return nil
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

func (s *NftLend) CreateNearAsset(ctx context.Context, contractAddress string, tokenID string) (*models.Asset, error) {
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
				var collectionName, description, assetDescription, assetName string
				collectionData, err := s.bcs.Near.GetNftMetadata(contractAddress)
				if err != nil {
					return errs.NewError(err)
				}
				if collectionData == nil {
					return errs.NewError(errs.ErrBadRequest)
				}
				collectionName = collectionData.Name
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
				var tokenURL, mimeType string
				parasCollectionID := contractAddress
				mediaURL := helpers.MergeMetaInfoURL(collectionData.BaseUri, tokenData.Metadata.Media)
				var tokenMetaData *saletrack.EvmNftMetaResp
				var sellerFeeRate float64
				seoURL := helpers.MakeSeoURL(fmt.Sprintf("%s-%s", models.NetworkNEAR, contractAddress))
				creator := contractAddress
				if tokenData.Metadata.Reference != "" {
					tokenURL = helpers.MergeMetaInfoURL(collectionData.BaseUri, tokenData.Metadata.Reference)
					tokenMetaData, err = s.stc.GetEvmNftMetaResp(helpers.ConvertImageDataURL(tokenURL))
					if err != nil {
						return errs.NewError(err, tokenURL)
					}
					if description == "" {
						description = tokenMetaData.Description
						assetDescription = tokenMetaData.Description
					}
					if assetName == "" {
						assetName = tokenMetaData.Name
					}
					mimeType = tokenMetaData.MimeType
					switch contractAddress {
					case "x.paras.near":
						{
							seriesData, err := s.bcs.Near.GetNftSeries(contractAddress, strings.Split(tokenID, ":")[0])
							if err != nil {
								return errs.NewError(err)
							}
							if seriesData == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							creator = seriesData.CreatorID
							seriesURL := helpers.MergeMetaInfoURL(collectionData.BaseUri, seriesData.Metadata.Reference)
							seriesMetaData, err := s.stc.GetEvmNftMetaResp(helpers.ConvertImageDataURL(seriesURL))
							if err != nil {
								return errs.NewError(err, seriesURL)
							}
							parasCollectionID = seriesMetaData.CollectionID
							if parasCollectionID == "" {
								return errs.NewError(errs.ErrBadRequest)
							}
							collectionName = seriesMetaData.CollectionName
							description = seriesMetaData.Description
							sellerFeeRate, _ = models.ConvertWeiToBigFloat(big.NewInt(seriesData.RoyaltyRate), 4).Float64()
							seoURL = helpers.MakeSeoURL(fmt.Sprintf("%s-%s-%s", models.NetworkNEAR, contractAddress, parasCollectionID))
						}
					}
				}
				if collectionName == "" {
					return errs.NewError(errs.ErrBadRequest)
				}
				collection, err := s.cld.First(
					tx,
					map[string][]interface{}{
						"seo_url = ?": []interface{}{seoURL},
					},
					map[string][]interface{}{},
					[]string{},
				)
				if err != nil {
					return errs.NewError(err)
				}
				if collection == nil {
					var isVerified bool
					parasProfiles, err := s.stc.GetParasProfile(creator)
					if err != nil {
						return errs.NewError(err)
					}
					collection = &models.Collection{
						Network:           models.NetworkNEAR,
						SeoURL:            seoURL,
						ContractAddress:   contractAddress,
						Name:              collectionName,
						Description:       description,
						Enabled:           true,
						Verified:          isVerified,
						ParasCollectionID: parasCollectionID,
						Creator:           creator,
					}
					if len(parasProfiles) > 0 {
						collection.Verified = parasProfiles[0].IsCreator
						collection.CoverURL = parasProfiles[0].CoverURL
						collection.ImageURL = parasProfiles[0].ImgURL
						collection.CreatorURL = parasProfiles[0].Website
						collection.TwitterID = parasProfiles[0].TwitterId
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
					SellerFeeRate:         sellerFeeRate,
					MetaJsonUrl:           tokenURL,
					OriginNetwork:         "",
					OriginContractAddress: "",
					OriginTokenID:         "",
					Description:           assetDescription,
					MimeType:              mimeType,
				}
				if tokenMetaData != nil {
					attributes, err := json.Marshal(tokenMetaData.Attributes)
					if err != nil {
						return errs.NewError(err)
					}
					asset.Attributes = string(attributes)
					metaJson, err := json.Marshal(tokenMetaData)
					if err != nil {
						return errs.NewError(err)
					}
					asset.MetaJson = string(metaJson)
				}
				asset.SearchText = strings.TrimSpace(strings.ToLower(fmt.Sprintf("%s %s %s %s %s", collection.Name, collection.Description, asset.ContractAddress, asset.Name, asset.Description)))
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
