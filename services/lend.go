package services

import (
	"context"

	"github.com/czConstant/blockchain-api/bcclient"
	"github.com/czConstant/constant-nftlend-api/daos"
	"github.com/czConstant/constant-nftlend-api/errs"
	"github.com/czConstant/constant-nftlend-api/models"
	"github.com/jinzhu/gorm"
)

type NftLend struct {
	bcs   *bcclient.Client
	nlcd  *daos.Currency
	nlcld *daos.Collection
	nlad  *daos.Asset
	nlld  *daos.Loan
	nllod *daos.LoanOffer
	nlltd *daos.LoanTransaction
	nlid  *daos.Instruction
}

func NewNftLend(
	bcs *bcclient.Client,
	nlcd *daos.Currency,
	nlcld *daos.Collection,
	nlad *daos.Asset,
	nlld *daos.Loan,
	nllod *daos.LoanOffer,
	nlltd *daos.LoanTransaction,
	nlid *daos.Instruction,

) *NftLend {
	return &NftLend{
		bcs:   bcs,
		nlcd:  nlcd,
		nlcld: nlcld,
		nlad:  nlad,
		nlld:  nlld,
		nllod: nllod,
		nlltd: nlltd,
		nlid:  nlid,
	}
}

func (s *NftLend) getLendCurrency(tx *gorm.DB, address string) (*models.Currency, error) {
	c, err := s.nlcd.First(
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

func (s *NftLend) GetAssetDetail(ctx context.Context, seoURL string) (*models.Asset, error) {
	m, err := s.nlad.First(
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
					return db.Order("nfty_lend_loan_offers.id DESC")
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
	categories, count, err := s.nlcld.Find4Page(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{},
		map[string][]interface{}{
			"ListingAsset": []interface{}{
				`id in (
					select nfty_lend_asset_id
					from nfty_lend_loans
					where nfty_lend_asset_id = nfty_lend_assets.id
					  and nfty_lend_loans.status in (?)
				)`,
				[]models.LoanOfferStatus{
					models.LoanOfferStatusNew,
				},
				func(db *gorm.DB) *gorm.DB {
					return db.Order(`
					(
						select max(nfty_lend_loans.created_at)
						from nfty_lend_loans
						where nfty_lend_asset_id = nfty_lend_assets.id
						  and nfty_lend_loans.status in ('new')
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
	m, err := s.nlcld.First(
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
	currencies, err := s.nlcd.Find(
		daos.GetDBMainCtx(ctx),
		map[string][]interface{}{},
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
	ms, err := s.nlad.GetRPTListingCollection(
		daos.GetDBMainCtx(ctx),
	)
	if err != nil {
		return nil, errs.NewError(err)
	}
	return ms, nil
}
