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
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) NearUpdateLoan(ctx context.Context, req *serializers.CreateLoanNearReq, lastUpdatedClient string) (*models.Loan, bool, error) {
	var isUpdated bool
	var loan *models.Loan
	if req.ContractAddress == "" ||
		req.TokenID == "" {
		return nil, false, errs.NewError(errs.ErrBadRequest)
	}
	req.ContractAddress = strings.ToLower(req.ContractAddress)
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			saleInfo, err := s.bcs.Near.GetNftpawnSale(s.conf.Contract.NearNftypawnAddress, fmt.Sprintf("%s||%s", req.ContractAddress, req.TokenID))
			if err != nil {
				return errs.NewError(err)
			}
			currency, err := s.getLendCurrency(tx, saleInfo.LoanCurrency)
			if err != nil {
				return errs.NewError(err)
			}
			asset, err := s.ad.First(
				tx,
				map[string][]interface{}{
					"network = ?":          []interface{}{models.NetworkNEAR},
					"contract_address = ?": []interface{}{req.ContractAddress},
					"token_id = ?":         []interface{}{req.TokenID},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if asset == nil {
				metaData, err := s.bcs.Near.GetNftMetadata(saleInfo.NftContractID, saleInfo.TokenID)
				if err != nil {
					return errs.NewError(err)
				}
				if metaData.Metadata.Reference == "" {
					return errs.NewError(errs.ErrBadRequest)
				}
				metaInfo, err := s.stc.GetNearNftMetaResp(helpers.ConvertImageDataURL(metaData.Metadata.Reference))
				if err != nil {
					return errs.NewError(err)
				}
				collection, err := s.cld.First(
					tx,
					map[string][]interface{}{
						"network = ?":          []interface{}{models.NetworkNEAR},
						"contract_address = ?": []interface{}{req.ContractAddress},
					},
					map[string][]interface{}{},
					[]string{},
				)
				if err != nil {
					return errs.NewError(err)
				}
				if collection == nil {
					collection = &models.Collection{
						Network:         models.NetworkNEAR,
						SeoURL:          helpers.MakeSeoURL(req.ContractAddress),
						ContractAddress: req.ContractAddress,
						Name:            metaInfo.Collection,
						Description:     metaInfo.Description,
						Enabled:         true,
					}
					err = s.cld.Create(
						tx,
						collection,
					)
					if err != nil {
						return errs.NewError(err)
					}
				}
				attributes, err := json.Marshal(metaInfo.Attributes)
				if err != nil {
					return errs.NewError(err)
				}
				metaJson, err := json.Marshal(metaInfo)
				if err != nil {
					return errs.NewError(err)
				}
				asset = &models.Asset{
					Network:               models.NetworkNEAR,
					CollectionID:          collection.ID,
					SeoURL:                helpers.MakeSeoURL(fmt.Sprintf("%s-%s", req.ContractAddress, req.TokenID)),
					ContractAddress:       collection.ContractAddress,
					TokenID:               req.TokenID,
					Symbol:                "",
					Name:                  metaData.Metadata.Title,
					TokenURL:              metaData.Metadata.Media,
					ExternalUrl:           metaData.Metadata.Reference,
					SellerFeeRate:         0,
					Attributes:            string(attributes),
					MetaJson:              string(metaJson),
					MetaJsonUrl:           "",
					OriginNetwork:         "",
					OriginContractAddress: "",
					OriginTokenID:         "",
				}
				err = s.ad.Create(
					tx,
					asset,
				)
				if err != nil {
					return errs.NewError(err)
				}
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
				loan = &models.Loan{
					Network:         models.NetworkNEAR,
					Owner:           saleInfo.OwnerID,
					PrincipalAmount: numeric.BigFloat{*big.NewFloat(principalAmount)},
					InterestRate:    interestRate,
					Duration:        uint(saleInfo.LoanDuration),
					StartedAt:       helpers.TimeNow(),
					ExpiredAt:       helpers.TimeAdd(time.Now(), time.Duration(saleInfo.LoanDuration)*time.Second),
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
			switch saleInfo.Status {
			case 0:
				{
					loan.Status = models.LoanStatusNew
				}
			case 1:
				{
					loan.Status = models.LoanStatusCreated
				}
			case 2:
				{
					loan.Status = models.LoanStatusDone
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
					}
					err = s.lod.Create(
						tx,
						offer,
					)
					if err != nil {
						return errs.NewError(err)
					}
					isUpdated = true
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
						startedAt := time.Unix(int64(v.Uint64())/1000, 0)
						offer.StartedAt = &startedAt
						offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(loan.Duration))
						offer.Status = models.LoanOfferStatusApproved
						isOffered = true
					}
				case 2:
					{
						v, err := models.ConvertString2BigInt(saleOffer.StartedAt)
						if err != nil {
							return errs.NewError(err)
						}
						startedAt := time.Unix(int64(v.Uint64())/1000, 0)
						offer.StartedAt = &startedAt
						offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(loan.Duration))
						offer.Status = models.LoanOfferStatusDone
						isOffered = true
					}
				case 3:
					{
						v, err := models.ConvertString2BigInt(saleOffer.StartedAt)
						if err != nil {
							return errs.NewError(err)
						}
						startedAt := time.Unix(int64(v.Uint64())/1000, 0)
						offer.StartedAt = &startedAt
						offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(loan.Duration))
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
