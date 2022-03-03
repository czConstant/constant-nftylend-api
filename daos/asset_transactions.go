package daos

import (
	"github.com/czConstant/constant-nftlend-api/errs"
	"github.com/czConstant/constant-nftlend-api/models"
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
