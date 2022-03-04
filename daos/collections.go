package daos

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type Collection struct {
	DAO
}

func (d *Collection) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.Collection, error) {
	var m models.Collection
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *Collection) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.Collection, error) {
	var m models.Collection
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *Collection) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.Collection, error) {
	var ms []*models.Collection
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *Collection) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.Collection, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.Collection
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.Collection{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}
