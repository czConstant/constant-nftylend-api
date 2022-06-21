package daos

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type IncentiveTransaction struct {
	DAO
}

func (d *IncentiveTransaction) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.IncentiveTransaction, error) {
	var m models.IncentiveTransaction
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *IncentiveTransaction) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.IncentiveTransaction, error) {
	var m models.IncentiveTransaction
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *IncentiveTransaction) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.IncentiveTransaction, error) {
	var ms []*models.IncentiveTransaction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *IncentiveTransaction) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.IncentiveTransaction, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.IncentiveTransaction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.IncentiveTransaction{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}

func (d *IncentiveTransaction) GetAffiliateStats(tx *gorm.DB, userID uint, currencyID uint) (*models.AffiliateStats, error) {
	var rs models.AffiliateStats
	err := tx.Raw(`
	select sum(amount)              total_commisions,
		count(distinct ref_user_id) total_users,
		count(1)                    total_transactions
	from incentive_transactions
	where 1 = 1
		and user_id = ?
		and currency_id = ?
		and type in (?)
	`,
		userID,
		currencyID,
		[]models.IncentiveTransactionType{
			models.IncentiveTransactionTypeAffiliateBorrowerLoanDone,
			models.IncentiveTransactionTypeAffiliateLenderLoanDone,
		},
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}
