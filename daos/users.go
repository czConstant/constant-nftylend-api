package daos

import (
	"math/big"

	"github.com/czConstant/constant-nftylend-api/errs"
	"github.com/czConstant/constant-nftylend-api/models"
	"github.com/czConstant/constant-nftylend-api/types/numeric"
	"github.com/jinzhu/gorm"
)

type User struct {
	DAO
}

func (d *User) FirstByID(tx *gorm.DB, id uint, preloads map[string][]interface{}, forUpdate bool) (*models.User, error) {
	var m models.User
	if err := d.first(tx, &m, map[string][]interface{}{"id = ?": []interface{}{id}}, preloads, nil, forUpdate); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *User) First(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string) (*models.User, error) {
	var m models.User
	if err := d.first(tx, &m, filters, preloads, orders, false); err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &m, nil
}

func (d *User) Find(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, offset int, limit int) ([]*models.User, error) {
	var ms []*models.User
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, err
	}
	return ms, nil
}

func (d *User) Find4Page(tx *gorm.DB, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.User, uint, error) {
	var (
		offset = (page - 1) * limit
	)
	var ms []*models.User
	if err := d.find(tx, &ms, filters, preloads, orders, offset, limit, false); err != nil {
		return nil, 0, errs.NewError(err)
	}
	c, err := d.count(tx, &models.User{}, filters)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}

func (d *User) GetRPTListingCollection(tx *gorm.DB) ([]*models.NftyRPTListingCollection, error) {
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

func (d *User) GetUserBorrowStats(tx *gorm.DB, borrowerID uint) (*models.UserBorrowStats, error) {
	var rs models.UserBorrowStats
	err := tx.Raw(`
	select ifnull(sum(1), 0)                                           total_loans,
       ifnull(sum(loans.offer_principal_amount * currencies.price), 0) total_volume
	from loans
			join currencies on loans.currency_id = currencies.id
	where loans.borrower_user_id = ?
	and loans.status not in ('new', 'cancelled')
	`,
		borrowerID,
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}

func (d *User) GetUserLendStats(tx *gorm.DB, lenderID uint) (*models.UserLendStats, error) {
	var rs models.UserLendStats
	err := tx.Raw(`
	select ifnull(sum(1), 0)                                           total_loans,
       ifnull(sum(loans.offer_principal_amount * currencies.price), 0) total_volume
	from loans
			join currencies on loans.currency_id = currencies.id
	where loans.lender_user_id = ?
	and loans.status not in ('new', 'cancelled')
	`,
		lenderID,
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}

func (d *User) GetUserCreditScore(tx *gorm.DB, userID uint) (numeric.BigFloat, error) {
	var rs struct {
		TotalCredit numeric.BigFloat
	}
	err := tx.Raw(
		"call CalculateCreditPoint(?)",
		userID,
	).Find(&rs).Error
	if err != nil {
		return numeric.BigFloat{*big.NewFloat(0)}, errs.NewError(err)
	}
	return rs.TotalCredit, nil
}
