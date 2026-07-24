package fixtures

import (
	"testing"

	"github.com/TNJKL/bookmark-management/pkg/sqldb"
	"gorm.io/gorm"
)

// Fixture defines the contract to setup, migrate, and seed data for the test database.
type Fixture interface {
	SetupDB(db *gorm.DB)
	Migrate() error
	GenerateData() error
	DB() *gorm.DB
}

// base is a helper struct that implements the DB storage for the Fixture interface.
type base struct {
	db *gorm.DB
}

// SetupDB assigns the database connection to the base fixture.
func (b *base) SetupDB(db *gorm.DB) {
	b.db = db
}

// DB returns the database connection of the base fixture.
func (b *base) DB() *gorm.DB {
	return b.db
}

// NewFixture initializes the test database, runs schema migrations, inserts sample data, and returns the configured DB connection.
func NewFixture(t *testing.T, fix Fixture) *gorm.DB {
	//create test DB
	fix.SetupDB(sqldb.InitMockDB(t))
	//migrate schema
	err := fix.Migrate()
	if err != nil {
		t.Fatalf("failed to migrate database: %s", err)
	}
	//generate data
	err = fix.GenerateData()
	if err != nil {
		t.Fatalf("failed to generate data: %s", err)
	}
	//return DB
	return fix.DB()
}
