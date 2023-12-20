package miflo

import (
	"context"
	"fmt"

	"github.com/gavsidhu/miflo/internal/database"
	"github.com/gavsidhu/miflo/internal/helpers"
)

func RevertMigrations(db database.Database, ctx context.Context, cwd string) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting transaction: %w", err)
	}
	defer tx.Rollback()

	batchNum, err := db.GetLastBatchNumber()
	if err != nil {
		return fmt.Errorf("error getting last batch number: %w", err)
	}

	migrationsToRevert, err := db.GetMigrationsToRevert(batchNum)
	if err != nil {
		return fmt.Errorf("error retrieving migrations to revert: %w", err)
	}

	if len(migrationsToRevert) < 1 {
		fmt.Println("no migrations to revert")
		return nil
	}

	helpers.SortDirMigrations(migrationsToRevert, false)

	for _, migration := range migrationsToRevert {
		if err := db.RevertMigration(ctx, tx, migration, cwd); err != nil {
			return err
		}

		if err := db.DeleteMigration(ctx, tx, batchNum); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction: %w", err)
	}

	fmt.Println("Migrations reverted successfully")

	return nil
}
