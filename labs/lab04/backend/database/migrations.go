package database

import (
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"
)

// RunMigrations runs database migrations using goose
func RunMigrations(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}

	// Set goose dialect for SQLite
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../migrations"
	if err := goose.Up(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %v", err)
	}

	return nil
}

// RollbackMigration rolls back the last migration using goose
func RollbackMigration(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}

	// Ensure we’re using the right dialect
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../migrations"
	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback last migration: %v", err)
	}

	return nil
}

// GetMigrationStatus checks migration status using goose
func GetMigrationStatus(db *sql.DB) error {
	if db == nil {
		return fmt.Errorf("database connection cannot be nil")
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../migrations"
	if err := goose.Status(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to get migration status: %v", err)
	}

	return nil
}

// CreateMigration creates a new SQL migration file with the given name
func CreateMigration(name string) error {
	if name == "" {
		return fmt.Errorf("migration name cannot be empty")
	}

	// Set dialect so goose knows how to format the timestamp/version
	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %v", err)
	}

	migrationsDir := "../migrations"
	// Create an SQL migration (timestamped) in the migrations directory
	if err := goose.Create(nil, migrationsDir, name, "sql"); err != nil {
		return fmt.Errorf("failed to create migration %q: %v", name, err)
	}

	return nil
}
