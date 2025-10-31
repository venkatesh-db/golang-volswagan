
package postgres

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	PGHost     = getenv("PG_HOST", "localhost")
	PGPort     = getenv("PG_PORT", "5432")
	PGUser     = getenv("PG_USER", "postgres")
	PGPassword = getenv("PG_PASSWORD", "postgres")
	PGDB       = getenv("PG_DB", "fleet")
)

func getenv(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

// NewPostgres creates a GORM PostgreSQL connection pool
func NewPostgres() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable TimeZone=UTC",
		PGHost, PGPort, PGUser, PGDB, PGPassword,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: false,
		Logger:                 logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	return db, nil
}


