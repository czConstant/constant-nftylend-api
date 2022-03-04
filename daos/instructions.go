package daos

import (
	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/jinzhu/gorm"
)

type Instruction struct {
	DAO
}

func (d *Instruction) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.Instruction, error) {
	var m models.Instruction
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *Instruction) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.Instruction, error) {
	var m models.Instruction
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *Instruction) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.Instruction, error) {
	var ms []*models.Instruction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *Instruction) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.Instruction, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.Instruction
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.Instruction{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}
