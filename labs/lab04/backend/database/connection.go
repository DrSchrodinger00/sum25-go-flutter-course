package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// Config holds database configuration
type Config struct {
	DatabasePath    string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConfig returns a default database configuration
func DefaultConfig() *Config {
	return &Config{
		DatabasePath:    "./lab04.db",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 2 * time.Minute,
	}
}

// InitDB initializes a SQLite database using the default configuration.
func InitDB() (*sql.DB, error) {
	return InitDBWithConfig(DefaultConfig())
}

// InitDBWithConfig initializes a SQLite database using the provided Config.
// It sets up the connection pool parameters and verifies the connection.
func InitDBWithConfig(config *Config) (*sql.DB, error) {
	// 1. Open the database
	db, err := sql.Open("sqlite3", config.DatabasePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite3 database at %q: %w", config.DatabasePath, err)
	}

	// 2. Configure the connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)
	db.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// 3. Verify the connection is alive
	if err := db.Ping(); err != nil {
		// If ping fails, close the opened DB to free resources
		_ = db.Close()
		return nil, fmt.Errorf("failed to ping sqlite3 database: %w", err)
	}

	return db, nil
}

// CloseDB cleanly closes the given *sql.DB. It is safe to call with a nil DB.
func CloseDB(db *sql.DB) error {
	if db == nil {
		return nil
	}
	if err := db.Close(); err != nil {
		return fmt.Errorf("error closing database: %w", err)
	}
	return nil
}
