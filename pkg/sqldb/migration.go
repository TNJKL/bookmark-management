package sqldb

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

// MigrateSQLDB initializes and executes PostgreSQL database migrations from SQL files at migrationPath
func MigrateSQLDB(db *gorm.DB, migrationPath string, mode string, steps int) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	pgDriver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(migrationPath, db.Name(), pgDriver)
	if err != nil {
		return err
	}
	return migrateSchema(m, mode, steps)
}

// migrateSchema handles migration execution modes (up, steps, down) for the migrate instance
func migrateSchema(m *migrate.Migrate, mode string, steps int) error {
	var migrationErr error

	switch mode {
	case "up":
		migrationErr = m.Up()
	case "steps":
		if steps == 0 {
			return errors.New("[Database migration] Steps must not be 0")
		}
		migrationErr = m.Steps(steps)
	case "down":
		migrationErr = m.Down()
	default:
		return errors.New("[Database migration] Invalid migration mode . Please use 'up' or 'steps'")
	}

	if migrationErr != nil && !errors.Is(migrationErr, migrate.ErrNoChange) {
		return fmt.Errorf("[Database migration] Migration error: %s", migrationErr.Error())
	}
	return nil
}
