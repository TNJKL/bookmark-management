package infrastructure

import (
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/sqldb"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"gorm.io/gorm"
)

// CreateSQLDBWithMigration creates a new sql db connection and migrates the database
func CreateSQLDBWithMigration() *gorm.DB {
	// Create sql db conn
	sqlDB, err := sqldb.NewClient("")
	utils.NoErr(err)

	err = MigrateDB(sqlDB)
	utils.NoErr(err)

	return sqlDB
}

// MigrateDB will migrate the database according to the User struct.
// It will create the table if it doesn't exist and update the schema if it's outdated.
func MigrateDB(db *gorm.DB) error { return db.AutoMigrate(&model.User{}) }
