package services

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/czConstant/constant-nftylend-api/daos"
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/helpers"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

func (s *NftLend) JobUpdateCurrencyPrice(ctx context.Context) error {
	var retErr error
	symbols := []string{"ETH", "NEAR"}
	for _, symbol := range symbols {
		err := (func() error {
			price, err := s.stc.GetCoingeckoPrice(symbol)
			if err != nil {
				return errs.NewError(err)
			}
			currencies, err := s.cd.Find(
				daos.GetDBMainCtx(ctx),
				map[string][]interface{}{
					"symbol = ?": []interface{}{symbol},
				},
				map[string][]interface{}{},
				[]string{},
				0,
				999999,
			)
			if err != nil {
				return errs.NewError(err)
			}
			for _, currency := range currencies {
				currency.Price = price
				err = s.cd.Save(
					daos.GetDBMainCtx(ctx),
					currency,
				)
				if err != nil {
					return errs.NewError(err)
				}
			}
			return nil
		})()
		if err != nil {
			retErr = errs.MergeError(retErr, err)
		}
	}
	return retErr
}

func (s *NftLend) LendNftLendUpdateBlock(ctx context.Context, block uint64) error {
	err := s.bcs.Solana.NftLendUpdateBlock(block)
	if err != nil {
		return errs.NewError(err)
	}
	return nil
}

func (s *NftLend) MoralisGetNFTs(ctx context.Context, chain string, address string, cursor string, limit int) (interface{}, error) {
	rs, err := s.mc.GetNFTs(chain, address, cursor, limit)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return rs, nil
}

func (s *NftLend) ProcessSolanaInstruction(ctx context.Context, insId uint) error {
	emailQueue := []*models.EmailQueue{}
	var loadAssetTransactionForId uint
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			ins, err := s.id.FirstByID(
				tx,
				insId,
				map[string][]interface{}{},
				true,
			)
			if err != nil {
				return errs.NewError(err)
			}
			switch ins.Network {
			case models.NetworkSOL:
				{
					switch ins.Instruction {
					case "InitLoan":
						{
							var req struct {
								LoanPrincipalAmount   uint64 `json:"loan_principal_amount"`
								LoanDuration          uint64 `json:"loan_duration"`
								InterestRate          uint32 `json:"interest_rate"`
								NftCollateralContract string `json:"nft_collateral_contract"`
								LoanCurrency          string `json:"loan_currency"`
								BorrowerAccount       string `json:"borrower_account"`
								TempNftAccount        string `json:"temp_nft_account"`
								TokenToReceiveAccount string `json:"token_to_receive_account"`
								LoanInfoAccount       string `json:"loan_info_account"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							currency, err := s.getLendCurrency(tx, req.LoanCurrency)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address =?": []interface{}{req.LoanInfoAccount},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan != nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							asset, err := s.ad.First(
								tx,
								map[string][]interface{}{
									"contract_address =?": []interface{}{req.NftCollateralContract},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if asset == nil {
								// parse info and new collection
								meta, err := s.bcs.Solana.GetMetadata(req.NftCollateralContract)
								if err != nil {
									return errs.NewError(err)
								}
								metaInfo, err := s.bcs.Solana.GetMetadataInfo(helpers.ConvertImageDataURL(meta.Data.Uri))
								if err != nil {
									return errs.NewError(err)
								}
								collection, tokenId, err := s.getSolanaCollectionVerified(
									tx,
									req.NftCollateralContract,
									meta,
									metaInfo,
								)
								if err != nil {
									return errs.NewError(err)
								}
								if collection == nil {
									// return errs.NewError(errs.ErrBadRequest)

									collectionName := metaInfo.Collection.Name
									if collectionName == "" {
										names := strings.Split(metaInfo.Name, "#")
										if len(names) >= 2 {
											collectionName = strings.TrimSpace(names[0])
										}
									}
									if collectionName == "" {
										return errs.NewError(errs.ErrBadRequest)
									}
									collection, err = s.cld.First(
										tx,
										map[string][]interface{}{
											"name = ?": []interface{}{collectionName},
										},
										map[string][]interface{}{},
										[]string{},
									)
									if err != nil {
										return errs.NewError(err)
									}
									if collection == nil {
										collection = &models.Collection{
											Network:     models.NetworkSOL,
											SeoURL:      helpers.MakeSeoURL(fmt.Sprintf("%s-%s", models.NetworkSOL, collectionName)),
											Name:        collectionName,
											Description: metaInfo.Description,
											Enabled:     true,
										}
										err = s.cld.Create(
											tx,
											collection,
										)
										if err != nil {
											return errs.NewError(err)
										}
									}
								}
								var sellerFeeBasisPoints int64
								switch metaInfo.SellerFeeBasisPoints.(type) {
								case string:
									{
										sellerFeeBasisPoints, _ = strconv.ParseInt(metaInfo.SellerFeeBasisPoints.(string), 10, 64)
									}
								case float64:
									{
										sellerFeeBasisPoints = int64(metaInfo.SellerFeeBasisPoints.(float64))
									}
								}
								sellerFeeRate, _ := models.ConvertWeiToBigFloat(big.NewInt(sellerFeeBasisPoints), 4).Float64()
								attributes, _ := json.Marshal(metaInfo.Attributes)
								metaJson, err := json.Marshal(metaInfo)
								if err != nil {
									return errs.NewError(err)
								}
								asset = &models.Asset{
									Network:               models.NetworkSOL,
									SeoURL:                helpers.MakeSeoURL(fmt.Sprintf("%s-%s", models.NetworkSOL, req.NftCollateralContract)),
									ContractAddress:       req.NftCollateralContract,
									CollectionID:          collection.ID,
									Symbol:                metaInfo.Symbol,
									Name:                  metaInfo.Name,
									TokenURL:              metaInfo.Image,
									ExternalUrl:           metaInfo.ExternalUrl,
									SellerFeeRate:         sellerFeeRate,
									Attributes:            string(attributes),
									MetaJson:              string(metaJson),
									MetaJsonUrl:           meta.Data.Uri,
									OriginNetwork:         collection.OriginNetwork,
									OriginContractAddress: collection.OriginContractAddress,
									OriginTokenID:         tokenId,
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
							principalAmount := models.ConvertWeiToBigFloat(big.NewInt(int64(req.LoanPrincipalAmount)), currency.Decimals)
							interestRate, _ := models.ConvertWeiToBigFloat(big.NewInt(int64(req.InterestRate)), 4).Float64()
							loan = &models.Loan{
								Network:          models.NetworkSOL,
								DataLoanAddress:  req.LoanInfoAccount,
								DataAssetAddress: req.TempNftAccount,
								Owner:            req.BorrowerAccount,
								PrincipalAmount:  numeric.BigFloat{*principalAmount},
								InterestRate:     interestRate,
								Duration:         uint(req.LoanDuration),
								StartedAt:        ins.BlockTime,
								ExpiredAt:        helpers.TimeAdd(*ins.BlockTime, time.Duration(req.LoanDuration)*time.Second),
								CurrencyID:       currency.ID,
								AssetID:          asset.ID,
								Status:           models.LoanStatusNew,
								InitTxHash:       ins.TransactionHash,
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
							loadAssetTransactionForId = asset.ID
						}
					case "MakeOffer":
						{
							var req struct {
								LoanID              string `json:"loan_id"`
								LoanPrincipalAmount uint64 `json:"loan_principal_amount"`
								LoanDuration        uint64 `json:"loan_duration"`
								InterestRate        uint64 `json:"interest_rate"`
								LoanCurrency        string `json:"loan_currency"`
								LenderAccount       string `json:"lender_account"`
								TempTokenAccount    string `json:"temp_token_account"`
								OfferInfoAccount    string `json:"offer_info_account"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							currency, err := s.getLendCurrency(tx, req.LoanCurrency)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address =?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{
									"Offers": []interface{}{},
								},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"data_offer_address =?": []interface{}{req.OfferInfoAccount},
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
							// parse info and new collection
							principalAmount := models.ConvertWeiToBigFloat(big.NewInt(int64(req.LoanPrincipalAmount)), currency.Decimals)
							interestRate, _ := models.ConvertWeiToBigFloat(big.NewInt(int64(req.InterestRate)), 4).Float64()
							offer = &models.LoanOffer{
								Network:             models.NetworkSOL,
								LoanID:              loan.ID,
								Lender:              req.LenderAccount,
								PrincipalAmount:     numeric.BigFloat{*principalAmount},
								InterestRate:        interestRate,
								Duration:            uint(req.LoanDuration),
								Status:              models.LoanOfferStatusNew,
								DataOfferAddress:    req.OfferInfoAccount,
								DataCurrencyAddress: req.TempTokenAccount,
								MakeTxHash:          ins.TransactionHash,
							}
							if loan.Status != models.LoanStatusNew {
								offer.Status = models.LoanOfferStatusRejected
							}
							err = s.lod.Create(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
							emailQueue = append(
								emailQueue,
								&models.EmailQueue{
									EmailType: models.EMAIL_BORROWER_OFFER_NEW,
									ObjectID:  offer.ID,
								},
							)
						}
					case "AcceptOffer":
						{
							var req struct {
								LoanID              string `json:"loan_id"`
								OfferID             string `json:"offer_id"`
								LoanPrincipalAmount uint64 `json:"loan_principal_amount"`
								LoanDuration        uint64 `json:"loan_duration"`
								InterestRate        uint64 `json:"interest_rate"`
								LoanCurrency        string `json:"loan_currency"`
								LenderAccount       string `json:"lender_account"`
								TempTokenAccount    string `json:"temp_token_account"`
								OfferInfoAccount    string `json:"offer_info_account"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address =?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{},
								[]string{},
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
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"data_offer_address =?": []interface{}{req.OfferID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if offer.Status != models.LoanOfferStatusNew {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.StartedAt = ins.BlockTime
							offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(offer.Duration))
							offer.Status = models.LoanOfferStatusApproved
							offer.AcceptTxHash = ins.TransactionHash
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
							loan.Lender = offer.Lender
							loan.OfferStartedAt = offer.StartedAt
							loan.OfferDuration = offer.Duration
							loan.OfferExpiredAt = offer.ExpiredAt
							loan.OfferOverdueAt = helpers.TimeAdd(*loan.OfferExpiredAt, (2 * 24 * time.Hour))
							loan.OfferPrincipalAmount = offer.PrincipalAmount
							loan.OfferInterestRate = offer.InterestRate
							loan.Status = models.LoanStatusCreated
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							for _, otherOffer := range loan.Offers {
								if otherOffer.ID != offer.ID {
									if otherOffer.Status == models.LoanOfferStatusNew {
										otherOffer.Status = models.LoanOfferStatusRejected
										err = s.lod.Save(
											tx,
											otherOffer,
										)
										if err != nil {
											return errs.NewError(err)
										}
									}
								}
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
							emailQueue = append(
								emailQueue,
								&models.EmailQueue{
									EmailType: models.EMAIL_BORROWER_LOAN_STARTED,
									ObjectID:  loan.ID,
								},
								&models.EmailQueue{
									EmailType: models.EMAIL_LENDER_OFFER_STARTED,
									ObjectID:  offer.ID,
								},
							)
						}
					case "CancelLoan":
						{
							var req struct {
								LoanID string `json:"loan_id"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address = ?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{
									"Offers": []interface{}{},
								},
								[]string{},
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
							loan.FinishedAt = ins.BlockTime
							loan.Status = models.LoanStatusCancelled
							loan.CancelTxHash = ins.TransactionHash
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							for _, otherOffer := range loan.Offers {
								if otherOffer.Status == models.LoanOfferStatusNew {
									otherOffer.Status = models.LoanOfferStatusRejected
									err = s.lod.Save(
										tx,
										otherOffer,
									)
									if err != nil {
										return errs.NewError(err)
									}
								}
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					case "CancelOffer":
						{
							var req struct {
								OfferID string `json:"offer_id"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"data_offer_address =?": []interface{}{req.OfferID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if offer.Status != models.LoanOfferStatusNew &&
								offer.Status != models.LoanOfferStatusRejected {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.FinishedAt = ins.BlockTime
							offer.Status = models.LoanOfferStatusCancelled
							offer.CancelTxHash = ins.TransactionHash
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					case "PayLoan":
						{
							var req struct {
								LoanID    string `json:"loan_id"`
								OfferID   string `json:"offer_id"`
								PayAmount uint64 `json:"pay_amount"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address =?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if loan.Status != models.LoanStatusCreated {
								return errs.NewError(errs.ErrBadRequest)
							}
							currency, err := s.cd.FirstByID(
								tx,
								loan.CurrencyID,
								map[string][]interface{}{},
								false,
							)
							if err != nil {
								return errs.NewError(err)
							}
							payAmount := models.ConvertWeiToBigFloat(big.NewInt(int64(req.PayAmount)), currency.Decimals)
							loan.RepaidAmount = numeric.BigFloat{*payAmount}
							loan.FinishedAt = ins.BlockTime
							loan.Status = models.LoanStatusDone
							loan.PayTxHash = ins.TransactionHash
							loan.FeeRate = 0.01
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"data_offer_address =?": []interface{}{req.OfferID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.RepaidAt = ins.BlockTime
							offer.RepaidAmount = numeric.BigFloat{*payAmount}
							offer.Status = models.LoanOfferStatusRepaid
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					case "LiquidateLoan":
						{
							var req struct {
								LoanID  string `json:"loan_id"`
								OfferID string `json:"offer_id"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address =?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if loan.Status != models.LoanStatusCreated {
								return errs.NewError(errs.ErrBadRequest)
							}
							loan.FinishedAt = ins.BlockTime
							loan.Status = models.LoanStatusLiquidated
							loan.LiquidateTxHash = ins.TransactionHash
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"data_offer_address =?": []interface{}{req.OfferID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.Status = models.LoanOfferStatusLiquidated
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					case "CloseOffer":
						{
							var req struct {
								OfferID string `json:"offer_id"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"data_offer_address =?": []interface{}{req.OfferID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if offer.Status != models.LoanOfferStatusRepaid {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.FinishedAt = ins.BlockTime
							offer.Status = models.LoanOfferStatusDone
							offer.CloseTxHash = ins.TransactionHash
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
						}
					case "Order":
						{
							var req struct {
								LoanID           string `json:"loan_id"`
								LenderAccount    string `json:"lender_account"`
								TempTokenAccount string `json:"temp_token_account"`
								OfferInfoAccount string `json:"offer_info_account"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"data_loan_address = ?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{
									"Offers": []interface{}{},
								},
								[]string{},
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
							offer := &models.LoanOffer{
								Network:             models.NetworkSOL,
								LoanID:              loan.ID,
								Lender:              req.LenderAccount,
								PrincipalAmount:     loan.PrincipalAmount,
								InterestRate:        loan.InterestRate,
								Duration:            loan.Duration,
								DataOfferAddress:    req.OfferInfoAccount,
								DataCurrencyAddress: req.TempTokenAccount,
								StartedAt:           ins.BlockTime,
								ExpiredAt:           helpers.TimeAdd(*ins.BlockTime, time.Second*time.Duration(loan.Duration)),
								Status:              models.LoanOfferStatusApproved,
								MakeTxHash:          ins.TransactionHash,
								AcceptTxHash:        ins.TransactionHash,
							}
							err = s.lod.Create(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
							loan.Lender = offer.Lender
							loan.OfferStartedAt = offer.StartedAt
							loan.OfferDuration = offer.Duration
							loan.OfferExpiredAt = offer.ExpiredAt
							loan.OfferOverdueAt = helpers.TimeAdd(*loan.OfferExpiredAt, (2 * 24 * time.Hour))
							loan.OfferPrincipalAmount = offer.PrincipalAmount
							loan.OfferInterestRate = offer.InterestRate
							loan.Status = models.LoanStatusCreated
							loan.InitTxHash = ins.TransactionHash
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							for _, otherOffer := range loan.Offers {
								if otherOffer.ID != offer.ID {
									if otherOffer.Status == models.LoanOfferStatusNew {
										otherOffer.Status = models.LoanOfferStatusRejected
										err = s.lod.Save(
											tx,
											otherOffer,
										)
										if err != nil {
											return errs.NewError(err)
										}
									}
								}
							}
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
									TxHash:          ins.TransactionHash,
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
				}
			case models.NetworkMATIC,
				models.NetworkAVAX,
				models.NetworkBSC,
				models.NetworkBOBA:
				{
					switch ins.Instruction {
					case "LoanStarted":
						{
							var req struct {
								LoanID           string `json:"loan_id"`
								Borrower         string `json:"borrower"`
								BorrowerNonceHex string `json:"borrower_nonce_hex"`
								Lender           string `json:"lender"`
								LenderNonceHex   string `json:"lender_nonce_hex"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							req.Lender = strings.ToLower(req.Lender)
							req.Borrower = strings.ToLower(req.Borrower)
							req.Lender = strings.ToLower(req.Lender)
							req.LenderNonceHex = strings.ToLower(req.LenderNonceHex)
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"network = ?":   []interface{}{ins.Network},
									"owner = ?":     []interface{}{req.Borrower},
									"nonce_hex = ?": []interface{}{req.BorrowerNonceHex},
								},
								map[string][]interface{}{
									"Offers": []interface{}{},
								},
								[]string{},
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
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"network = ?":  []interface{}{ins.Network},
									"lender = ?":   []interface{}{req.Lender},
									"nonce_hex =?": []interface{}{req.LenderNonceHex},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if offer.Status != models.LoanOfferStatusNew {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.StartedAt = ins.BlockTime
							offer.ExpiredAt = helpers.TimeAdd(*offer.StartedAt, time.Second*time.Duration(offer.Duration))
							offer.Status = models.LoanOfferStatusApproved
							offer.AcceptTxHash = ins.TransactionHash
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
							loan.DataLoanAddress = req.LoanID
							loan.Lender = offer.Lender
							loan.OfferStartedAt = offer.StartedAt
							loan.OfferDuration = offer.Duration
							loan.OfferExpiredAt = offer.ExpiredAt
							loan.OfferOverdueAt = helpers.TimeAdd(*loan.OfferExpiredAt, (2 * 24 * time.Hour))
							loan.OfferPrincipalAmount = offer.PrincipalAmount
							loan.OfferInterestRate = offer.InterestRate
							loan.InitTxHash = ins.TransactionHash
							loan.Status = models.LoanStatusCreated
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							for _, otherOffer := range loan.Offers {
								if otherOffer.ID != offer.ID {
									if otherOffer.Status == models.LoanOfferStatusNew {
										otherOffer.Status = models.LoanOfferStatusRejected
										err = s.lod.Save(
											tx,
											otherOffer,
										)
										if err != nil {
											return errs.NewError(err)
										}
									}
								}
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
							emailQueue = append(
								emailQueue,
								&models.EmailQueue{
									EmailType: models.EMAIL_BORROWER_LOAN_STARTED,
									ObjectID:  loan.ID,
								},
								&models.EmailQueue{
									EmailType: models.EMAIL_LENDER_OFFER_STARTED,
									ObjectID:  offer.ID,
								},
							)
						}
					case "CancelNonce":
						{
							var req struct {
								Sender   string `json:"sender"`
								NonceHex string `json:"nonce_hex"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							req.NonceHex = strings.ToLower(req.NonceHex)
							req.Sender = strings.ToLower(req.Sender)
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"network = ?":   []interface{}{ins.Network},
									"owner = ?":     []interface{}{req.Sender},
									"nonce_hex = ?": []interface{}{req.NonceHex},
								},
								map[string][]interface{}{
									"Offers": []interface{}{},
								},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan != nil {
								if loan.Status != models.LoanStatusNew {
									return errs.NewError(errs.ErrBadRequest)
								}
								loan.CancelTxHash = ins.TransactionHash
								loan.Status = models.LoanStatusCancelled
								loan.FinishedAt = helpers.TimeNow()
								err = s.ld.Save(
									tx,
									loan,
								)
								if err != nil {
									return errs.NewError(err)
								}
								for _, otherOffer := range loan.Offers {
									if otherOffer.Status == models.LoanOfferStatusNew {
										otherOffer.Status = models.LoanOfferStatusRejected
										err = s.lod.Save(
											tx,
											otherOffer,
										)
										if err != nil {
											return errs.NewError(err)
										}
									}
								}
							} else {
								offer, err := s.lod.First(
									tx,
									map[string][]interface{}{
										"network = ?":  []interface{}{ins.Network},
										"lender = ?":   []interface{}{req.Sender},
										"nonce_hex =?": []interface{}{req.NonceHex},
									},
									map[string][]interface{}{},
									[]string{},
								)
								if err != nil {
									return errs.NewError(err)
								}
								if offer != nil {
									if offer.Status != models.LoanOfferStatusNew {
										return errs.NewError(errs.ErrBadRequest)
									}
									offer.Status = models.LoanOfferStatusCancelled
									offer.CancelTxHash = ins.TransactionHash
									err = s.lod.Save(
										tx,
										offer,
									)
									if err != nil {
										return errs.NewError(err)
									}
								}
							}
						}
					case "LoanRepaid":
						{
							var req struct {
								LoanID string `json:"loan_id"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"network = ?":           []interface{}{ins.Network},
									"data_loan_address = ?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if loan.Status != models.LoanStatusCreated {
								return errs.NewError(errs.ErrBadRequest)
							}
							currency, err := s.cd.FirstByID(
								tx,
								loan.CurrencyID,
								map[string][]interface{}{},
								false,
							)
							if err != nil {
								return errs.NewError(err)
							}
							payAmount := models.ConvertWeiToBigFloat(big.NewInt(0), currency.Decimals)
							loan.RepaidAmount = numeric.BigFloat{*payAmount}
							loan.FinishedAt = ins.BlockTime
							loan.Status = models.LoanStatusDone
							loan.PayTxHash = ins.TransactionHash
							loan.FeeRate = 0.01
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"loan_id = ?": []interface{}{loan.ID},
									"status = ?":  []interface{}{models.LoanOfferStatusApproved},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.RepaidAt = ins.BlockTime
							offer.RepaidAmount = numeric.BigFloat{*payAmount}
							offer.Status = models.LoanOfferStatusRepaid
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
							emailQueue = append(
								emailQueue,
								&models.EmailQueue{
									EmailType: models.EMAIL_LENDER_LOAN_REPAID,
									ObjectID:  offer.ID,
								},
							)
						}
					case "LoanLiquidated":
						{
							var req struct {
								LoanID string `json:"loan_id"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"network = ?":           []interface{}{ins.Network},
									"data_loan_address = ?": []interface{}{req.LoanID},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if loan == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							if loan.Status != models.LoanStatusCreated {
								return errs.NewError(errs.ErrBadRequest)
							}
							currency, err := s.cd.FirstByID(
								tx,
								loan.CurrencyID,
								map[string][]interface{}{},
								false,
							)
							if err != nil {
								return errs.NewError(err)
							}
							payAmount := models.ConvertWeiToBigFloat(big.NewInt(0), currency.Decimals)
							loan.RepaidAmount = numeric.BigFloat{*payAmount}
							loan.FinishedAt = ins.BlockTime
							loan.Status = models.LoanStatusLiquidated
							loan.LiquidateTxHash = ins.TransactionHash
							loan.FeeRate = 0.01
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							offer, err := s.lod.First(
								tx,
								map[string][]interface{}{
									"loan_id = ?": []interface{}{loan.ID},
									"status = ?":  []interface{}{models.LoanOfferStatusApproved},
								},
								map[string][]interface{}{},
								[]string{},
							)
							if err != nil {
								return errs.NewError(err)
							}
							if offer == nil {
								return errs.NewError(errs.ErrBadRequest)
							}
							offer.RepaidAt = ins.BlockTime
							offer.RepaidAmount = numeric.BigFloat{*payAmount}
							offer.Status = models.LoanOfferStatusLiquidated
							err = s.lod.Save(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
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
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
							emailQueue = append(
								emailQueue,
								&models.EmailQueue{
									EmailType: models.EMAIL_BORROWER_LOAN_LIQUIDATED,
									ObjectID:  loan.ID,
								},
							)
						}
					case "OfferNow":
						{
							var req struct {
								LoanID           string `json:"loan_id"`
								Borrower         string `json:"borrower"`
								BorrowerNonceHex string `json:"borrower_nonce_hex"`
								Lender           string `json:"lender"`
								LenderNonceHex   string `json:"lender_nonce_hex"`
							}
							err = json.Unmarshal([]byte(ins.Data), &req)
							if err != nil {
								return errs.NewError(err)
							}
							req.Lender = strings.ToLower(req.Lender)
							req.Borrower = strings.ToLower(req.Borrower)
							req.Lender = strings.ToLower(req.Lender)
							req.LenderNonceHex = strings.ToLower(req.LenderNonceHex)
							loan, err := s.ld.First(
								tx,
								map[string][]interface{}{
									"network = ?":   []interface{}{ins.Network},
									"owner = ?":     []interface{}{req.Borrower},
									"nonce_hex = ?": []interface{}{req.BorrowerNonceHex},
								},
								map[string][]interface{}{
									"Offers": []interface{}{},
								},
								[]string{},
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
							offer := &models.LoanOffer{
								Network:         ins.Network,
								LoanID:          loan.ID,
								Lender:          req.Lender,
								PrincipalAmount: loan.PrincipalAmount,
								InterestRate:    loan.InterestRate,
								Duration:        loan.Duration,
								StartedAt:       ins.BlockTime,
								ExpiredAt:       helpers.TimeAdd(*ins.BlockTime, time.Second*time.Duration(loan.Duration)),
								Status:          models.LoanOfferStatusApproved,
								MakeTxHash:      ins.TransactionHash,
								AcceptTxHash:    ins.TransactionHash,
								NonceHex:        req.LenderNonceHex,
							}
							err = s.lod.Create(
								tx,
								offer,
							)
							if err != nil {
								return errs.NewError(err)
							}
							loan.DataLoanAddress = req.LoanID
							loan.Lender = offer.Lender
							loan.OfferStartedAt = offer.StartedAt
							loan.OfferDuration = offer.Duration
							loan.OfferExpiredAt = offer.ExpiredAt
							loan.OfferOverdueAt = helpers.TimeAdd(*loan.OfferExpiredAt, (2 * 24 * time.Hour))
							loan.OfferPrincipalAmount = offer.PrincipalAmount
							loan.OfferInterestRate = offer.InterestRate
							loan.Status = models.LoanStatusCreated
							loan.InitTxHash = ins.TransactionHash
							err = s.ld.Save(
								tx,
								loan,
							)
							if err != nil {
								return errs.NewError(err)
							}
							for _, otherOffer := range loan.Offers {
								if otherOffer.ID != offer.ID {
									if otherOffer.Status == models.LoanOfferStatusNew {
										otherOffer.Status = models.LoanOfferStatusRejected
										err = s.lod.Save(
											tx,
											otherOffer,
										)
										if err != nil {
											return errs.NewError(err)
										}
									}
								}
							}
							err = s.ltd.Create(
								tx,
								&models.LoanTransaction{
									Network:         ins.Network,
									Type:            models.LoanTransactionTypeOffered,
									LoanID:          loan.ID,
									Borrower:        loan.Owner,
									Lender:          offer.Lender,
									PrincipalAmount: offer.PrincipalAmount,
									InterestRate:    offer.InterestRate,
									StartedAt:       offer.StartedAt,
									Duration:        offer.Duration,
									ExpiredAt:       offer.ExpiredAt,
									TxHash:          ins.TransactionHash,
								},
							)
							if err != nil {
								return errs.NewError(err)
							}
							emailQueue = append(
								emailQueue,
								&models.EmailQueue{
									EmailType: models.EMAIL_BORROWER_LOAN_STARTED,
									ObjectID:  loan.ID,
								},
								&models.EmailQueue{
									EmailType: models.EMAIL_LENDER_OFFER_STARTED,
									ObjectID:  offer.ID,
								},
							)
						}
					default:
						{
							return errs.NewError(errs.ErrBadRequest)
						}
					}
				}
			default:
				{
					return errs.NewError(errs.ErrBadRequest)
				}
			}
			ins.Status = "done"
			err = s.id.Save(
				tx,
				ins,
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
	if loadAssetTransactionForId > 0 {
		s.updateAssetTransactions(ctx, loadAssetTransactionForId)
	}
	{
		s.EmailForReference(ctx, emailQueue)
	}
	return nil
}

func (s *NftLend) InternalHookSolanaInstruction(ctx context.Context, network models.Network, blockNumber uint64, blockTime uint64, transactionHash string, transactionIndex uint, instructionIndex uint, program string, instruction string, data interface{}) error {
	dataJson, err := json.Marshal(&data)
	if err != nil {
		return errs.NewError(err)
	}
	var isProcess bool
	var ins *models.Instruction
	err = daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			ins, err = s.id.First(
				tx,
				map[string][]interface{}{
					"transaction_hash = ?":  []interface{}{transactionHash},
					"instruction_index = ?": []interface{}{instructionIndex},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if ins != nil {
				if ins.Status == "new" {
					isProcess = true
				}
				return nil
			}
			bt := time.Unix(int64(blockTime), 0)
			ins = &models.Instruction{
				Network:          network,
				BlockNumber:      blockNumber,
				BlockTime:        &bt,
				TransactionHash:  transactionHash,
				TransactionIndex: transactionIndex,
				InstructionIndex: instructionIndex,
				Program:          program,
				Instruction:      instruction,
				Data:             string(dataJson),
				Status:           "new",
			}
			err = s.id.Create(
				tx,
				ins,
			)
			if err != nil {
				return errs.NewError(err)
			}
			isProcess = true
			return nil
		},
	)
	if err != nil {
		return errs.NewError(err)
	}
	if isProcess {
		err = s.ProcessSolanaInstruction(ctx, ins.ID)
		if err != nil {
			return errs.NewError(err)
		}
	}
	return nil
}

func (s *NftLend) UpdateAssetInfo(ctx context.Context, address string) error {
	err := daos.WithTransaction(
		daos.GetDBMainCtx(ctx),
		func(tx *gorm.DB) error {
			asset, err := s.ad.First(
				tx,
				map[string][]interface{}{
					"contract_address =?": []interface{}{address},
				},
				map[string][]interface{}{},
				[]string{},
			)
			if err != nil {
				return errs.NewError(err)
			}
			if asset == nil {
				return errs.NewError(errs.ErrBadRequest)
			}
			// parse info and new collection
			meta, err := s.bcs.Solana.GetMetadata(asset.ContractAddress)
			if err != nil {
				return errs.NewError(err)
			}
			metaInfo, err := s.bcs.Solana.GetMetadataInfo(meta.Data.Uri)
			if err != nil {
				return errs.NewError(err)
			}
			var sellerFeeBasisPoints int64
			switch metaInfo.SellerFeeBasisPoints.(type) {
			case string:
				{
					sellerFeeBasisPoints, _ = strconv.ParseInt(metaInfo.SellerFeeBasisPoints.(string), 10, 64)
				}
			case float64:
				{
					sellerFeeBasisPoints = int64(metaInfo.SellerFeeBasisPoints.(float64))
				}
			}
			metaJson, err := json.Marshal(metaInfo)
			if err != nil {
				return errs.NewError(err)
			}
			sellerFeeRate, _ := models.ConvertWeiToBigFloat(big.NewInt(sellerFeeBasisPoints), 4).Float64()
			attributes, _ := json.Marshal(metaInfo.Attributes)
			asset.Symbol = metaInfo.Symbol
			asset.Name = metaInfo.Name
			asset.TokenURL = metaInfo.Image
			asset.ExternalUrl = metaInfo.ExternalUrl
			asset.SellerFeeRate = sellerFeeRate
			asset.Attributes = string(attributes)
			asset.MetaJson = string(metaJson)
			asset.MetaJsonUrl = meta.Data.Uri
			err = s.ad.Save(
				tx,
				asset,
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

func (s *NftLend) JobEvmNftypawnFilterLogs(ctx context.Context, network models.Network, block uint64) error {
	var retErr error
	evmClient := s.getEvmClientByNetwork(network)
	if evmClient != nil {
		contractAddress := s.getEvmContractAddress(network)
		if contractAddress != "" {
			resps, err := s.getEvmClientByNetwork(network).NftypawnFilterLogs(s.getEvmContractAddress(network), block)
			if err != nil {
				return errs.NewError(err)
			}
			for _, resp := range resps {
				err = s.InternalHookSolanaInstruction(
					ctx,
					network,
					uint64(resp.BlockNumber),
					uint64(resp.Timestamp),
					resp.Hash,
					resp.Index,
					resp.Index,
					"",
					resp.Event,
					resp.Data,
				)
				if err != nil {
					retErr = errs.MergeError(retErr, err)
				}
			}
		}
	}
	return retErr
}
