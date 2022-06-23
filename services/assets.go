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

func (s *NftLend) GetAssetDetail(ctx context.Context, seoURL string) (*models.Asset, error) {
	m, err := s.ad.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"network in (?)": []interface{}{s.getSupportedNetworks()},
			"seo_url = ?":    []interface{}{seoURL},
		},
		map[string][]interface{}{
			"Collection": []interface{}{},
			"NewLoan": []interface{}{
				"status in (?)",
				[]models.LoanStatus{
					models.LoanStatusNew,
					models.LoanStatusCreated,
				},
			},
			"NewLoan.Currency": []interface{}{},
			"NewLoan.Offers": []interface{}{
				func(db *gorm.DB) *gorm.DB {
					return db.Order("loan_offers.id DESC")
				},
			},
			"NewLoan.ApprovedOffer": []interface{}{
				"status = ?",
				models.LoanOfferStatusApproved,
			},
		},
		[]string{"id desc"},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetAseetTransactions(ctx context.Context, assetId uint, page int, limit int) ([]*models.AssetTransaction, uint, error) {
	filters := map[string][]interface{}{}
	if assetId > 0 {
		err := s.updateAssetTransactions(ctx, assetId)
		if err != nil {
			return nil, 0, errs.NewError(err)
		}
		filters["asset_id = ?"] = []interface{}{assetId}
	}
	txns, count, err := s.atd.Find4Page(
		daos.GetDBMainCtx(ctx),
		filters,
		map[string][]interface{}{
			"Asset":            []interface{}{},
			"Asset.Collection": []interface{}{},
			"Currency":         []interface{}{},
		},
		[]string{"transaction_at desc"},
		page,
		limit,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return txns, count, nil
}

func (s *NftLend) updateAssetTransactions(ctx context.Context, assetId uint) error {
	asset, err := s.ad.FirstByID(
		daos.GetDBMainCtx(ctx),
		assetId,
		map[string][]interface{}{},
		false,
	)
	if err != nil {
		return errs.NewError(err)
	}
	if asset == nil {
		return errs.NewError(errs.ErrBadRequest)
	}
	if s.conf.Contract.ProgramID != "" {
		if asset.Network == models.NetworkSOL &&
			(asset.MagicEdenCrawAt == nil ||
				asset.MagicEdenCrawAt.Before(time.Now().Add(-24*time.Hour))) {
			c, err := s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), models.NetworkSOL, "SOL")
			if err != nil {
				return errs.NewError(err)
			}
			tokenAddress := asset.ContractAddress
			if asset.TestContractAddress != "" {
				tokenAddress = asset.TestContractAddress
			}
			rs, _ := s.stc.GetMagicEdenSaleHistories(tokenAddress)
			for _, r := range rs {
				if r.TxType == "exchange" {
					txnAt := time.Unix(r.BlockTime, 0)
					_ = s.atd.Create(
						daos.GetDBMainCtx(ctx),
						&models.AssetTransaction{
							Source:        "magiceden.io",
							Network:       models.NetworkSOL,
							AssetID:       asset.ID,
							Type:          models.AssetTransactionTypeExchange,
							Seller:        r.SellerAddress,
							Buyer:         r.BuyerAddress,
							TransactionID: r.TransactionID,
							TransactionAt: &txnAt,
							Amount:        numeric.BigFloat{*models.ConvertWeiToBigFloat(big.NewInt(int64(r.ParsedTransaction.TotalAmount)), 9)},
							CurrencyID:    c.ID,
						},
					)
				}
			}
			err = daos.WithTransaction(
				daos.GetDBMainCtx(ctx),
				func(tx *gorm.DB) error {
					asset, err := s.ad.FirstByID(
						tx,
						asset.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if asset == nil {
						return errs.NewError(errs.ErrBadRequest)
					}
					asset.MagicEdenCrawAt = helpers.TimeNow()
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
		}
		if asset.Network == models.NetworkSOL &&
			(asset.SolanartCrawAt == nil ||
				asset.SolanartCrawAt.Before(time.Now().Add(-24*time.Hour))) {
			tokenAddress := asset.ContractAddress
			if asset.TestContractAddress != "" {
				tokenAddress = asset.TestContractAddress
			}
			rs, _ := s.stc.GetSolnartSaleHistories(tokenAddress)
			for _, r := range rs {
				c, err := s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), models.NetworkSOL, r.Currency)
				if err != nil {
					return errs.NewError(err)
				}
				_ = s.atd.Create(
					daos.GetDBMainCtx(ctx),
					&models.AssetTransaction{
						Source:        "solanart.io",
						Network:       models.NetworkSOL,
						AssetID:       asset.ID,
						Type:          models.AssetTransactionTypeExchange,
						Seller:        r.SellerAddress,
						Buyer:         r.BuyerAdd,
						TransactionAt: r.Date,
						Amount:        numeric.BigFloat{*big.NewFloat(r.Price)},
						CurrencyID:    c.ID,
					},
				)
			}
			err = daos.WithTransaction(
				daos.GetDBMainCtx(ctx),
				func(tx *gorm.DB) error {
					asset, err := s.ad.FirstByID(
						tx,
						asset.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if asset == nil {
						return errs.NewError(errs.ErrBadRequest)
					}
					asset.SolanartCrawAt = helpers.TimeNow()
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
		}
		if asset.Network == models.NetworkSOL &&
			(asset.SolSeaCrawAt == nil ||
				asset.SolSeaCrawAt.Before(time.Now().Add(-24*time.Hour))) {
			tokenAddress := asset.ContractAddress
			if asset.TestContractAddress != "" {
				tokenAddress = asset.TestContractAddress
			}
			s.stc.PubSolseaMsg(
				fmt.Sprintf(`421["find","listed-archive",{"Mint":"%s","status":"SOLD"}]`, tokenAddress),
			)
			err = daos.WithTransaction(
				daos.GetDBMainCtx(ctx),
				func(tx *gorm.DB) error {
					asset, err := s.ad.FirstByID(
						tx,
						asset.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if asset == nil {
						return errs.NewError(errs.ErrBadRequest)
					}
					asset.SolSeaCrawAt = helpers.TimeNow()
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
		}
	}
	if s.conf.Contract.NearNftypawnAddress != "" {
		if asset.Network == models.NetworkNEAR &&
			(asset.ParasCrawAt == nil ||
				asset.ParasCrawAt.Before(time.Now().Add(-24*time.Hour))) {
			c, err := s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), models.NetworkNEAR, "NEAR")
			if err != nil {
				return errs.NewError(err)
			}
			contractAddress := asset.ContractAddress
			tokenID := asset.TokenID
			if asset.TestContractAddress != "" {
				contractAddress = asset.TestContractAddress
				tokenID = asset.TestTokenID
			}
			rs, _ := s.stc.GetParasSaleHistories(contractAddress, tokenID)
			for i := len(rs) - 1; i >= 0; i-- {
				r := rs[i]
				txnAt := time.Unix(r.IssuedAt/1000, 0)
				_ = s.atd.Create(
					daos.GetDBMainCtx(ctx),
					&models.AssetTransaction{
						Source:        "paras.id",
						Network:       models.NetworkNEAR,
						AssetID:       asset.ID,
						Type:          models.AssetTransactionTypeExchange,
						Seller:        r.From,
						Buyer:         r.To,
						TransactionID: r.TransactionHash,
						TransactionAt: &txnAt,
						Amount:        numeric.BigFloat{*models.ConvertWeiToBigFloat(&r.Msg.Params.Price.Int, c.Decimals)},
						CurrencyID:    c.ID,
					},
				)
			}
			err = daos.WithTransaction(
				daos.GetDBMainCtx(ctx),
				func(tx *gorm.DB) error {
					asset, err := s.ad.FirstByID(
						tx,
						asset.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if asset == nil {
						return errs.NewError(errs.ErrBadRequest)
					}
					asset.ParasCrawAt = helpers.TimeNow()
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
		}
	}
	if s.conf.Contract.MaticNftypawnAddress != "" {
		if asset.Network == models.NetworkMATIC &&
			(asset.NftbankCrawAt == nil ||
				asset.NftbankCrawAt.Before(time.Now().Add(-24*time.Hour))) {
			c, err := s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), models.NetworkMATIC, "ETH")
			if err != nil {
				return errs.NewError(err)
			}
			rs, err := s.stc.GetNftbankSaleHistories(asset.GetContractAddress(), asset.GetTokenID(), "MATIC")
			if err != nil {
				return errs.NewError(err)
			}
			for i := len(rs) - 1; i >= 0; i-- {
				r := rs[i]
				txnAt, err := time.Parse("Mon, 02 Jan 2006 15:04:05 GMT", r.BlockTimestamp)
				if err != nil {
					return errs.NewError(err)
				}
				_ = s.atd.Create(
					daos.GetDBMainCtx(ctx),
					&models.AssetTransaction{
						Source:        "nftbank.ai",
						Network:       models.NetworkMATIC,
						AssetID:       asset.ID,
						Type:          models.AssetTransactionTypeExchange,
						Seller:        r.SellerAddress,
						Buyer:         r.BuyerAddress,
						TransactionID: r.TransactionHash,
						TransactionAt: &txnAt,
						Amount:        r.SoldPriceEth,
						CurrencyID:    c.ID,
					},
				)
			}
			err = daos.WithTransaction(
				daos.GetDBMainCtx(ctx),
				func(tx *gorm.DB) error {
					asset, err := s.ad.FirstByID(
						tx,
						asset.ID,
						map[string][]interface{}{},
						true,
					)
					if err != nil {
						return errs.NewError(err)
					}
					if asset == nil {
						return errs.NewError(errs.ErrBadRequest)
					}
					asset.NftbankCrawAt = helpers.TimeNow()
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
		}
	}
	return nil
}

func (s *NftLend) solseaMsgReceived(msg string) {
	if strings.HasPrefix(msg, "431") {
		msg = strings.TrimLeft(msg, "431")
		resps := []*struct {
			Data []*struct {
				Mint      string     `json:"mint"`
				Price     uint64     `json:"price"`
				SellerKey string     `json:"sellerKey"`
				BuyerKey  string     `json:"buyerKey"`
				Status    string     `json:"status"`
				ListedAt  *time.Time `json:"listedAt"`
			} `json:"data"`
		}{}
		err := json.Unmarshal([]byte(msg), &resps)
		if err != nil {
			return
		}
		c, err := s.getLendCurrencyBySymbol(daos.GetDBMainCtx(context.Background()), models.NetworkSOL, "SOL")
		if err != nil {
			return
		}
		for _, resp := range resps {
			if resp != nil {
				for _, d := range resp.Data {
					asset, err := s.ad.First(
						daos.GetDBMainCtx(context.Background()),
						map[string][]interface{}{
							"contract_address = ?": []interface{}{d.Mint},
						},
						map[string][]interface{}{},
						[]string{"id desc"},
					)
					if err != nil {
						return
					}
					if asset == nil {
						asset, err = s.ad.First(
							daos.GetDBMainCtx(context.Background()),
							map[string][]interface{}{
								"test_contract_address = ?": []interface{}{d.Mint},
							},
							map[string][]interface{}{},
							[]string{"id desc"},
						)
						if err != nil {
							return
						}
						if asset == nil {
							return
						}
					}
					_ = s.atd.Create(
						daos.GetDBMainCtx(context.Background()),
						&models.AssetTransaction{
							Source:        "solsea.io",
							Network:       models.NetworkSOL,
							AssetID:       asset.ID,
							Type:          models.AssetTransactionTypeExchange,
							Seller:        d.SellerKey,
							Buyer:         d.BuyerKey,
							TransactionAt: d.ListedAt,
							Amount:        numeric.BigFloat{*models.ConvertWeiToBigFloat(big.NewInt(int64(d.Price)), 9)},
							CurrencyID:    c.ID,
						},
					)
				}
			}
		}
	}
}

func (s *NftLend) GetRPTAssetLoanToValue(ctx context.Context, assetID uint) (numeric.BigFloat, error) {
	v, err := s.atd.GetRPTAssetLoanToValue(
		daos.GetDBMainCtx(ctx),
		assetID,
	)
	if err != nil {
		return v, errs.NewError(err)
	}
	return v, nil
}

func (s *NftLend) GetAssetStatsInfo(ctx context.Context, assetID uint) (*serializers.AssetStatsResp, error) {
	m, err := s.ad.FirstByID(
		daos.GetDBMainCtx(ctx),
		assetID,
		map[string][]interface{}{},
		false,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if m == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	resp := &serializers.AssetStatsResp{}
	floorPice, saleCurrency, err := s.GetAssetFloorPrice(ctx, m.ID)
	if err != nil {
		return nil, errs.NewError(err)
	}
	resp.FloorPrice = floorPice
	avgPrice, err := s.atd.GetAssetAvgPrice(daos.GetDBMainCtx(ctx), m.ID)
	if err != nil {
		return nil, errs.NewError(err)
	}
	resp.AvgPrice = avgPrice
	resp.Currency = serializers.NewCurrencyResp(saleCurrency)
	return resp, nil
}

func (s *NftLend) GetAssetFloorPrice(ctx context.Context, assetID uint) (numeric.BigFloat, *models.Currency, error) {
	m, err := s.ad.FirstByID(
		daos.GetDBMainCtx(ctx),
		assetID,
		map[string][]interface{}{
			"Collection": []interface{}{},
		},
		false,
	)
	if err != nil {
		return numeric.BigFloat{}, nil, errs.NewError(err)
	}
	if m == nil {
		return numeric.BigFloat{}, nil, errs.NewError(errs.ErrBadRequest)
	}
	var saleCurrency *models.Currency
	switch m.Network {
	case models.NetworkMATIC:
		{
			saleCurrency, err = s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), models.NetworkNEAR, "NEAR")
			if err != nil {
				return numeric.BigFloat{}, nil, errs.NewError(err)
			}
		}
	case models.NetworkNEAR:
		{
			saleCurrency, err = s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), models.NetworkNEAR, "NEAR")
			if err != nil {
				return numeric.BigFloat{}, nil, errs.NewError(err)
			}
		}
	}
	if m.FloorPriceAt == nil ||
		m.FloorPriceAt.Before(time.Now().Add(-24*time.Hour)) {
		assetFloorPrice := numeric.BigFloat{*big.NewFloat(0)}
		switch m.Network {
		case models.NetworkMATIC:
			{
				nftbankStats, _ := s.stc.GetNftbankFloorPrice(m.GetContractAddress(), "MATIC")
				if nftbankStats != nil && len(nftbankStats) > 0 {
					for _, v := range nftbankStats[0].FloorPrice {
						if v.CurrencySymbol == "ETH" {
							assetFloorPrice = v.FloorPrice
						}
					}
				}
			}
		case models.NetworkNEAR:
			{
				parasStats, _ := s.stc.GetParasCollectionStats(m.Collection.ParasCollectionID)
				if parasStats != nil {
					floorPrice := models.ConvertWeiToBigFloat(&parasStats.FloorPrice.Int, saleCurrency.Decimals)
					assetFloorPrice = numeric.BigFloat{*floorPrice}
				}
			}
		}
		err = daos.WithTransaction(
			daos.GetDBMainCtx(ctx),
			func(tx *gorm.DB) error {
				m, err = s.ad.FirstByID(
					tx,
					assetID,
					map[string][]interface{}{},
					true,
				)
				if err != nil {
					return errs.NewError(err)
				}
				if m == nil {
					return errs.NewError(errs.ErrBadRequest)
				}
				m.FloorPrice = assetFloorPrice
				m.FloorPriceAt = helpers.TimeNow()
				err = s.ad.Save(
					tx,
					m,
				)
				if err != nil {
					return errs.NewError(err)
				}
				return nil
			},
		)
		if err != nil {
			return m.FloorPrice, nil, errs.NewError(err)
		}
	}
	return m.FloorPrice, saleCurrency, nil
}
