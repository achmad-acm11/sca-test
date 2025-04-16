package test

import (
	"database/sql"
	"database/sql/driver"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func DbMock(t *testing.T) (*sql.DB, *gorm.DB, sqlmock.Sqlmock) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		Conn:                 sqlDB,
		DriverName:           "postgres",
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return sqlDB, gormDB, mock
}

type AnyThing struct{}

// Match satisfies sqlmock.Argument interface
func (a AnyThing) Match(v driver.Value) bool {
	return true
}

func AnyThingArgs(size int) []driver.Value {
	var values []driver.Value
	for i := 0; i < size; i++ {
		values = append(values, AnyThing{})
	}

	return values
}
