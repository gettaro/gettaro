package database

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	// Get database configuration from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Construct the DSN string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Configure GORM
	db, err := gorm.Open(gormpostgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	migrationDSN := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	m, err := migrate.New(
		"file://database/migrations",
		migrationDSN,
	)
	if err != nil {
		log.Fatal("Failed to create migrate instance:", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	DB = db
}
