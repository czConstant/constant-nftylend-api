package daos

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type Loan struct {
	DAO
}

func (d *Loan) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.Loan, error) {
	var m models.Loan
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *Loan) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.Loan, error) {
	var m models.Loan
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *Loan) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.Loan, error) {
	var ms []*models.Loan
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *Loan) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.Loan, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.Loan
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.Loan{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}
func (d *Loan) GetRPTCollectionLoan(tx *gorm.DB, collectionId uint) (*models.NftyRPTCollectionLoan, error) {
	var rs models.NftyRPTCollectionLoan
	err := tx.Raw(`
	select (
		select sum(nll.offer_principal_amount)
		from loans nll
				 join assets nla on nll.asset_id = nla.id
		where nla.collection_id = ?
		  and nll.offer_principal_amount > 0
	) total_volume,
	(
		select count(1)
		from loans nll
				 join assets nla on nll.asset_id = nla.id
		where nla.collection_id = ?
		  and nll.offer_principal_amount > 0
	) total_listed,
	(
		select avg(nll.offer_principal_amount)
		from loans nll
				 join assets nla on nll.asset_id = nla.id
		where nla.collection_id = ?
		  and nll.offer_principal_amount > 0
		  and nll.offer_started_at >= adddate(now(), interval -24 hour)
	) avg24h_amount;
	`,
		collectionId,
		collectionId,
		collectionId,
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}
