package db

import (
	"fmt"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sqlx.DB

func LoadEnv(paths ...string) error {
	if len(paths) == 0 {
		return godotenv.Load()
	}
	return godotenv.Load(paths...)
}

func InitDBFromEnv() error {
	driver := os.Getenv("DB_DRIVER")
	dsn := os.Getenv("DB_DSN")
	if driver == "" || dsn == "" {
		return fmt.Errorf("DB_DRIVER and DB_DSN must be set in .env file")
	}
	return InitDB(driver, dsn)
}

func InitDB(driver string, dsn string) error {
	if driver == "sqlite3" {
		dir := "data"
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create data directory: %w", err)
		}
	}

	var err error
	DB, err = sqlx.Connect(driver, dsn)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	if driver == "mysql" {
		DB.SetMaxOpenConns(25)
		DB.SetMaxIdleConns(5)
	}

	if driver == "sqlite3" {
		if err := enableSQLiteWAL(); err != nil {
			return fmt.Errorf("failed to enable SQLite WAL: %w", err)
		}
	}

	if err := initTables(driver); err != nil {
		return fmt.Errorf("failed to init tables: %w", err)
	}

	return nil
}

func enableSQLiteWAL() error {
	if _, err := DB.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return err
	}
	if _, err := DB.Exec("PRAGMA synchronous=NORMAL"); err != nil {
		return err
	}
	if _, err := DB.Exec("PRAGMA foreign_keys=ON"); err != nil {
		return err
	}
	return nil
}

func initTables(driver string) error {
	if err := createMessageTable(driver); err != nil {
		return err
	}
	if err := createUserTable(driver); err != nil {
		return err
	}
	return nil
}

func createMessageTable(driver string) error {
	var tableName string
	var existsQuery string
	var createQuery string

	if driver == "mysql" {
		tableName = "message"
		existsQuery = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
		createQuery = `
			CREATE TABLE IF NOT EXISTS message (
				id INT AUTO_INCREMENT PRIMARY KEY,
				msg TEXT NOT NULL
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	} else {
		tableName = "message"
		existsQuery = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		createQuery = `
			CREATE TABLE IF NOT EXISTS message (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				msg TEXT NOT NULL
			);
		`
	}

	var count int
	if err := DB.Get(&count, existsQuery, tableName); err != nil {
		return fmt.Errorf("failed to check message table: %w", err)
	}

	if count == 0 {
		if _, err := DB.Exec(createQuery); err != nil {
			return fmt.Errorf("failed to create message table: %w", err)
		}
	}

	return nil
}

func createUserTable(driver string) error {
	var tableName string
	var existsQuery string
	var createQuery string

	if driver == "mysql" {
		tableName = "users"
		existsQuery = "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema = DATABASE() AND table_name = ?"
		createQuery = `
			CREATE TABLE IF NOT EXISTS users (
				id INT AUTO_INCREMENT PRIMARY KEY,
				name VARCHAR(255) NOT NULL
			) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
		`
	} else {
		tableName = "users"
		existsQuery = "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		createQuery = `
			CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name TEXT NOT NULL
			);
		`
	}

	var count int
	if err := DB.Get(&count, existsQuery, tableName); err != nil {
		return fmt.Errorf("failed to check users table: %w", err)
	}

	if count == 0 {
		if _, err := DB.Exec(createQuery); err != nil {
			return fmt.Errorf("failed to create users table: %w", err)
		}
	}

	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
