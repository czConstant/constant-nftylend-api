package daos

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type AffiliateSubmission struct {
	DAO
}

func (d *AffiliateSubmission) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.AffiliateSubmission, error) {
	var m models.AffiliateSubmission
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *AffiliateSubmission) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.AffiliateSubmission, error) {
	var m models.AffiliateSubmission
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *AffiliateSubmission) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.AffiliateSubmission, error) {
	var ms []*models.AffiliateSubmission
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *AffiliateSubmission) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.AffiliateSubmission, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.AffiliateSubmission
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.AffiliateSubmission{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}
