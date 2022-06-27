package daos

import (
	"math/big"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type UserBalance struct {
	DAO
}

func (d *UserBalance) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.UserBalance, error) {
	var m models.UserBalance
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if m.GetAvaiableBalance().Cmp(big.NewFloat(0)) < 0 {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return &m, nil
}

func (d *UserBalance) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.UserBalance, error) {
	var m models.UserBalance
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	if m.GetAvaiableBalance().Cmp(big.NewFloat(0)) < 0 {
		return nil, errs.NewError(errs.ErrBadRequest)
	}
	return &m, nil
}

func (d *UserBalance) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.UserBalance, error) {
	var ms []*models.UserBalance
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *UserBalance) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.UserBalance, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.UserBalance
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.UserBalance{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}

func (d *UserBalance) BalanceChecked(tx *gorm.DB, id uint) error {
	m, err := d.FirstByID(
		tx,
		id,
		map[string][]interface{}{},
		true,
	)
	if err != nil {
		return errs.NewError(err)
	}
	if m.GetAvaiableBalance().Cmp(big.NewFloat(0)) < 0 {
		return errs.NewError(errs.ErrBadRequest)
	}
	return nil
}
