package databases

import (
	"fmt"

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
		(*models.CollectionSubmitted)(nil),
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
		(*models.UserBalanceHistory)(nil),
		(*models.IncentiveProgram)(nil),
		(*models.IncentiveProgramDetail)(nil),
		(*models.IncentiveTransaction)(nil),
	}
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		fmt.Println(err)
		return err
	}
	db.Model(&models.Collection{}).AddUniqueIndex("collections_main_uindex", "seo_url")
	db.Model(&models.Asset{}).AddUniqueIndex("assets_main_uindex", "seo_url")
	db.Model(&models.User{}).AddUniqueIndex("users_main_uindex", "network", "address_checked")
	db.Model(&models.UserBalance{}).AddUniqueIndex("user_balances_main_uindex", "network", "address_checked", "currency_id")
	db.Model(&models.UserBalanceHistory{}).AddUniqueIndex("user_balance_histories_main_uindex", "user_balance_id", "type", "reference")
	db.Model(&models.IncentiveTransaction{}).AddUniqueIndex("incentive_transactions_main_uindex", "network", "incentive_program_id", "type", "address", "loan_id")
	return nil
}
