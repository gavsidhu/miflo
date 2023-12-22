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

type PostgresDB struct {
	*sql.DB
}

func (db *PostgresDB) ApplyMigration(ctx context.Context, tx *sql.Tx, migrationName string, cwd string) error {
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

func (db *PostgresDB) RecordMigration(ctx context.Context, tx *sql.Tx, migrationName string, batchNum int) error {
	if _, err := tx.ExecContext(ctx, "INSERT INTO migrations (name, batch, applied) VALUES ($1, $2, $3)", migrationName, batchNum, true); err != nil {
		return fmt.Errorf("error executing migration row insert: %w", err)
	}

	return nil
}

func (db *PostgresDB) RevertMigration(ctx context.Context, tx *sql.Tx, migrationName string, cwd string) error {
	downFilePath := path.Join(cwd, "migrations", migrationName, "down.sql")
	sqlBytes, err := os.ReadFile(downFilePath)
	if err != nil {
		return fmt.Errorf("error reading SQL file %s: %w", downFilePath, err)
	}

	_, err = tx.ExecContext(ctx, string(sqlBytes))
	if err != nil {
		return fmt.Errorf("error executing migration down file: %s %w", downFilePath, err)
	}

	return nil
}

func (db *PostgresDB) DeleteMigration(ctx context.Context, tx *sql.Tx, batchNum int) error {
	_, err := tx.ExecContext(ctx, "DELETE FROM migrations where batch = $1", batchNum)
	if err != nil {
		return fmt.Errorf("error exectuting migration row delete: %w", err)
	}

	return nil
}

func (db *PostgresDB) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return db.DB.BeginTx(ctx, opts)
}

func (db *PostgresDB) Close() error {
	return db.DB.Close()
}

func (db *PostgresDB) GetUnappliedMigrations(cwd string) ([]string, error) {
	dirMigrations, err := helpers.GetDirMigrations(cwd)
	if err != nil {
		return nil, err
	}

	appliedMigrationsRows, err := db.GetAppliedMigrations()
	if err != nil {
		fmt.Println("error getting applied mitgrations:", err)
		return nil, fmt.Errorf("error getting applied migrations: %w", err)
	}

	defer appliedMigrationsRows.Close()

	appliedMigrations, err := helpers.GetAppliedMigrationNames(appliedMigrationsRows)
	if err != nil {
	}

	var pendingMigrations []string
	for _, migration := range dirMigrations {
		if !helpers.Contains(appliedMigrations, migration) {
			pendingMigrations = append(pendingMigrations, migration)
		}
	}

	return pendingMigrations, nil
}

func (db *PostgresDB) GetAppliedMigrations() (*sql.Rows, error) {
	rows, err := db.Query("SELECT name FROM migrations WHERE applied = TRUE")
	if err != nil {
		fmt.Println("Error getting applied migrations:", err)
		return nil, fmt.Errorf("error querying for applied migrations: %w", err)
	}

	return rows, nil
}

func (db *PostgresDB) GetMigrationsToRevert(batch int) ([]string, error) {
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
		fmt.Printf("error during iteration over applied migrations: %v\n", err)
		return nil, err
	}

	return migrationsToRevert, nil
}

func (db *PostgresDB) GetAppliedMigrationsByBatch(batch int) (*sql.Rows, error) {
	query := "SELECT name FROM migrations WHERE applied = true AND batch = $1"
	rows, err := db.Query(query, batch)
	if err != nil {
		fmt.Printf("error getting applied migrations for batch %d: %v\n", batch, err)
		return nil, err
	}

	return rows, nil
}

func (db *PostgresDB) GetNextBatchNumber() (int, error) {
	var maxBatchNum int
	err := db.QueryRow("SELECT COALESCE(MAX(batch), 0) + 1 FROM migrations").Scan(&maxBatchNum)
	if err != nil {
		return 0, err
	}
	return maxBatchNum, nil
}

func (db *PostgresDB) GetLastBatchNumber() (int, error) {
	var lastBatchNum int
	err := db.QueryRow("SELECT COALESCE(MAX(batch), 0) FROM migrations").Scan(&lastBatchNum)
	if err != nil {
		return 0, err
	}

	return lastBatchNum, nil

}

func (db *PostgresDB) ensureMigrationsTable() error {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS migrations (
        id SERIAL PRIMARY KEY,
        name VARCHAR(255) UNIQUE NOT NULL,
        batch INTEGER NOT NULL,
        applied BOOLEAN NOT NULL,
        applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
    );`
	_, err := db.Exec(createTableSQL)
	return err
}
