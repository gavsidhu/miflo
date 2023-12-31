package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"

	"github.com/gavsidhu/miflo/internal/helpers"
)

type SQLiteDB struct {
	*sql.DB
}

func (db *SQLiteDB) ApplyMigration(ctx context.Context, tx *sql.Tx, migrationName string, cwd string) error {
	upFilePath := path.Join(cwd, "migrations", migrationName, "up.sql")
	sqlBytes, err := os.ReadFile(upFilePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file %s: %w", upFilePath, err)
	}

	if _, err = tx.ExecContext(ctx, string(sqlBytes)); err != nil {
		return fmt.Errorf("error executing migration file %s: %w", upFilePath, err)
	}

	return nil
}

func (db *SQLiteDB) RecordMigration(ctx context.Context, tx *sql.Tx, migrationName string, batchNum int) error {
	if _, err := tx.ExecContext(ctx, "INSERT INTO miflo_migrations (name, batch, applied) VALUES (?, ?, ?)", migrationName, batchNum, true); err != nil {
		return fmt.Errorf("error executing migration row insert: %w", err)
	}

	return nil
}

func (db *SQLiteDB) RevertMigration(ctx context.Context, tx *sql.Tx, migrationName string, cwd string) error {
	downFilePath := path.Join(cwd, "migrations", migrationName, "down.sql")
	sqlBytes, err := os.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file %s: %w", downFilePath, err)
	}

	_, err = tx.ExecContext(ctx, string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error exectuing migration down file: %s %w", downFilePath, err)
	}

	return nil
}

func (db *SQLiteDB) DeleteMigration(ctx context.Context, tx *sql.Tx, batchNum int) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM miflo_migrations where batch = ?", batchNum)
	if err != nil {
		return fmt.Errorf("error executing migration row delete: %w", err)
	}

	return nil
}

func (db *SQLiteDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

func (db *SQLiteDB) Close() error {
	return db.DB.Close()
}

func (db *SQLiteDB) GetUnappliedMigrations(cwd string) ([]string, error) {
	dirMigrations, err := helpers.GetDirMigrations(cwd)
	if err != nil {
		return nil, err
	}

	appliedMigrationsRows, err := db.GetAppliedMigrations()
	if err != nil {
		fmt.Println("error getting applied migrations:", err)
		return nil, fmt.Errorf("error getting applied migrations: %w", err)
	}

	defer appliedMigrationsRows.Close()

	appliedMigrations, err := helpers.GetAppliedMigrationNames(appliedMigrationsRows)
	if err != nil {
		return nil, err
	}

	var pendingMigrations []string
	for _, migration := range dirMigrations {
		if !helpers.Contains(appliedMigrations, migration) {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	return pendingMigrations, nil
}

func (db *SQLiteDB) GetAppliedMigrations() (*sql.Rows, error) {
	rows, err := db.Query("SELECT name FROM miflo_migrations WHERE applied = TRUE")
	if err != nil {
		fmt.Println("Error getting applied migrations:", err)
		return nil, fmt.Errorf("error querying for applied migrations: %w", err)
	}

	return rows, nil
}

func (db *SQLiteDB) GetMigrationsToRevert(batch int) ([]string, error) {
	appliedMigrationsByBatch, err := db.GetAppliedMigrationsByBatch(batch)
	if err != nil {
		return nil, err
	}

	if appliedMigrationsByBatch == nil {
		return nil, errors.New("no applied migrations found")
	}

	defer appliedMigrationsByBatch.Close()

	var migrationsToRevert []string
	for appliedMigrationsByBatch.Next() {
		var name string
		if err := appliedMigrationsByBatch.Scan(&name); err != nil {
			fmt.Printf("Error scanning migration name: %v\n", err)
			return nil, err
		}
		migrationsToRevert = append(migrationsToRevert, name)
	}

	if err := appliedMigrationsByBatch.Err(); err != nil {
		fmt.Printf("Error during iteration over applied migrations: %v\n", err)
		return nil, err
	}

	return migrationsToRevert, nil
}

func (db *SQLiteDB) GetAppliedMigrationsByBatch(batch int) (*sql.Rows, error) {
	query := "SELECT name FROM miflo_migrations WHERE applied = 1 AND batch = ?"
	rows, err := db.Query(query, batch)
	if err != nil {
		fmt.Printf("Error getting applied migrations for batch %d: %v\n", batch, err)
		return nil, err
	}

	return rows, nil
}

func (db *SQLiteDB) GetNextBatchNumber() (int, error) {
	var maxBatchNum int
	err := db.QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM miflo_migrations").Scan(&maxBatchNum)
	if err != nil {
		return 0, err
	}
	return maxBatchNum, nil
}

func (db *SQLiteDB) GetLastBatchNumber() (int, error) {
	var lastBatchNum int
	err := db.QueryRow("SELECT COALESCE(MAX(batch), 0) FROM miflo_migrations").Scan(&lastBatchNum)
	if err != nil {
		return 0, err
	}

	return lastBatchNum, nil

}

func (db *SQLiteDB) ensureMigrationsTable() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS miflo_migrations (
        id INTEGER PRIMARY KEY,
        name TEXT UNIQUE NOT NULL,
        batch INTEGER NOT NULL,
        applied BOOLEAN NOT NULL,
        applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(createTableSQL)
	return err
}
