package services

import (
	"context"
	"math/big"
	"time"

	"github.com/czConstant/blockchain-api/bcclient"
	"github.com/czConstant/constant-nftlend-api/daos"
	"github.com/czConstant/constant-nftlend-api/errs"
	"github.com/czConstant/constant-nftlend-api/models"
	"github.com/czConstant/constant-nftlend-api/services/3rd/saletrack"
	"github.com/czConstant/constant-nftlend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type NftLend struct {
	bcs *bcclient.Client
	stc *saletrack.Client
	cd  *daos.Currency
	cld *daos.Collection
	ad  *daos.Asset
	atd *daos.AssetTransaction
	ld  *daos.Loan
	lod *daos.LoanOffer
	ltd *daos.LoanTransaction
	id  *daos.Instruction
}

func NewNftLend(
	bcs *bcclient.Client,
	stc *saletrack.Client,
	cd *daos.Currency,
	cld *daos.Collection,
	ad *daos.Asset,
	atd *daos.AssetTransaction,
	ld *daos.Loan,
	lod *daos.LoanOffer,
	ltd *daos.LoanTransaction,
	id *daos.Instruction,

) *NftLend {
	return &NftLend{
		bcs: bcs,
		stc: stc,
		cd:  cd,
		cld: cld,
		ad:  ad,
		atd: atd,
		ld:  ld,
		lod: lod,
		ltd: ltd,
		id:  id,
	}
}

func (s *NftLend) getLendCurrency(tx *gorm.DB, address string) (*models.Currency, error) {
	c, err := s.cd.First(
		tx,
		map[string][]interface{}{
			"contract_address = ?": []interface{}{address},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if c == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return c, nil
}

func (s *NftLend) getLendCurrencyBySymbol(tx *gorm.DB, symbol string) (*models.Currency, error) {
	c, err := s.cd.First(
		tx,
		map[string][]interface{}{
			"symbol = ?": []interface{}{symbol},
		},
		map[string][]interface{}{},
		[]string{},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	if c == nil {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return c, nil
}

func (s *NftLend) GetAssetDetail(ctx context.Context, seoURL string) (*models.Asset, error) {
	m, err := s.ad.First(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"seo_url = ?": []interface{}{seoURL},
		},
		map[string][]interface{}{
			"Collection": []interface{}{},
			"NewLoan": []interface{}{
				"status = ?",
				models.LoanStatusNew,
			},
			"NewLoan.Currency": []interface{}{},
			"NewLoan.Offers": []interface{}{
				func(db *gorm.DB) *gorm.DB {
					return db.Order("loan_offers.id DESC")
				},
			},
		},
		[]string{"id desc"},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetCollections(ctx context.Context, page int, limit int) ([]*models.Collection, uint, error) {
	categories, count, err := s.cld.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{},
		map[string][]interface{}{
			"ListingAsset": []interface{}{
				`id in (
					select asset_id
					from loans
					where asset_id = assets.id
					  and loans.status in (?)
				)`,
				[]models.LoanOfferStatus{
					models.LoanOfferStatusNew,
				},
				func(db *gorm.DB) *gorm.DB {
					return db.Order(`
					(
						select max(loans.created_at)
						from loans
						where asset_id = assets.id
						  and loans.status in ('new')
					) desc
					`)
				},
			},
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
			"seo_url = ?": []interface{}{seoURL},
		},
		map[string][]interface{}{},
		[]string{"id desc"},
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return m, nil
}

func (s *NftLend) GetCurrencies(ctx context.Context) ([]*models.Currency, error) {
	currencies, err := s.cd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{
			"enabled = ?": []interface{}{true},
		},
		map[string][]interface{}{},
		[]string{"id desc"},
		0,
		99999999,
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return currencies, nil
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

func (s *NftLend) GetAseetTransactions(ctx context.Context, assetId uint, page int, limit int) ([]*models.AssetTransaction, uint, error) {
	asset, err := s.ad.FirstByID(
		daos.GetDBMainCtx(ctx),
		assetId,
		map[string][]interface{}{},
		false,
	)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	if asset == nil {
		return nil, 0, errs.NewError(errs.ErrBadRequest)
	}
	if asset.MagicEdenCrawAt == nil ||
		asset.MagicEdenCrawAt.Before(time.Now().Add(-24*time.Hour)) {
		c, err := s.getLendCurrencyBySymbol(daos.GetDBMainCtx(ctx), "SOL")
		if err != nil {
			return nil, 0, errs.NewError(err)
		}
		tokenAddress := asset.ContractAddress
		if asset.TestContractAddress == "" {
			tokenAddress = asset.TestContractAddress
		}
		rs, _ := s.stc.GetMagicEdenSaleHistories(tokenAddress)
		for _, r := range rs {
			txnAt := time.Unix(r.BlockTime, 0)
			m := models.AssetTransaction{
				Network:       models.ChainSOL,
				AssetID:       asset.ID,
				Type:          models.AssetTransactionTypeExchange,
				Seller:        r.SellerAddress,
				Buyer:         r.BuyerAddress,
				TransactionAt: &txnAt,
				Amount:        numeric.BigFloat{*models.ConvertWeiToBigFloat(big.NewInt(int64(r.ParsedTransaction.TotalAmount)), 9)},
				CurrencyID:    c.ID,
			}
			_ = s.atd.Create(
				daos.GetDBMainCtx(ctx),
				m,
			)
		}
	}
	return nil, 0, nil
}
