package daos

import (
	"time"

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

func (d *Loan) FindJoin4Page(tx *gorm.DB, joins map[string][]interface{}, filters map[string][]interface{}, preloads map[string][]interface{}, orders []string, page int, limit int) ([]*models.Loan, uint, error) {
	var ms []*models.Loan
	c, err := d.findJoin4Page(tx, &models.Loan{}, &ms, joins, filters, preloads, orders, uint(page), uint(limit), false)
	if err != nil {
		return nil, 0, errs.NewError(err)
	}
	return ms, c, nil
}

func (d *Loan) GetCollectionStats(tx *gorm.DB, collectionId uint) (*models.CollectionStats, error) {
	var rs models.CollectionStats
	err := tx.Raw(`
	select (
		select sum(nll.offer_principal_amount * nll.currency_price)
		from loans nll
		where nll.collection_id = ?
		  and nll.offer_principal_amount > 0
	) total_volume,
	(
		select count(1)
		from loans nll
		where nll.collection_id = ?
		  and nll.offer_principal_amount > 0
	) total_listed,
	(
		select avg(nll.offer_principal_amount * nll.currency_price)
		from loans nll
		where nll.collection_id = ?
		  and nll.offer_principal_amount > 0
		  and nll.offer_started_at >= adddate(now(), interval -24 hour)
	) avg24h_amount,
	(
		select min(nll.offer_principal_amount * nll.currency_price)
		from loans nll
		where nll.collection_id = ?
		  and nll.offer_principal_amount > 0
	) min_amount;
	`,
		collectionId,
		collectionId,
		collectionId,
		collectionId,
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}

func (d *Loan) GetBorrowerStats(tx *gorm.DB, borrowerUserID uint) (*models.BorrowerStats, error) {
	var rs models.BorrowerStats
	err := tx.Raw(`
	select ifnull(
				sum(
						case
							when status in (
											'done',
											'liquidated',
											'expired'
								) then 1
							else 0 end
					), 0
			) total_loans,
		ifnull(
				sum(
						case
							when status in (
								'done'
								) then 1
							else 0 end
					), 0
			) total_done_loans,
		ifnull(
				sum(
						case
							when status in (
								'done',
								'liquidated',
								'expired'
								) then offer_principal_amount
							else 0 end
					), 0
			) total_volume
		from loans
		where borrower_user_id = ?;
	`,
		borrowerUserID,
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}

func (d *Loan) GetPlatformStats(tx *gorm.DB) (*models.PlatformStats, error) {
	var rs models.PlatformStats
	err := tx.Raw(`
	select ifnull(
		sum(
				case
					when loans.status in (
										  'created',
										  'done',
										  'liquidated',
										  'expired'
						) then 1
					else 0 end
			), 0
			) total_loans,
		ifnull(
				sum(
						case
							when loans.status in (
												'done',
												'liquidated',
												'expired'
								) then 1
							else 0 end
					), 0
			) total_done_loans,
		ifnull(
				sum(
						case
							when loans.status in (
												'liquidated',
												'expired'
								) then 1
							else 0 end
					), 0
			) total_defaulted_loans,
		ifnull(
				sum(
						case
							when loans.status in (
												'created',
												'done',
												'liquidated',
												'expired'
								) then loans.offer_principal_amount * loans.currency_price
							else 0 end
					), 0
			) total_volume
		from loans;
	`,
	).Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return &rs, nil
}

func (d *Loan) GetLeaderBoardByMonth(tx *gorm.DB, network models.Network, t time.Time) ([]*models.LeaderBoardData, error) {
	var rs []*models.LeaderBoardData
	err := tx.Raw(`
	select user_id,
       sum(matching_point)                 matching_point,
       sum(matched_point)                  matched_point,
       sum(matching_point + matched_point) total_point
	from (
			select *
			from (
					select borrower_user_id user_id,
							sum(case
									when lender_user_id = 0
										and (finished_at is null or finished_at >= valid_at)
										then 1
									else 0
								end)         matching_point,
							sum(case
									when lender_user_id > 0
										then 1
									else 0
								end)         matched_point
					from loans
					where 1 = 1
						and network = ?
						and created_at >= ?
						and created_at < add_date(?, interval 1 month)
					group by borrower_user_id
				) rs
			union all
			select *
			from (
					select lender_user_id user_id,
							0              matching_point,
							sum(case
									when lender_user_id > 0
										then 2
									else 0
								end)       matched_point
					from loans
					where 1 = 1
						and lender_user_id > 0
						and network = ?
						and created_at >= ?
						and created_at < add_date(?, interval 1 month)
					group by lender_user_id
				) rs
		) rs
	group by user_id
	having sum(matching_point + matched_point) > 0
	order by sum(matching_point + matched_point) desc
	`,
		network,
		t,
		t,
		network,
		t,
		t,
	).
		Preload("User").
		Find(&rs).Error
	if err != nil {
		return nil, errs.NewError(err)
	}
	return rs, nil
}
