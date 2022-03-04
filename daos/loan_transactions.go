package daos

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type LoanTransaction struct {
	DAO
}

func (d *LoanTransaction) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.LoanTransaction, error) {
	var m models.LoanTransaction
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *LoanTransaction) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.LoanTransaction, error) {
	var m models.LoanTransaction
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *LoanTransaction) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.LoanTransaction, error) {
	var ms []*models.LoanTransaction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *LoanTransaction) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.LoanTransaction, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.LoanTransaction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.LoanTransaction{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}
