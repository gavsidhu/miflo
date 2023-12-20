package miflo

import (
	"context"
	"fmt"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/helpers"
)

func ApplyMigrations(db database.Database, ctx context.Context, cwd string) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	batchNum, err := db.GetNextBatchNumber()
	if err != nil {
		return fmt.Errorf("error getting next batch number: %w", err)
	}

	pendingMigrations, err := db.GetUnappliedMigrations()
	if err != nil {
		return fmt.Errorf("error retrieving unapplied migrations: %w", err)
	}

	if len(pendingMigrations) < 1 {
		fmt.Println("no pending migrations to apply")
		return nil
	}

	helpers.SortDirMigrations(pendingMigrations, true)

	for _, migration := range pendingMigrations {
		if err := db.ApplyMigration(ctx, tx, migration, cwd); err != nil {
			return err
		}

		if err := db.RecordMigration(ctx, tx, migration, batchNum); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	fmt.Println("Migrations applied successfully")

	return nil
}
