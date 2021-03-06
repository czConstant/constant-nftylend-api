package databases

import (
	"github.com/czConstant/constant-nftylend-api/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/jinzhu/gorm"
)

// Init : config
func Init(dbURL string, migrateFunc func(db *gorm.DB) error, idleNum int, openNum int, debug bool) (*gorm.DB, error) {
	dbConn, err := gormtrace.Open("mysql", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "gorm.Open")
	}
	dbConn.LogMode(debug)
	dbConn = dbConn.Set("gorm:save_associations", false)
	dbConn = dbConn.Set("gorm:association_save_reference", false)
	dbConn.DB().SetMaxIdleConns(idleNum)
	dbConn.DB().SetMaxOpenConns(openNum)
	if migrateFunc != nil {
		err = migrateFunc(dbConn)
		if err != nil {
			return dbConn, err
		}
	}
	return dbConn, nil
}

func MigrateDBMain(db *gorm.DB) error {
	allTables := []interface{}{
		(*models.User)(nil),
		(*models.Currency)(nil),
		(*models.Asset)(nil),
		(*models.AssetTransaction)(nil),
		(*models.Collection)(nil),
		(*models.CollectionSubmission)(nil),
		(*models.Loan)(nil),
		(*models.LoanOffer)(nil),
		(*models.LoanTransaction)(nil),
		(*models.Instruction)(nil),
		(*models.Proposal)(nil),
		(*models.ProposalChoice)(nil),
		(*models.ProposalVote)(nil),
		(*models.NotificationTemplate)(nil),
		(*models.Notification)(nil),
		(*models.UserBalance)(nil),
		(*models.UserBalanceTransaction)(nil),
		(*models.UserBalanceHistory)(nil),
		(*models.IncentiveProgram)(nil),
		(*models.IncentiveProgramDetail)(nil),
		(*models.IncentiveTransaction)(nil),
		(*models.Verification)(nil),
		(*models.AffiliateSubmission)(nil),
		(*models.Leaderboard)(nil),
	}
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		return err
	}
	db.Model(&models.Collection{}).AddUniqueIndex("collections_main_uindex", "seo_url")
	db.Model(&models.Collection{}).AddIndex("collections_creator_index", "creator")

	db.Model(&models.CollectionSubmission{}).AddIndex("collection_submissions_creator_index", "creator")

	db.Model(&models.Asset{}).AddUniqueIndex("assets_main_uindex", "seo_url")
	db.Model(&models.Asset{}).AddUniqueIndex("assets_info_uindex", "network", "contract_address", "token_id")
	db.Model(&models.Asset{}).AddIndex("assets_collection_id_index", "collection_id")
	db.Model(&models.Asset{}).AddIndex("assets_search_text_index", "search_text")

	db.Model(&models.AssetTransaction{}).AddUniqueIndex("asset_transactions_main_uindex", "asset_id", "transaction_at")
	db.Model(&models.AssetTransaction{}).AddIndex("asset_transactions_asset_id_index", "asset_id")

	db.Model(&models.Loan{}).AddUniqueIndex("loans_main_uindex", "network", "borrower_user_id", "asset_id", "started_at")
	db.Model(&models.Loan{}).AddIndex("loans_collection_id_index", "collection_id")
	db.Model(&models.Loan{}).AddIndex("loans_asset_id_index", "asset_id")
	db.Model(&models.Loan{}).AddIndex("loans_borrower_user_id_index", "borrower_user_id")
	db.Model(&models.Loan{}).AddIndex("loans_lender_user_id_index", "lender_user_id")
	db.Model(&models.Loan{}).AddIndex("loans_offer_expired_at_index", "offer_expired_at")
	db.Model(&models.Loan{}).AddIndex("loans_offer_overdue_at_index", "offer_overdue_at")

	db.Model(&models.LoanOffer{}).AddUniqueIndex("loan_offers_main_index", "loan_id", "nonce_hex")
	db.Model(&models.LoanOffer{}).AddIndex("loan_offers_loan_id_index", "loan_id")
	db.Model(&models.LoanOffer{}).AddIndex("loan_offers_lender_user_id_index", "lender_user_id")

	db.Model(&models.LoanTransaction{}).AddIndex("loan_transactions_loan_id_index", "loan_id")

	db.Model(&models.NotificationTemplate{}).AddUniqueIndex("notification_templates_type_uindex", "type")

	db.Model(&models.Notification{}).AddIndex("notifications_user_id_index", "user_id")

	db.Model(&models.Proposal{}).AddIndex("proposals_user_id_index", "user_id")
	db.Model(&models.Proposal{}).AddIndex("proposals_type_index", "type")
	db.Model(&models.Proposal{}).AddIndex("proposals_status_index", "status")
	db.Model(&models.Proposal{}).AddIndex("proposals_start_index", "start")
	db.Model(&models.Proposal{}).AddIndex("proposals_end_index", "end")

	db.Model(&models.ProposalChoice{}).AddIndex("proposal_choices_proposal_id_index", "proposal_id")

	db.Model(&models.ProposalVote{}).AddIndex("proposal_votes_proposal_id_index", "proposal_id")
	db.Model(&models.ProposalVote{}).AddIndex("proposal_votes_user_id_index", "user_id")

	db.Model(&models.User{}).AddUniqueIndex("users_main_uindex", "network", "address_checked")
	db.Model(&models.User{}).AddUniqueIndex("users_referral_uindex", "network", "username")

	db.Model(&models.UserBalance{}).AddUniqueIndex("user_balances_main_uindex", "user_id", "currency_id")
	db.Model(&models.UserBalance{}).AddIndex("user_balances_user_id_index", "user_id")

	db.Model(&models.UserBalanceTransaction{}).AddIndex("user_balance_transactions_user_id_index", "user_id")
	db.Model(&models.UserBalanceTransaction{}).AddIndex("user_balance_transactions_user_balance_id_index", "user_balance_id")

	db.Model(&models.UserBalanceHistory{}).AddUniqueIndex("user_balance_histories_main_uindex", "user_balance_id", "type", "reference")

	db.Model(&models.IncentiveTransaction{}).AddUniqueIndex("incentive_transactions_main_uindex", "user_id", "incentive_program_id", "type", "loan_id", "version")
	db.Model(&models.IncentiveTransaction{}).AddIndex("incentive_transactions_user_id_index", "user_id")
	db.Model(&models.IncentiveTransaction{}).AddIndex("incentive_transactions_type_index", "type")
	db.Model(&models.IncentiveTransaction{}).AddIndex("incentive_transactions_lock_until_at_index", "lock_until_at")

	db.Model(&models.Verification{}).AddUniqueIndex("verifications_token_uindex", "token")
	db.Model(&models.Verification{}).AddIndex("verifications_created_at_index", "created_at")
	db.Model(&models.Verification{}).AddIndex("verifications_user_id_index", "user_id")

	db.Model(&models.Leaderboard{}).AddUniqueIndex("leaderboards_rpt_date_uindex", "rpt_date")

	return nil
}
