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
		(*models.CollectionSubmitted)(nil),
		(*models.Loan)(nil),
		(*models.LoanOffer)(nil),
		(*models.LoanTransaction)(nil),
		(*models.Instruction)(nil),
		(*models.NotificationTemplate)(nil),
		(*models.Notification)(nil),
	}
	if err := db.AutoMigrate(allTables...).Error; err != nil {
		return err
	}
	return nil
}
