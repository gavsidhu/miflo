package database

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
)

type Database interface {
	ApplyMigration(ctx context.Context, tx *sql.Tx, migrationName string, cwd string) error
	RecordMigration(ctx context.Context, tx *sql.Tx, migrationName string, batchNum int) error
	RevertMigration(ctx context.Context, tx *sql.Tx, migrationName string, cwd string) error
	DeleteMigration(ctx context.Context, tx *sql.Tx, batchNum int) error
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
	GetNextBatchNumber() (int, error)
	GetLastBatchNumber() (int, error)
	GetAppliedMigrations() (*sql.Rows, error)
	GetUnappliedMigrations(cwd string) ([]string, error)
	GetMigrationsToRevert(batch int) ([]string, error)
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	Close() error
}

func NewDatabase(databaseURL string) (Database, error) {
	u, err := url.Parse(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("error parsing database URL: %w", err)
	}

	switch u.Scheme {
	case "sqlite":
		path := strings.TrimPrefix(databaseURL, "sqlite:")
		db, err := sql.Open("sqlite3", path)
		if err != nil {
			return nil, fmt.Errorf("error opening SQLite database: %w", err)
		}

		sqliteDB := &SQLiteDB{db}
		if err := sqliteDB.ensureMigrationsTable(); err != nil {
			return nil, fmt.Errorf("error setting up SQLite migrations table: %w", err)
		}
		return sqliteDB, nil

	case "postgresql", "postgres":
		db, err := sql.Open("postgres", databaseURL)
		if err != nil {
			return nil, fmt.Errorf("error opening PostgreSQL database: %w", err)
		}

		postgresDB := &PostgresDB{db}
		if err := postgresDB.ensureMigrationsTable(); err != nil {
			return nil, fmt.Errorf("error setting up PostgreSQL migrations table: %w", err)
		}
		return postgresDB, nil

	default:
		return nil, fmt.Errorf("unsupported database type: %s", u.Scheme)
	}

}
