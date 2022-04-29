package daos

import (
	"math/big"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type AssetTransaction struct {
	DAO
}

func (d *AssetTransaction) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.AssetTransaction, error) {
	var m models.AssetTransaction
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *AssetTransaction) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.AssetTransaction, error) {
	var m models.AssetTransaction
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *AssetTransaction) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.AssetTransaction, error) {
	var ms []*models.AssetTransaction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *AssetTransaction) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.AssetTransaction, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.AssetTransaction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.AssetTransaction{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}

func (d *AssetTransaction) GetRPTListingCollection(tx *gorm.DB) ([]*models.NftyRPTListingCollection, error) {
	var rs []*models.NftyRPTListingCollection
	err := tx.Raw(`
	select collection_id, count(1) total
	from assets
	where exists(
				select 1
				from loans
				where asset_id = assets.id
					and loans.status in (?)
			)
	group by collection_id
	`,
		[]models.LoanStatus{
			models.LoanStatusNew,
		},
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return rs, nil
}

func (d *AssetTransaction) GetRPTAssetLoanToValue(tx *gorm.DB, assetID uint) (numeric.BigFloat, error) {
	var rs struct {
		UsdAmount numeric.BigFloat
	}
	err := tx.Raw(`
	select ifnull(sum(asset_transactions.amount * currencies.price *0.2), 0) usd_amount
	from asset_transactions join currencies on asset_transactions.currency_id = currencies.id
	where asset_transactions.asset_id = ?
		and asset_transactions.transaction_at >= adddate(now(), interval -12 month)
	`,
		assetID,
	).Find(&rs).Error
	if err != nil {
		return numeric.BigFloat{*big.NewFloat(0)}, errs.NewError(err)
	}
	return rs.UsdAmount, nil
}

func (d *AssetTransaction) GetAssetAvgPrice(tx *gorm.DB, assetID uint) (numeric.BigFloat, error) {
	var rs struct {
		AvgPrice numeric.BigFloat
	}
	err := tx.Raw(`
	select ifnull(sum(asset_transactions.amount), 0) avg_price
	from asset_transactions
	where asset_transactions.asset_id = ?
		and asset_transactions.transaction_at >= adddate(now(), interval -12 month)
	`,
		assetID,
	).Find(&rs).Error
	if err != nil {
		return numeric.BigFloat{*big.NewFloat(0)}, errs.NewError(err)
	}
	return rs.AvgPrice, nil
}
